package batch

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/entity"
)

func TestConvertToPhotoForm(t *testing.T) {
	t.Run("UpdateTitleCaptionTypeAndBooleans", func(t *testing.T) {
		photo := &entity.Photo{
			PhotoTitle:   "Old Title",
			PhotoCaption: "Old Caption",
			PhotoType:    "image",
		}

		v := &PhotosForm{}
		v.PhotoTitle.Action = ActionUpdate
		v.PhotoTitle.Value = "New Title"

		v.PhotoCaption.Action = ActionRemove // remove caption
		v.PhotoCaption.Value = "whatever"

		v.PhotoType.Action = ActionUpdate
		v.PhotoType.Value = "live"

		v.PhotoFavorite.Action = ActionUpdate
		v.PhotoFavorite.Value = true

		v.PhotoPrivate.Action = ActionUpdate
		v.PhotoPrivate.Value = true

		v.PhotoScan.Action = ActionUpdate
		v.PhotoScan.Value = true

		v.PhotoPanorama.Action = ActionUpdate
		v.PhotoPanorama.Value = true

		form, err := ConvertToPhotoForm(photo, v)
		assert.NoError(t, err)

		assert.Equal(t, "New Title", form.PhotoTitle)
		assert.Equal(t, entity.SrcBatch, form.TitleSrc)

		assert.Equal(t, "", form.PhotoCaption)
		assert.Equal(t, entity.SrcBatch, form.CaptionSrc)

		assert.Equal(t, "live", form.PhotoType)
		assert.Equal(t, entity.SrcBatch, form.TypeSrc)

		assert.True(t, form.PhotoFavorite)
		assert.True(t, form.PhotoPrivate)
		assert.True(t, form.PhotoScan)
		assert.True(t, form.PhotoPanorama)
	})

	t.Run("UpdateLocationSetsPlaceSrc", func(t *testing.T) {
		photo := &entity.Photo{}

		v := &PhotosForm{}
		v.PhotoLat.Action = ActionUpdate
		v.PhotoLat.Value = 37.5
		v.PhotoLng.Action = ActionUpdate
		v.PhotoLng.Value = -122.4
		v.PhotoCountry.Action = ActionUpdate
		v.PhotoCountry.Value = "us"
		v.PhotoAltitude.Action = ActionUpdate
		v.PhotoAltitude.Value = 15

		form, err := ConvertToPhotoForm(photo, v)
		assert.NoError(t, err)

		assert.Equal(t, 37.5, form.PhotoLat)
		assert.Equal(t, -122.4, form.PhotoLng)
		assert.Equal(t, "us", form.PhotoCountry)
		assert.Equal(t, 15, form.PhotoAltitude)
		assert.Equal(t, entity.SrcBatch, form.PlaceSrc)
	})

	t.Run("UpdateDetails", func(t *testing.T) {
		photo := &entity.Photo{}

		v := &PhotosForm{}
		v.DetailsSubject.Action = ActionUpdate
		v.DetailsSubject.Value = "Subject"
		v.DetailsArtist.Action = ActionUpdate
		v.DetailsArtist.Value = "Artist"
		v.DetailsCopyright.Action = ActionUpdate
		v.DetailsCopyright.Value = "Copyright"
		v.DetailsLicense.Action = ActionUpdate
		v.DetailsLicense.Value = "License"

		form, err := ConvertToPhotoForm(photo, v)
		assert.NoError(t, err)

		assert.Equal(t, "Subject", form.Details.Subject)
		assert.Equal(t, entity.SrcBatch, form.Details.SubjectSrc)
		assert.Equal(t, "Artist", form.Details.Artist)
		assert.Equal(t, entity.SrcBatch, form.Details.ArtistSrc)
		assert.Equal(t, "Copyright", form.Details.Copyright)
		assert.Equal(t, entity.SrcBatch, form.Details.CopyrightSrc)
		assert.Equal(t, "License", form.Details.License)
		assert.Equal(t, entity.SrcBatch, form.Details.LicenseSrc)
	})
}
