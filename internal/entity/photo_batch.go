package entity

import (
	"fmt"
	"time"

	"github.com/photoprism/photoprism/internal/form"
)

type BatchAction string

const (
	BatchActionNone   BatchAction = "none"
	BatchActionUpdate BatchAction = "update"
	BatchActionRemove BatchAction = "remove"
)

// BatchValue wrappers replicate the shape of batch form values without creating a package cycle.
type BatchString struct {
	Value  string
	Mixed  bool
	Action BatchAction
}

type BatchBool struct {
	Value  bool
	Mixed  bool
	Action BatchAction
}

type BatchInt struct {
	Value  int
	Mixed  bool
	Action BatchAction
}

type BatchFloat64 struct {
	Value  float64
	Mixed  bool
	Action BatchAction
}

type BatchTime struct {
	Value  time.Time
	Mixed  bool
	Action BatchAction
}

// BatchPhotoValues contains the subset of batch values relevant for converting to a normal photo form.
type BatchPhotoValues struct {
	PhotoType    BatchString
	PhotoTitle   BatchString
	PhotoCaption BatchString

	TakenAt      BatchTime
	TakenAtLocal BatchTime
	TimeZone     BatchString
	PhotoYear    BatchInt
	PhotoMonth   BatchInt
	PhotoDay     BatchInt

	PhotoLat      BatchFloat64
	PhotoLng      BatchFloat64
	PhotoCountry  BatchString
	PhotoAltitude BatchInt

	PhotoFavorite BatchBool
	PhotoPrivate  BatchBool
	PhotoScan     BatchBool
	PhotoPanorama BatchBool

	DetailsSubject   BatchString
	DetailsArtist    BatchString
	DetailsCopyright BatchString
	DetailsLicense   BatchString
}

