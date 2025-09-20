package dl

import (
	"testing"
	"time"

	"github.com/photoprism/photoprism/pkg/fs"
)

func TestRemuxOptionsFromInfo(t *testing.T) {
	info := Info{
		Title:       " My Title ",
		AltTitle:    "Alt",
		Description: "  Desc  ",
		Artist:      "Artist Name",
		Timestamp:   float64(time.Date(2024, 12, 31, 23, 59, 58, 0, time.UTC).Unix()),
	}
	opt := RemuxOptionsFromInfo("ffmpeg", fs.VideoMp4, info, "https://example.com/v")
	if opt.Title != "My Title" {
		t.Fatalf("Title mismatch: %q", opt.Title)
	}
	if opt.Description != "Desc" {
		t.Fatalf("Description mismatch: %q", opt.Description)
	}
	if opt.Author != "Artist Name" {
		t.Fatalf("Author mismatch: %q", opt.Author)
	}
	if opt.Comment != "https://example.com/v" {
		t.Fatalf("Comment mismatch: %q", opt.Comment)
	}
	if opt.Created.IsZero() {
		t.Fatalf("Created timestamp should be set")
	}
}
