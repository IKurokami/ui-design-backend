package domain

import (
	"encoding/json"
	"strings"

	"golang.org/x/net/html"
)

const StoryPreviewLimit = 3

func (s *Story) RefreshPreviewImages() {
	s.PreviewImages = []string{}
	s.AddPreviewImage(s.CoverURL)
	for _, ch := range s.Chapters {
		s.AddChapterPreviewImages(ch)
	}
}

// The cover comes first; additional layers must use different images from the story.
func (s *Story) AddPreviewImage(source string) {
	source = strings.TrimSpace(source)
	if len(s.PreviewImages) >= StoryPreviewLimit || !readableImageSource(source) {
		return
	}
	for _, existing := range s.PreviewImages {
		if existing == source {
			return
		}
	}
	s.PreviewImages = append(s.PreviewImages, source)
}

func (s *Story) AddChapterPreviewImages(ch Chapter) {
	if len(s.PreviewImages) >= StoryPreviewLimit {
		return
	}
	var illustrations []string
	if json.Unmarshal([]byte(ch.Illustrations), &illustrations) == nil {
		for _, source := range illustrations {
			s.AddPreviewImage(source)
		}
	}
	if len(s.PreviewImages) >= StoryPreviewLimit {
		return
	}
	// Older chapters may contain embedded images without an illustrations array.
	z := html.NewTokenizer(strings.NewReader(ch.Content))
	for len(s.PreviewImages) < StoryPreviewLimit {
		token := z.Next()
		if token == html.ErrorToken {
			return
		}
		if token != html.StartTagToken && token != html.SelfClosingTagToken {
			continue
		}
		t := z.Token()
		if t.Data != "img" {
			continue
		}
		for _, attr := range t.Attr {
			if attr.Key == "src" {
				s.AddPreviewImage(attr.Val)
			}
		}
	}
}
