package batch

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/entity/search"
	"github.com/photoprism/photoprism/pkg/fs"
)

func TestNewPhotosForm(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var photos search.PhotoResults

		dataFile := fs.Abs("./testdata/photos.json")
		data, dataErr := os.ReadFile(dataFile)

		if dataErr != nil {
			t.Fatal(dataErr)
		}

		jsonErr := json.Unmarshal(data, &photos)

		if jsonErr != nil {
			t.Fatal(jsonErr)
		}

		frm := NewPhotosForm(photos)

		// Photo metadata.
		assert.Equal(t, "", frm.PhotoType.Value)
		assert.Equal(t, true, frm.PhotoType.Mixed)
		assert.Equal(t, "", frm.PhotoTitle.Value)
		assert.Equal(t, true, frm.PhotoTitle.Mixed)
		assert.Equal(t, "", frm.PhotoCaption.Value)
		assert.Equal(t, true, frm.PhotoCaption.Mixed)
		assert.Equal(t, false, frm.PhotoFavorite.Value)
		assert.Equal(t, true, frm.PhotoFavorite.Mixed)
		assert.Equal(t, false, frm.PhotoPrivate.Value)
		assert.Equal(t, false, frm.PhotoPrivate.Mixed)
		assert.Equal(t, 1000003, frm.CameraID.Value)
		assert.Equal(t, false, frm.CameraID.Mixed)
		assert.Equal(t, 1000000, frm.LensID.Value)
		assert.Equal(t, false, frm.LensID.Mixed)
		assert.Equal(t, 0, frm.PhotoIso.Value)
		assert.Equal(t, true, frm.PhotoIso.Mixed)
		assert.Equal(t, float32(0), frm.PhotoFNumber.Value)
		assert.Equal(t, true, frm.PhotoFNumber.Mixed)
		assert.Equal(t, 0, frm.PhotoFocalLength.Value)
		assert.Equal(t, true, frm.PhotoFocalLength.Mixed)

		// Additional details.
		assert.Equal(t, "", frm.DetailsKeywords.Value)
		assert.Equal(t, true, frm.DetailsKeywords.Mixed)
		assert.Equal(t, "", frm.DetailsSubject.Value)
		assert.Equal(t, true, frm.DetailsSubject.Mixed)
		assert.Equal(t, "", frm.DetailsArtist.Value)
		assert.Equal(t, true, frm.DetailsArtist.Mixed)
		assert.Equal(t, "", frm.DetailsCopyright.Value)
		assert.Equal(t, true, frm.DetailsCopyright.Mixed)
		assert.Equal(t, "", frm.DetailsLicense.Value)
		assert.Equal(t, true, frm.DetailsLicense.Mixed)
	})
	t.Run("Success", func(t *testing.T) {
		photo1 := search.Photo{
			ID:            111115411,
			TakenSrc:      "",
			TimeZone:      "",
			PhotoUID:      "ps6sg6be2lvl0x41",
			PhotoType:     "image",
			PhotoTitle:    "Same Title",
			PhotoCountry:  "de",
			PhotoPrivate:  true,
			PhotoPanorama: true,
			PhotoScan:     true,
			PhotoFavorite: false,
			CameraID:      1,
			LensID:        2,
		}

		photo2 := search.Photo{
			ID:            111115511,
			CreatedAt:     time.Time{},
			TakenAt:       time.Time{},
			TakenAtLocal:  time.Time{},
			TakenSrc:      "",
			TimeZone:      "",
			PhotoUID:      "ps6sg6be2lvlx986",
			PhotoType:     "image",
			PhotoTitle:    "Same Title",
			PhotoCountry:  "ca",
			PhotoPrivate:  false,
			PhotoPanorama: false,
			PhotoScan:     false,
			PhotoFavorite: true,
			CameraID:      1,
			LensID:        2,
		}

		photos := search.PhotoResults{photo1, photo2}

		frm := NewPhotosForm(photos)

		// Photo metadata.
		assert.Equal(t, "image", frm.PhotoType.Value)
		assert.Equal(t, false, frm.PhotoType.Mixed)
		assert.Equal(t, "Same Title", frm.PhotoTitle.Value)
		assert.Equal(t, false, frm.PhotoTitle.Mixed)
		assert.Equal(t, "", frm.PhotoCaption.Value)
		assert.Equal(t, false, frm.PhotoCaption.Mixed)
		assert.Equal(t, false, frm.PhotoFavorite.Value)
		assert.Equal(t, true, frm.PhotoFavorite.Mixed)
		assert.Equal(t, false, frm.PhotoPrivate.Value)
		assert.Equal(t, true, frm.PhotoPrivate.Mixed)
		assert.Equal(t, false, frm.PhotoScan.Value)
		assert.Equal(t, true, frm.PhotoScan.Mixed)
		assert.Equal(t, false, frm.PhotoPanorama.Value)
		assert.Equal(t, true, frm.PhotoPanorama.Mixed)
		assert.Equal(t, false, frm.CameraID.Mixed)
		assert.Equal(t, 1, frm.CameraID.Value)
		assert.Equal(t, false, frm.CameraID.Mixed)
		assert.Equal(t, 2, frm.LensID.Value)
		assert.Equal(t, false, frm.LensID.Mixed)
		assert.Equal(t, "", frm.PhotoCountry.Value)
		assert.Equal(t, true, frm.PhotoCountry.Mixed)
		assert.Equal(t, 0, frm.PhotoIso.Value)
		assert.Equal(t, false, frm.PhotoIso.Mixed)
		assert.Equal(t, float32(0), frm.PhotoFNumber.Value)
		assert.Equal(t, false, frm.PhotoFNumber.Mixed)
		assert.Equal(t, 0, frm.PhotoFocalLength.Value)
		assert.Equal(t, false, frm.PhotoFocalLength.Mixed)
	})
	t.Run("Success", func(t *testing.T) {
		photo1 := search.Photo{
			ID:           111115411,
			TakenSrc:     "",
			TimeZone:     "",
			PhotoUID:     "ps6sg6be2lvl0x41",
			PhotoType:    "image",
			PhotoTitle:   "Same Title",
			PhotoCountry: "de",
			CameraID:     1000001,
			LensID:       1000001,
		}

		photo2 := search.Photo{
			ID:           111115511,
			CreatedAt:    time.Time{},
			TakenAt:      time.Time{},
			TakenAtLocal: time.Time{},
			TakenSrc:     "",
			TimeZone:     "",
			PhotoUID:     "ps6sg6be2lvlx986",
			PhotoType:    "image",
			PhotoTitle:   "Same Title",
			PhotoCountry: "",
			CameraID:     1000000,
			LensID:       1000000,
		}

		photos := search.PhotoResults{photo1, photo2}

		frm := NewPhotosForm(photos)

		// Photo metadata.
		assert.Equal(t, "image", frm.PhotoType.Value)
		assert.Equal(t, false, frm.PhotoType.Mixed)
		assert.Equal(t, "Same Title", frm.PhotoTitle.Value)
		assert.Equal(t, false, frm.PhotoTitle.Mixed)
		assert.Equal(t, "", frm.PhotoCaption.Value)
		assert.Equal(t, false, frm.PhotoCaption.Mixed)
		assert.Equal(t, false, frm.PhotoFavorite.Value)
		assert.Equal(t, false, frm.PhotoFavorite.Mixed)
		assert.Equal(t, false, frm.PhotoPrivate.Value)
		assert.Equal(t, false, frm.PhotoPrivate.Mixed)
		assert.Equal(t, false, frm.PhotoScan.Value)
		assert.Equal(t, false, frm.PhotoScan.Mixed)
		assert.Equal(t, false, frm.PhotoPanorama.Value)
		assert.Equal(t, false, frm.PhotoPanorama.Mixed)
		assert.Equal(t, true, frm.CameraID.Mixed)
		assert.Equal(t, -2, frm.CameraID.Value)
		assert.Equal(t, true, frm.LensID.Mixed)
		assert.Equal(t, -2, frm.LensID.Value)
		assert.Equal(t, true, frm.LensID.Mixed)
		assert.Equal(t, "", frm.PhotoCountry.Value)
		assert.Equal(t, true, frm.PhotoCountry.Mixed)
		assert.Equal(t, 0, frm.PhotoIso.Value)
		assert.Equal(t, false, frm.PhotoIso.Mixed)
		assert.Equal(t, float32(0), frm.PhotoFNumber.Value)
		assert.Equal(t, false, frm.PhotoFNumber.Mixed)
		assert.Equal(t, 0, frm.PhotoFocalLength.Value)
		assert.Equal(t, false, frm.PhotoFocalLength.Mixed)
	})

}