// ConvertBatchToPhotoForm converts entity.BatchPhotoValues into a regular form.Photo while
// preserving unchanged fields and marking source fields with SrcBatch where applicable.
func ConvertBatchToPhotoForm(photo *Photo, v *BatchPhotoValues) (*form.Photo, error) {
	if photo == nil || v == nil {
		return nil, fmt.Errorf("photo or batch values is nil")
	}

	// Start with a form created from the current photo
	photoForm, err := form.NewPhoto(photo)
	if err != nil {
		return nil, fmt.Errorf("failed to create form from photo: %w", err)
	}

	switch v.PhotoTitle.Action {
	case BatchActionUpdate:
		photoForm.PhotoTitle = v.PhotoTitle.Value
		photoForm.TitleSrc = SrcBatch
	case BatchActionRemove:
		photoForm.PhotoTitle = ""
		photoForm.TitleSrc = SrcBatch
	}

	switch v.PhotoCaption.Action {
	case BatchActionUpdate:
		photoForm.PhotoCaption = v.PhotoCaption.Value
		photoForm.CaptionSrc = SrcBatch
	case BatchActionRemove:
		photoForm.PhotoCaption = ""
		photoForm.CaptionSrc = SrcBatch
	}

	if v.PhotoType.Action == BatchActionUpdate {
		photoForm.PhotoType = v.PhotoType.Value
		photoForm.TypeSrc = SrcBatch
	}

	// Date/time fields
	timeChanged := false
	if v.PhotoDay.Action == BatchActionUpdate {
		photoForm.PhotoDay = v.PhotoDay.Value
		timeChanged = true
	}
	if v.PhotoMonth.Action == BatchActionUpdate {
		photoForm.PhotoMonth = v.PhotoMonth.Value
		timeChanged = true
	}
	if v.PhotoYear.Action == BatchActionUpdate {
		photoForm.PhotoYear = v.PhotoYear.Value
		timeChanged = true
	}
	if v.TimeZone.Action == BatchActionUpdate {
		photoForm.TimeZone = v.TimeZone.Value
		timeChanged = true
	}
	if timeChanged {
		photoForm.TakenSrc = SrcBatch
	}

	// Location fields
	locationChanged := false
	if v.PhotoLat.Action == BatchActionUpdate {
		photoForm.PhotoLat = v.PhotoLat.Value
		locationChanged = true
	}
	if v.PhotoLng.Action == BatchActionUpdate {
		photoForm.PhotoLng = v.PhotoLng.Value
		locationChanged = true
	}
	if v.PhotoCountry.Action == BatchActionUpdate {
		photoForm.PhotoCountry = v.PhotoCountry.Value
		locationChanged = true
	}
	if v.PhotoAltitude.Action == BatchActionUpdate {
		photoForm.PhotoAltitude = v.PhotoAltitude.Value
		locationChanged = true
	}
	if locationChanged {
		photoForm.PlaceSrc = SrcBatch
	}

	// Boolean flags
	if v.PhotoFavorite.Action == BatchActionUpdate {
		photoForm.PhotoFavorite = v.PhotoFavorite.Value
	}
	if v.PhotoPrivate.Action == BatchActionUpdate {
		photoForm.PhotoPrivate = v.PhotoPrivate.Value
	}
	if v.PhotoScan.Action == BatchActionUpdate {
		photoForm.PhotoScan = v.PhotoScan.Value
	}
	if v.PhotoPanorama.Action == BatchActionUpdate {
		photoForm.PhotoPanorama = v.PhotoPanorama.Value
	}

	// Details fields - preserve existing values, only update changed ones
	currentDetails := photo.GetDetails()
	if currentDetails != nil {
		// Start with current values to preserve unchanged fields
		photoForm.Details.Subject = currentDetails.Subject
		photoForm.Details.SubjectSrc = currentDetails.SubjectSrc
		photoForm.Details.Artist = currentDetails.Artist
		photoForm.Details.ArtistSrc = currentDetails.ArtistSrc
		photoForm.Details.Copyright = currentDetails.Copyright
		photoForm.Details.CopyrightSrc = currentDetails.CopyrightSrc
		photoForm.Details.License = currentDetails.License
		photoForm.Details.LicenseSrc = currentDetails.LicenseSrc
		photoForm.Details.Keywords = currentDetails.Keywords
		photoForm.Details.KeywordsSrc = currentDetails.KeywordsSrc
		photoForm.Details.Notes = currentDetails.Notes
		photoForm.Details.NotesSrc = currentDetails.NotesSrc
	}

	switch v.DetailsSubject.Action {
	case BatchActionUpdate:
		photoForm.Details.Subject = v.DetailsSubject.Value
		photoForm.Details.SubjectSrc = SrcBatch
	case BatchActionRemove:
		photoForm.Details.Subject = ""
		photoForm.Details.SubjectSrc = SrcBatch
	}

	switch v.DetailsArtist.Action {
	case BatchActionUpdate:
		photoForm.Details.Artist = v.DetailsArtist.Value
		photoForm.Details.ArtistSrc = SrcBatch
	case BatchActionRemove:
		photoForm.Details.Artist = ""
		photoForm.Details.ArtistSrc = SrcBatch
	}

	switch v.DetailsCopyright.Action {
	case BatchActionUpdate:
		photoForm.Details.Copyright = v.DetailsCopyright.Value
		photoForm.Details.CopyrightSrc = SrcBatch
	case BatchActionRemove:
		photoForm.Details.Copyright = ""
		photoForm.Details.CopyrightSrc = SrcBatch
	}

	switch v.DetailsLicense.Action {
	case BatchActionUpdate:
		photoForm.Details.License = v.DetailsLicense.Value
		photoForm.Details.LicenseSrc = SrcBatch
	case BatchActionRemove:
		photoForm.Details.License = ""
		photoForm.Details.LicenseSrc = SrcBatch
	}

	// Set the PhotoID for details
	photoForm.Details.PhotoID = photo.ID

	return &photoForm, nil
}
