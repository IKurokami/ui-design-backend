package domain

import (
	"reflect"
	"testing"
)

func TestStoryPreviewImagesUseDistinctRealSources(t *testing.T) {
	s := Story{CoverURL: "/uploads/cover.jpg", Chapters: []Chapter{
		{Illustrations: `["/uploads/cover.jpg", " /uploads/scene-1.jpg ", "/uploads/scene-1.jpg"]`, Content: `<p><img src="/uploads/scene-2.jpg"></p><img src="/uploads/scene-3.jpg">`},
		{Illustrations: `["/uploads/later.jpg"]`},
	}}
	s.RefreshPreviewImages()
	expected := []string{"/uploads/cover.jpg", "/uploads/scene-1.jpg", "/uploads/scene-2.jpg"}
	if !reflect.DeepEqual(s.PreviewImages, expected) {
		t.Fatalf("Expected %v, got %v", expected, s.PreviewImages)
	}
	s.CoverURL = "/uploads/scene-2.jpg"
	s.RefreshPreviewImages()
	if s.PreviewImages[0] != s.CoverURL {
		t.Fatal("Changing the cover must change the first preview")
	}
}

func TestStoryPreviewImagesDoNotFillMissingLayers(t *testing.T) {
	s := Story{CoverURL: "/uploads/only.jpg", Chapters: []Chapter{{Illustrations: `["/uploads/only.jpg"]`, Content: `<img src="/uploads/only.jpg">`}}}
	s.RefreshPreviewImages()
	if !reflect.DeepEqual(s.PreviewImages, []string{"/uploads/only.jpg"}) {
		t.Fatal("A single image must remain a single image")
	}
	s = Story{Chapters: []Chapter{{Illustrations: "malformed JSON", Content: `<img src="javascript:alert(1)"><img src="/uploads/embedded.jpg">`}}}
	s.RefreshPreviewImages()
	if !reflect.DeepEqual(s.PreviewImages, []string{"/uploads/embedded.jpg"}) {
		t.Fatal("Embedded images must work without stored illustrations")
	}
	s = Story{}
	s.RefreshPreviewImages()
	if len(s.PreviewImages) != 0 {
		t.Fatal("An imageless story must not invent previews")
	}
}
