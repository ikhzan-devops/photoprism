package thumb

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"runtime/debug"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"

	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/fs"
)

// FileInfo returns the image header info containing width and height.
func FileInfo(fileName string) (info image.Config, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic %s while decoding %s file info\nstack: %s", r, clean.Log(fileName), debug.Stack())
		}
	}()

	// Resolve symlinks.
	if fileName, err = fs.Resolve(fileName); err != nil {
		return info, err
	}

	file, err := os.Open(fileName)

	if err != nil || file == nil {
		return info, err
	}

	defer file.Close()

	// Reset file offset.
	// see https://github.com/golang/go/issues/45902#issuecomment-1007953723
	_, err = file.Seek(0, 0)

	if err != nil {
		return info, fmt.Errorf("%s on seek", err)
	}

	// Decode image config (dimensions).
	info, _, err = image.DecodeConfig(file)

	if err != nil {
		return info, fmt.Errorf("%s while decoding file info", err)
	}

	return info, err
}
