package domain

import "testing"

func TestHasChapterContent(t *testing.T) {
	for _, content := range []string{"", "<p><br></p>", "<p>&nbsp;\u200b\ufeff</p>", "<!--Nội dung-->", "<script>Code</script><style>CSS</style>", "<iframe>Fallback</iframe>", "<img>", "<img src='javascript:alert(1)'>", "<svg><text>SVG</text></svg>"} {
		if HasChapterContent(content) {
			t.Errorf("Expected empty content: %q", content)
		}
	}
	for _, content := range []string{"Văn bản", "<p><strong>Nội dung</strong></p>", "<p>&amp;</p>", "<img src='/uploads/scene.png'>"} {
		if !HasChapterContent(content) {
			t.Errorf("Expected readable content: %q", content)
		}
	}
}
