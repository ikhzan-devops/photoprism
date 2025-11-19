package batch

import (
	"fmt"
	"strings"
	"time"

	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/entity/query"
	"github.com/photoprism/photoprism/internal/entity/search"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/pkg/txt"
)

// PhotoSaveRequest bundles the context required to persist a batch update for a single photo.
type PhotoSaveRequest struct {
	Photo  *entity.Photo
	Form   *form.Photo
	Values *PhotosForm
}

// SaveBatchResult describes the outcome of PrepareAndSavePhotos so callers can publish events
// and build responses without re-walking the selection.
type SaveBatchResult struct {
	Requests     []*PhotoSaveRequest
	Results      []bool
	Preloaded    map[string]*entity.Photo
	UpdatedCount int
	SavedAny     bool
	Stats        MutationStats
}

// MutationStats captures how many photos received supporting mutations and
// how many errors occurred while applying them.
type MutationStats struct {
	AlbumMutations int
	LabelMutations int
	AlbumErrors    int
	LabelErrors    int
}

// NewPhotoSaveRequest converts the batch values into a form.Photo and bundles it with the
// target entity so callers outside this package do not have to depend on ConvertToPhotoForm.
func NewPhotoSaveRequest(photo *entity.Photo, values *PhotosForm) (*PhotoSaveRequest, error) {
	if values == nil {
		return nil, fmt.Errorf("batch: values are required")
	}

	frm, err := ConvertToPhotoForm(photo, values)
	if err != nil {
		return nil, err
	}

	return &PhotoSaveRequest{Photo: photo, Form: frm, Values: values}, nil
}

// PreparePhotoSaveRequests converts the given selection into save requests, ensuring each photo
// entity is hydrated and album/label actions are applied once per selection. The returned map is
// the (possibly newly populated) preloaded set so callers can reuse it for response rendering.
func PreparePhotoSaveRequests(photos search.PhotoResults, preloaded map[string]*entity.Photo, values *PhotosForm) ([]*PhotoSaveRequest, map[string]*entity.Photo, MutationStats) {
	if values == nil {
		return nil, preloaded, MutationStats{}
	}

	if preloaded == nil {
		preloaded = map[string]*entity.Photo{}
	}

	log.Infof("batch: updating %d photos", len(photos))

	saveRequests := make([]*PhotoSaveRequest, 0, len(photos))
	stats := MutationStats{}

	for _, result := range photos {
		photoID := result.PhotoUID

		if photoID == "" {
			continue
		}

		fullPhoto := preloaded[photoID]

		if fullPhoto == nil {
			loaded, err := query.PhotoPreloadByUID(photoID)
			if err != nil {
				log.Errorf("batch: failed to load photo %s: %s", photoID, err)
				continue
			}
			fullPhoto = &loaded
			preloaded[photoID] = fullPhoto
		}

		saveReq, err := NewPhotoSaveRequest(fullPhoto, values)

		if err != nil {
			log.Errorf("batch: failed to build save request for photo %s: %s", photoID, err)
			continue
		}

		saveRequests = append(saveRequests, saveReq)

		if values.Albums.Action == ActionUpdate {
			if errs := ApplyAlbums(fullPhoto, values.Albums); errs != nil {
				log.Errorf("batch: failed to update albums for photo %s: (%s)", photoID, errs)
				stats.AlbumErrors += len(errs)
			}

			if len(values.Albums.Items) > 0 {
				stats.AlbumMutations++
			}
		}

		if values.Labels.Action == ActionUpdate {
			if errs := ApplyLabels(fullPhoto, values.Labels); errs != nil {
				log.Errorf("batch: failed to update labels for photo %s (%s)", photoID, errs)
				stats.LabelErrors += len(errs)
			}

			if len(values.Labels.Items) > 0 {
				stats.LabelMutations++
			}
		}
	}

	return saveRequests, preloaded, stats
}

