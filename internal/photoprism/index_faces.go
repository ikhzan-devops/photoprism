package photoprism

import (
	"time"

	"github.com/dustin/go-humanize/english"

	"github.com/photoprism/photoprism/internal/ai/face"
	"github.com/photoprism/photoprism/internal/ai/vision"
	"github.com/photoprism/photoprism/internal/thumb"
	"github.com/photoprism/photoprism/pkg/clean"
)

// Faces finds faces in JPEG media files and returns them.
func (ind *Index) Faces(jpeg *MediaFile, expected int) face.Faces {
	if jpeg == nil {
		return face.Faces{}
	}

	var thumbSize thumb.Name

	// Select best thumbnail depending on configured size.
	if Config().ThumbSizePrecached() < 1280 {
		thumbSize = thumb.Fit720
	} else {
		thumbSize = thumb.Fit1280
	}

	thumbName, err := jpeg.Thumbnail(Config().ThumbCachePath(), thumbSize)

	if err != nil {
		log.Debugf("vision: %s in %s (detect faces)", err, clean.Log(jpeg.BaseName()))
		return face.Faces{}
	}

	if thumbName == "" {
		log.Debugf("vision: thumb %s not found in %s (detect faces)", thumbSize, clean.Log(jpeg.BaseName()))
		return face.Faces{}
	}

	start := time.Now()

	faces, err := vision.Faces(thumbName, Config().FaceSize(), true, expected)

	if err != nil {
		log.Debugf("vision: %s in %s (detect faces)", err, clean.Log(jpeg.BaseName()))
	}

	if thumbSize != thumb.Fit1280 {
		needRetry := len(faces) == 0

		if !needRetry && expected > 0 && len(faces) < expected {
			needRetry = true
		}

		if !needRetry && len(faces) > 0 && faces.MaxScale() < 96 {
			needRetry = true
		}

		if needRetry {
			if altThumb, altErr := jpeg.Thumbnail(Config().ThumbCachePath(), thumb.Fit1280); altErr != nil {
				log.Debugf("vision: %s in %s (detect faces @1280)", altErr, clean.Log(jpeg.BaseName()))
			} else if altThumb == "" {
				log.Debugf("vision: thumb %s not found in %s (detect faces @1280)", thumb.Fit1280, clean.Log(jpeg.BaseName()))
			} else if retryFaces, retryErr := vision.Faces(altThumb, Config().FaceSize(), true, expected); retryErr != nil {
				log.Debugf("vision: %s in %s (detect faces @1280)", retryErr, clean.Log(jpeg.BaseName()))
			} else if len(retryFaces) > 0 {
				log.Debugf("vision: retry face detection for %s using %s", clean.Log(jpeg.BaseName()), thumb.Fit1280)
				faces = retryFaces
			}
		}
	}

	if l := len(faces); l > 0 {
		log.Infof("vision: found %s in %s [%s]", english.Plural(l, "face", "faces"), clean.Log(jpeg.BaseName()), time.Since(start))
	}

	return faces
}
