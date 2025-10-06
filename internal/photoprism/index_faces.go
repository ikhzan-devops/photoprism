package photoprism

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize/english"

	"github.com/photoprism/photoprism/internal/ai/face"
	"github.com/photoprism/photoprism/internal/ai/vision"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/thumb"
	"github.com/photoprism/photoprism/pkg/clean"
)

// DetectFaces finds faces in JPEG media files and returns them.
func DetectFaces(jpeg *MediaFile, expected int) (face.Faces, error) {
	if jpeg == nil {
		return face.Faces{}, fmt.Errorf("missing media file")
	}

	engine := face.ActiveEngine()
	engineName := ""
	if engine != nil {
		engineName = engine.Name()
	}

	var thumbSize thumb.Name

	if engineName == face.EngineONNX {
		thumbSize = thumb.Fit720
	} else if Config().ThumbSizePrecached() < 1280 {
		thumbSize = thumb.Fit720
	} else {
		thumbSize = thumb.Fit1280
	}

	thumbName, err := jpeg.Thumbnail(Config().ThumbCachePath(), thumbSize)

	if err != nil {
		log.Debugf("vision: %s in %s (detect faces)", err, clean.Log(jpeg.BaseName()))
		return face.Faces{}, err
	}

	if thumbName == "" {
		log.Debugf("vision: thumb %s not found in %s (detect faces)", thumbSize, clean.Log(jpeg.BaseName()))
		return face.Faces{}, fmt.Errorf("thumbnail %s not found", thumbSize)
	}

	start := time.Now()
	var detectErr error
	allowRetry := Config().FaceEngineRetry() && (engineName == "" || engineName == face.EnginePigo)

	faces, err := vision.Faces(thumbName, Config().FaceSize(), true, expected)

	if err != nil {
		log.Debugf("vision: %s in %s (detect faces)", err, clean.Log(jpeg.BaseName()))
		detectErr = err
	}

	if allowRetry && thumbSize != thumb.Fit1280 {
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
				detectErr = retryErr
			} else if len(retryFaces) > 0 {
				log.Debugf("vision: retry face detection for %s using %s", clean.Log(jpeg.BaseName()), thumb.Fit1280)
				faces = retryFaces
				detectErr = nil
			}
		}
	}

	if l := len(faces); l > 0 {
		log.Infof("vision: found %s in %s [%s]", english.Plural(l, "face", "faces"), clean.Log(jpeg.BaseName()), time.Since(start))
	}

	return faces, detectErr
}

// ApplyDetectedFaces persists detected faces on the given file and updates face counts.
func ApplyDetectedFaces(file *entity.File, faces face.Faces) (saved bool, count int, err error) {
	if file == nil {
		return false, 0, fmt.Errorf("faces: file is nil")
	}

	if len(faces) == 0 {
		return false, 0, nil
	}

	file.AddFaces(faces)

	savedMarkers, saveErr := file.SaveMarkers()
	if saveErr != nil {
		return false, 0, saveErr
	}

	if savedMarkers == 0 {
		return false, 0, nil
	}

	count, updateErr := file.UpdatePhotoFaceCount()
	if updateErr != nil {
		return true, 0, updateErr
	}

	return true, count, nil
}

// Faces finds faces in JPEG media files and returns them.
func (ind *Index) Faces(jpeg *MediaFile, expected int) face.Faces {
	faces, _ := DetectFaces(jpeg, expected)
	return faces
}