// PrepareAndSavePhotos hydrates the photo selection, applies album/label mutations, builds save
// requests, and persists the changes in one step. It returns a SaveBatchResult so callers can run
// follow-up work (events, cache flushes) without re-querying state.
func PrepareAndSavePhotos(photos search.PhotoResults, preloaded map[string]*entity.Photo, values *PhotosForm) (*SaveBatchResult, error) {
	start := time.Now()
	result := &SaveBatchResult{Preloaded: preloaded}

	if values == nil {
		if result.Preloaded == nil {
			result.Preloaded = map[string]*entity.Photo{}
		}
		return result, nil
	}

	requests, preloaded, mutationStats := PreparePhotoSaveRequests(photos, result.Preloaded, values)

	result.Requests = requests
	result.Preloaded = preloaded
	result.Stats = mutationStats

	if len(requests) == 0 {
		return result, nil
	}

	saveResults, err := SavePhotos(requests)

	if err != nil {
		return nil, err
	}

	result.Results = saveResults

	for i, saved := range saveResults {
		if saved {
			result.UpdatedCount++
			result.SavedAny = true
			log.Debugf("batch: successfully updated photo %s", requests[i].Photo.PhotoUID)
		}
	}

	var logFields []string

	if result.UpdatedCount > 0 {
		logFields = append(logFields, fmt.Sprintf("metadata (%d/%d)", result.UpdatedCount, len(photos)))
	}

	if result.Stats.LabelMutations > 0 {
		entry := fmt.Sprintf("labels (%d/%d)", result.Stats.LabelMutations, len(photos))
		if result.Stats.LabelErrors > 0 {
			entry = fmt.Sprintf("%s with %d errors", entry, result.Stats.LabelErrors)
		}
		logFields = append(logFields, entry)
	}

	if result.Stats.AlbumMutations > 0 {
		entry := fmt.Sprintf("albums (%d/%d)", result.Stats.AlbumMutations, len(photos))
		if result.Stats.AlbumErrors > 0 {
			entry = fmt.Sprintf("%s with %d errors", entry, result.Stats.AlbumErrors)
		}
		logFields = append(logFields, entry)
	}

	if len(logFields) > 0 {
		log.Infof("batch: updated photo %s [%s]", txt.JoinAnd(logFields), time.Since(start))
	} else {
		log.Infof("batch: no photos have been updated [%s]", time.Since(start))
	}

	return result, nil
}

// SavePhotos persists the batch updates described by the provided requests while skipping the
// heavyweight per-photo maintenance performed by entity.SavePhotoForm. It only updates the
// columns that actually changed, flags the photo for background metadata refresh by clearing
// CheckedAt, and updates the shared counts once per batch instead of once per photo.
func SavePhotos(requests []*PhotoSaveRequest) ([]bool, error) {
	results := make([]bool, len(requests))
	anySaved := false

	for i, req := range requests {
		saved, err := savePhoto(req)

		if err != nil {
			return results, err
		}

		results[i] = saved

		if saved {
			anySaved = true
		}
	}

	if anySaved {
		entity.UpdateCountsAsync()
	}

	return results, nil
}

