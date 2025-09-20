package dl

import (
	"math"
	"strings"
	"time"

	"github.com/photoprism/photoprism/internal/ffmpeg/encode"
	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/fs"
)

// RemuxOptionsFromInfo builds ffmpeg remux options (container + metadata)
// based on yt-dlp Info and the source URL. The returned options enforce the
// target container and set title, description, author, comment, and created
// timestamp when provided by the extractor.
func RemuxOptionsFromInfo(ffmpegBin string, container fs.Type, info Info, sourceURL string) encode.Options {
	opt := encode.NewRemuxOptions(ffmpegBin, container, false)

	if title := clean.Name(info.Title); title != "" {
		opt.Title = title
	} else if title = clean.Name(info.AltTitle); title != "" {
		opt.Title = title
	}

	if desc := strings.TrimSpace(info.Description); desc != "" {
		opt.Description = desc
	}
	if u := strings.TrimSpace(sourceURL); u != "" {
		opt.Comment = u
	}

	if author := clean.Name(info.Artist); author != "" {
		opt.Author = author
	} else if author = clean.Name(info.AlbumArtist); author != "" {
		opt.Author = author
	} else if author = clean.Name(info.Creator); author != "" {
		opt.Author = author
	} else if author = clean.Name(info.License); author != "" {
		opt.Author = author
	}

	if info.Timestamp > 1 {
		sec, dec := math.Modf(info.Timestamp)
		opt.Created = time.Unix(int64(sec), int64(dec*(1e9)))
	}

	return opt
}
