package batch

import (
	"fmt"
	"time"

	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/form"
)

// ConvertToPhotoForm converts PhotosForm into a regular form.Photo while
// preserving unchanged fields and marking source fields with SrcBatch where applicable.
func ConvertToPhotoForm(photo *entity.Photo, v *PhotosForm) (*form.Photo, error) {
	if photo == nil || v == nil {
		return nil, fmt.Errorf("photo or batch values is nil")
	}

	// Start with a form created from the current photo
	photoForm, err := form.NewPhoto(photo)
	if err != nil {
		return nil, fmt.Errorf("failed to create form from photo: %w", err)
	}

	switch v.PhotoTitle.Action {
	case ActionUpdate:
		photoForm.PhotoTitle = v.PhotoTitle.Value
		photoForm.TitleSrc = entity.SrcBatch
	case ActionRemove:
		photoForm.PhotoTitle = ""
		photoForm.TitleSrc = entity.SrcBatch
	}

	switch v.PhotoCaption.Action {
	case ActionUpdate:
		photoForm.PhotoCaption = v.PhotoCaption.Value
		photoForm.CaptionSrc = entity.SrcBatch
	case ActionRemove:
		photoForm.PhotoCaption = ""
		photoForm.CaptionSrc = entity.SrcBatch
	}

	if v.PhotoType.Action == ActionUpdate {
		photoForm.PhotoType = v.PhotoType.Value
		photoForm.TypeSrc = entity.SrcBatch
	}

	// Date/time fields
	timeChanged := false
	if v.PhotoDay.Action == ActionUpdate {
		photoForm.PhotoDay = v.PhotoDay.Value
		timeChanged = true
	}
	if v.PhotoMonth.Action == ActionUpdate {
		photoForm.PhotoMonth = v.PhotoMonth.Value
		timeChanged = true
	}
	if v.PhotoYear.Action == ActionUpdate {
		photoForm.PhotoYear = v.PhotoYear.Value
		timeChanged = true
	}
	if v.TimeZone.Action == ActionUpdate {
		photoForm.TimeZone = v.TimeZone.Value
		timeChanged = true
	}
	if timeChanged {
		photoForm.TakenSrc = entity.SrcBatch
	}

	if v.PhotoDay.Action == ActionUpdate && v.PhotoDay.Value == -1 {
		photoForm.PhotoYear = -1
	}
	if v.PhotoMonth.Action == ActionUpdate && v.PhotoMonth.Value == -1 {
		photoForm.PhotoYear = -1
	}

	// If any date field changed, recompute TakenAtLocal per photo and clamp invalid day.
	if v.PhotoDay.Action == ActionUpdate || v.PhotoMonth.Action == ActionUpdate || v.PhotoYear.Action == ActionUpdate {
		base := photo.TakenAtLocal

		// Determine target year: apply updated value if > 0; keep base if unknown or unchanged.
		year := base.Year()
		if v.PhotoYear.Action == ActionUpdate {
			if v.PhotoYear.Value > 0 {
				year = v.PhotoYear.Value
			}
		}

		// Determine target month: apply updated value if > 0; keep base if unknown or unchanged.
		month := int(base.Month())
		if v.PhotoMonth.Action == ActionUpdate {
			if v.PhotoMonth.Value > 0 {
				month = v.PhotoMonth.Value
			}
		}

		// Determine target day: -1 becomes 1; otherwise apply updated value; if unchanged keep base day.
		day := base.Day()
		if v.PhotoDay.Action == ActionUpdate {
			if v.PhotoDay.Value == -1 {
				day = 1
			} else if v.PhotoDay.Value > 0 {
				day = v.PhotoDay.Value
			}
		}

		// Clamp day to last valid day of month/year.
		lastDay := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
		if day > lastDay {
			day = lastDay
		}

		// Preserve time-of-day from base, construct new local date.
		newLocal := time.Date(year, time.Month(month), day, base.Hour(), base.Minute(), base.Second(), 0, time.UTC)
		photoForm.TakenAtLocal = newLocal
		photoForm.TakenSrc = entity.SrcBatch

		// Ensure PhotoDay is consistent with the clamped value when user provided a positive day.
		if v.PhotoDay.Action == ActionUpdate && v.PhotoDay.Value > 0 {
			photoForm.PhotoDay = day
		}

		// If only Month and/or Year changed (Day not explicitly updated),
		// keep PhotoDay consistent with the computed/clamped value.
		if v.PhotoDay.Action != ActionUpdate {
			photoForm.PhotoDay = day
		}
	}

	// Location fields
	locationChanged := false
	if v.PhotoLat.Action == ActionUpdate {
		photoForm.PhotoLat = v.PhotoLat.Value
		locationChanged = true
	}
	if v.PhotoLng.Action == ActionUpdate {
		photoForm.PhotoLng = v.PhotoLng.Value
		locationChanged = true
	}
	if v.PhotoCountry.Action == ActionUpdate {
		photoForm.PhotoCountry = v.PhotoCountry.Value
		locationChanged = true
	}
	if v.PhotoAltitude.Action == ActionUpdate {
		photoForm.PhotoAltitude = v.PhotoAltitude.Value
		locationChanged = true
	}
	if locationChanged {
		photoForm.PlaceSrc = entity.SrcBatch
	}

	// Boolean flags
	if v.PhotoFavorite.Action == ActionUpdate {
		photoForm.PhotoFavorite = v.PhotoFavorite.Value
	}
	if v.PhotoPrivate.Action == ActionUpdate {
		photoForm.PhotoPrivate = v.PhotoPrivate.Value
	}
	if v.PhotoScan.Action == ActionUpdate {
		photoForm.PhotoScan = v.PhotoScan.Value
	}
	if v.PhotoPanorama.Action == ActionUpdate {
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
	case ActionUpdate:
		photoForm.Details.Subject = v.DetailsSubject.Value
		photoForm.Details.SubjectSrc = entity.SrcBatch
	case ActionRemove:
		photoForm.Details.Subject = ""
		photoForm.Details.SubjectSrc = entity.SrcBatch
	}

	switch v.DetailsArtist.Action {
	case ActionUpdate:
		photoForm.Details.Artist = v.DetailsArtist.Value
		photoForm.Details.ArtistSrc = entity.SrcBatch
	case ActionRemove:
		photoForm.Details.Artist = ""
		photoForm.Details.ArtistSrc = entity.SrcBatch
	}

	switch v.DetailsCopyright.Action {
	case ActionUpdate:
		photoForm.Details.Copyright = v.DetailsCopyright.Value
		photoForm.Details.CopyrightSrc = entity.SrcBatch
	case ActionRemove:
		photoForm.Details.Copyright = ""
		photoForm.Details.CopyrightSrc = entity.SrcBatch
	}

	switch v.DetailsLicense.Action {
	case ActionUpdate:
		photoForm.Details.License = v.DetailsLicense.Value
		photoForm.Details.LicenseSrc = entity.SrcBatch
	case ActionRemove:
		photoForm.Details.License = ""
		photoForm.Details.LicenseSrc = entity.SrcBatch
	}

	// Set the PhotoID for details
	photoForm.Details.PhotoID = photo.ID

	return &photoForm, nil
}
