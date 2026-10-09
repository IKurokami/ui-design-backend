package repository

import (
	"bytes"
	"database/sql"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Remove retired presentation metadata while preserving chapters and reading history.
func (db *DB) migrateChapterOnly() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var count int
	if err = tx.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = 'reader-chapter-only-v1'`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return tx.Commit()
	}
	for _, target := range []struct {
		table   string
		columns []string
	}{
		{"stories", []string{"book_layout"}},
		{"reading_progress", []string{"mode", "page_index", "block_index", "text_offset", "scroll_progress"}},
	} {
		columns, err := schemaColumns(tx, db.Driver, target.table)
		if err != nil {
			return err
		}
		for _, column := range target.columns {
			if columns[column] {
				if _, err = tx.Exec("ALTER TABLE " + target.table + " DROP COLUMN " + column); err != nil {
					return err
				}
			}
		}
	}
	rows, err := tx.Query(`SELECT id, content FROM chapters WHERE content LIKE '%data-book-%'`)
	if err != nil {
		return err
	}
	type chapter struct{ id, content string }
	var chapters []chapter
	for rows.Next() {
		var ch chapter
		if err = rows.Scan(&ch.id, &ch.content); err != nil {
			rows.Close()
			return err
		}
		chapters = append(chapters, ch)
	}
	if err = rows.Err(); err != nil {
		rows.Close()
		return err
	}
	if err = rows.Close(); err != nil {
		return err
	}
	for _, ch := range chapters {
		content, err := plainChapterContent(ch.content)
		if err != nil {
			return err
		}
		if _, err = tx.Exec(`UPDATE chapters SET content = $1 WHERE id = $2`, content, ch.id); err != nil {
			return err
		}
	}
	if _, err = tx.Exec(`DELETE FROM schema_migrations WHERE version = 'reader-book-v1'`); err != nil {
		return err
	}
	if _, err = tx.Exec(`INSERT INTO schema_migrations(version) VALUES ('reader-chapter-only-v1')`); err != nil {
		return err
	}
	return tx.Commit()
}

func schemaColumns(tx *sql.Tx, driver, table string) (map[string]bool, error) {
	var rows *sql.Rows
	var err error
	if driver == "sqlite" {
		rows, err = tx.Query("PRAGMA table_info(" + table + ")")
	} else {
		rows, err = tx.Query(`SELECT column_name FROM information_schema.columns WHERE table_schema = current_schema() AND table_name = $1`, table)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns := make(map[string]bool)
	for rows.Next() {
		var name string
		if driver == "sqlite" {
			var index, required, key int
			var kind string
			var fallback sql.NullString
			err = rows.Scan(&index, &name, &kind, &required, &fallback, &key)
		} else {
			err = rows.Scan(&name)
		}
		if err != nil {
			return nil, err
		}
		columns[name] = true
	}
	return columns, rows.Err()
}

func plainChapterContent(content string) (string, error) {
	root := &html.Node{Type: html.ElementNode, Data: "div", DataAtom: atom.Div}
	nodes, err := html.ParseFragment(strings.NewReader(content), root)
	if err != nil {
		return "", err
	}
	for _, node := range nodes {
		root.AppendChild(node)
	}
	var clean func(*html.Node, bool)
	clean = func(node *html.Node, controlled bool) {
		unwrap := false
		for _, attr := range node.Attr {
			if strings.HasPrefix(attr.Key, "data-book-") {
				controlled = true
			}
			if attr.Key == "data-book-frame" || attr.Key == "data-book-tiles" || attr.Key == "data-book-tile-text" {
				unwrap = true
			}
		}
		attrs := node.Attr[:0]
		for _, attr := range node.Attr {
			if strings.HasPrefix(attr.Key, "data-book-") || attr.Key == "data-block" || (controlled && (attr.Key == "style" || attr.Key == "contenteditable")) {
				continue
			}
			attrs = append(attrs, attr)
		}
		node.Attr = attrs
		for child := node.FirstChild; child != nil; {
			next := child.NextSibling
			clean(child, controlled)
			child = next
		}
		if unwrap && node.Parent != nil {
			parent := node.Parent
			for node.FirstChild != nil {
				child := node.FirstChild
				node.RemoveChild(child)
				parent.InsertBefore(child, node)
			}
			parent.RemoveChild(node)
		}
	}
	clean(root, false)
	var result bytes.Buffer
	for node := root.FirstChild; node != nil; node = node.NextSibling {
		if err = html.Render(&result, node); err != nil {
			return "", err
		}
	}
	return result.String(), nil
}