// savePhoto applies the batched form values to a single entity.Photo and writes
// only the changed columns back to the database.
func savePhoto(req *PhotoSaveRequest) (bool, error) {
	if req == nil || req.Photo == nil || req.Form == nil || req.Values == nil {
		return false, fmt.Errorf("batch: invalid save request")
	}

	p := req.Photo
	formValues := req.Values
	formData := req.Form

	if !p.HasID() {
		return false, fmt.Errorf("batch: photo has no id")
	}

	original := *p
	details := p.GetDetails()
	origDetails := *details

	if shouldUpdateString(formValues.PhotoTitle) {
		p.PhotoTitle = formData.PhotoTitle
		p.TitleSrc = formData.TitleSrc
	}

	if shouldUpdateString(formValues.PhotoCaption) {
		p.PhotoCaption = formData.PhotoCaption
		p.CaptionSrc = formData.CaptionSrc
	}

	if formValues.PhotoType.Action == ActionUpdate {
		p.PhotoType = formData.PhotoType
		p.TypeSrc = formData.TypeSrc
	}

	if formValues.PhotoFavorite.Action == ActionUpdate {
		p.PhotoFavorite = formValues.PhotoFavorite.Value
	}

	if formValues.PhotoPrivate.Action == ActionUpdate {
		p.PhotoPrivate = formValues.PhotoPrivate.Value
	}

	if formValues.PhotoScan.Action == ActionUpdate {
		p.PhotoScan = formValues.PhotoScan.Value
	}

	if formValues.PhotoPanorama.Action == ActionUpdate {
		p.PhotoPanorama = formValues.PhotoPanorama.Value
	}

	timeChanged := formValues.PhotoDay.Action == ActionUpdate ||
		formValues.PhotoMonth.Action == ActionUpdate ||
		formValues.PhotoYear.Action == ActionUpdate ||
		formValues.TimeZone.Action == ActionUpdate

	if timeChanged {
		p.PhotoYear = formData.PhotoYear
		p.PhotoMonth = formData.PhotoMonth
		p.PhotoDay = formData.PhotoDay
		if formValues.TimeZone.Action == ActionUpdate {
			p.TimeZone = formData.TimeZone
		}
		p.TakenAtLocal = formData.TakenAtLocal
		p.TakenSrc = formData.TakenSrc

		p.NormalizeValues()

		if p.TimeZoneUTC() {
			p.TakenAt = p.TakenAtLocal
		} else {
			p.TakenAt = p.GetTakenAt()
		}

		p.UpdateDateFields()
	}

	locationChanged := formValues.PhotoLat.Action == ActionUpdate ||
		formValues.PhotoLng.Action == ActionUpdate ||
		formValues.PhotoCountry.Action == ActionUpdate ||
		formValues.PhotoAltitude.Action == ActionUpdate

	if formValues.PhotoLat.Action == ActionUpdate {
		p.PhotoLat = formValues.PhotoLat.Value
	}

	if formValues.PhotoLng.Action == ActionUpdate {
		p.PhotoLng = formValues.PhotoLng.Value
	}

	if formValues.PhotoAltitude.Action == ActionUpdate {
		p.PhotoAltitude = formValues.PhotoAltitude.Value
	}

	if formValues.PhotoCountry.Action == ActionUpdate {
		p.PhotoCountry = formValues.PhotoCountry.Value
	}

	if locationChanged {
		p.PlaceSrc = entity.SrcBatch
		locKeywords, locLabels := p.UpdateLocation()
		if len(locLabels) > 0 {
			p.AddLabels(locLabels)
		}
		if len(locKeywords) > 0 {
			words := txt.UniqueWords(txt.Words(details.Keywords))
			words = append(words, locKeywords...)
			details.Keywords = strings.Join(txt.UniqueWords(words), ", ")
		}
	}

	if formValues.DetailsSubject.Action == ActionUpdate || formValues.DetailsSubject.Action == ActionRemove {
		details.Subject = formData.Details.Subject
		details.SubjectSrc = formData.Details.SubjectSrc
	}

	if formValues.DetailsArtist.Action == ActionUpdate || formValues.DetailsArtist.Action == ActionRemove {
		details.Artist = formData.Details.Artist
		details.ArtistSrc = formData.Details.ArtistSrc
	}

	if formValues.DetailsCopyright.Action == ActionUpdate || formValues.DetailsCopyright.Action == ActionRemove {
		details.Copyright = formData.Details.Copyright
		details.CopyrightSrc = formData.Details.CopyrightSrc
	}

	if formValues.DetailsLicense.Action == ActionUpdate || formValues.DetailsLicense.Action == ActionRemove {
		details.License = formData.Details.License
		details.LicenseSrc = formData.Details.LicenseSrc
	}

	updates := entity.Values{}
	addUpdate := func(column string, changed bool, value interface{}) {
		if changed {
			updates[column] = value
		}
	}

	addUpdate("photo_title", p.PhotoTitle != original.PhotoTitle, p.PhotoTitle)
	addUpdate("title_src", p.TitleSrc != original.TitleSrc, p.TitleSrc)
	addUpdate("photo_caption", p.PhotoCaption != original.PhotoCaption, p.PhotoCaption)
	addUpdate("caption_src", p.CaptionSrc != original.CaptionSrc, p.CaptionSrc)
	addUpdate("photo_type", p.PhotoType != original.PhotoType, p.PhotoType)
	addUpdate("type_src", p.TypeSrc != original.TypeSrc, p.TypeSrc)
	addUpdate("photo_favorite", p.PhotoFavorite != original.PhotoFavorite, p.PhotoFavorite)
	addUpdate("photo_private", p.PhotoPrivate != original.PhotoPrivate, p.PhotoPrivate)
	addUpdate("photo_scan", p.PhotoScan != original.PhotoScan, p.PhotoScan)
	addUpdate("photo_panorama", p.PhotoPanorama != original.PhotoPanorama, p.PhotoPanorama)
	addUpdate("photo_year", p.PhotoYear != original.PhotoYear, p.PhotoYear)
	addUpdate("photo_month", p.PhotoMonth != original.PhotoMonth, p.PhotoMonth)
	addUpdate("photo_day", p.PhotoDay != original.PhotoDay, p.PhotoDay)
	addUpdate("time_zone", p.TimeZone != original.TimeZone, p.TimeZone)
	addUpdate("taken_src", p.TakenSrc != original.TakenSrc, p.TakenSrc)
	addUpdate("taken_at", !p.TakenAt.Equal(original.TakenAt), p.TakenAt)
	addUpdate("taken_at_local", !p.TakenAtLocal.Equal(original.TakenAtLocal), p.TakenAtLocal)
	addUpdate("photo_lat", p.PhotoLat != original.PhotoLat, p.PhotoLat)
	addUpdate("photo_lng", p.PhotoLng != original.PhotoLng, p.PhotoLng)
	addUpdate("photo_altitude", p.PhotoAltitude != original.PhotoAltitude, p.PhotoAltitude)
	addUpdate("photo_country", p.PhotoCountry != original.PhotoCountry, p.PhotoCountry)
	addUpdate("place_id", p.PlaceID != original.PlaceID, p.PlaceID)
	addUpdate("cell_id", p.CellID != original.CellID, p.CellID)
	addUpdate("place_src", p.PlaceSrc != original.PlaceSrc, p.PlaceSrc)

	detailUpdates := entity.Values{}
	if details.Subject != origDetails.Subject {
		detailUpdates["subject"] = details.Subject
	}
	if details.SubjectSrc != origDetails.SubjectSrc {
		detailUpdates["subject_src"] = details.SubjectSrc
	}
	if details.Artist != origDetails.Artist {
		detailUpdates["artist"] = details.Artist
	}
	if details.ArtistSrc != origDetails.ArtistSrc {
		detailUpdates["artist_src"] = details.ArtistSrc
	}
	if details.Copyright != origDetails.Copyright {
		detailUpdates["copyright"] = details.Copyright
	}
	if details.CopyrightSrc != origDetails.CopyrightSrc {
		detailUpdates["copyright_src"] = details.CopyrightSrc
	}
	if details.License != origDetails.License {
		detailUpdates["license"] = details.License
	}
	if details.LicenseSrc != origDetails.LicenseSrc {
		detailUpdates["license_src"] = details.LicenseSrc
	}
	if details.Keywords != origDetails.Keywords {
		detailUpdates["keywords"] = details.Keywords
	}

	if len(updates) == 0 && len(detailUpdates) == 0 {
		return false, nil
	}

	log.Debugf("batch: saving photo %s with updates=%v details=%v", p.PhotoUID, updates, detailUpdates)

	edited := entity.Now()
	updates["edited_at"] = edited
	p.EditedAt = &edited
	updates["updated_at"] = edited
	p.UpdatedAt = edited
	updates["checked_at"] = nil
	p.CheckedAt = nil

	if err := p.Updates(updates); err != nil {
		return false, err
	}

	if len(detailUpdates) > 0 {
		if err := details.Updates(detailUpdates); err != nil {
			return false, err
		}
	}

	return true, nil
}

// shouldUpdateString reports whether the provided field requests an update or removal.
func shouldUpdateString(v String) bool {
	return v.Action == ActionUpdate || v.Action == ActionRemove
}
