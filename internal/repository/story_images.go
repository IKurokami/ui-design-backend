package repository

import (
	"context"
	"fmt"
	"strings"
	"ui-design-backend/internal/domain"
)

// Fetch preview sources in batches rather than issuing a detail request for every story.
// Full chapter contents stay on the server; list responses carry at most three URLs.
func (r *StoryRepo) populatePreviewImages(ctx context.Context, stories []*domain.Story) error {
	byID := make(map[string]*domain.Story, len(stories))
	for _, story := range stories {
		story.PreviewImages = []string{}
		story.AddPreviewImage(story.CoverURL)
		byID[story.ID] = story
	}
	const batchSize = 500
	for offset := 0; offset < len(stories); offset += batchSize {
		end := min(offset+batchSize, len(stories))
		args := make([]any, 0, end-offset)
		placeholders := make([]string, 0, end-offset)
		for _, story := range stories[offset:end] {
			args = append(args, story.ID)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}
		rows, err := r.db.QueryContext(ctx, `SELECT story_id, COALESCE(illustrations, ''), COALESCE(content, '')
		FROM chapters WHERE story_id IN (`+strings.Join(placeholders, ",")+`)
		ORDER BY story_id, chapter_number ASC, created_at ASC, id ASC`, args...)
		if err != nil {
			return err
		}
		for rows.Next() {
			var ch domain.Chapter
			if err := rows.Scan(&ch.StoryID, &ch.Illustrations, &ch.Content); err != nil {
				rows.Close()
				return err
			}
			byID[ch.StoryID].AddChapterPreviewImages(ch)
		}
		err = rows.Err()
		rows.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
