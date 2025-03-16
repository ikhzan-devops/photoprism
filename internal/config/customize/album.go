package customize

import (
	"github.com/photoprism/photoprism/internal/entity/sortby"
)

// AlbumsSettings represents album defaults and preferences.
type AlbumsSettings struct {
	DefaultOrder string           `json:"defaultOrder" yaml:"DefaultOrder"`
	Download     DownloadSettings `json:"download" yaml:"Download"`
}

// NewAlbumSettings creates album settings with defaults.
func NewAlbumSettings() AlbumsSettings {
	return AlbumsSettings{
		DefaultOrder: sortby.Oldest,
		Download: DownloadSettings{
			Name:         DownloadNameShare,
			Disabled:     false,
			Originals:    true,
			MediaRaw:     false,
			MediaSidecar: false,
		},
	}
}
