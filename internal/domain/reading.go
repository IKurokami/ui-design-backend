package domain

import (
	"net/url"
	"strings"
	"time"
	"unicode"

	"golang.org/x/net/html"
)

type ReadingProgress struct {
	ChapterID     string     `json:"chapterId"`
	ChapterNumber int        `json:"chapterNumber,omitempty"`
	ChapterTitle  string     `json:"chapterTitle,omitempty"`
	UpdatedAt     *time.Time `json:"updatedAt,omitempty"`
}

// Ignore empty rich-text markup, nonbreaking spaces and invisible characters.
func HasChapterContent(content string) bool {
	z := html.NewTokenizer(strings.NewReader(content))
	ignored := ""
	for {
		switch z.Next() {
		case html.ErrorToken:
			return false
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if ignored != "" {
				continue
			}
			if t.Data == "script" || t.Data == "style" || t.Data == "template" || t.Data == "iframe" || t.Data == "object" || t.Data == "svg" || t.Data == "math" {
				ignored = t.Data
				continue
			}
			if t.Data == "img" {
				for _, a := range t.Attr {
					if a.Key == "src" && readableImageSource(a.Val) {
						return true
					}
				}
			}
		case html.EndTagToken:
			if z.Token().Data == ignored {
				ignored = ""
			}
		case html.TextToken:
			if ignored == "" {
				for _, r := range string(z.Text()) {
					if !unicode.IsSpace(r) && !unicode.Is(unicode.Cf, r) {
						return true
					}
				}
			}
		}
	}
}

func readableImageSource(source string) bool {
	source = strings.TrimSpace(source)
	if source == "" {
		return false
	}
	u, err := url.Parse(source)
	if err != nil {
		return false
	}
	switch strings.ToLower(u.Scheme) {
	case "", "http", "https":
		return true
	case "data":
		return strings.HasPrefix(strings.ToLower(source), "data:image/")
	default:
		return false
	}
}
