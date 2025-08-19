package api

import (
	"sort"

	"github.com/photoprism/photoprism/internal/entity/query"
	"github.com/photoprism/photoprism/internal/entity/search"
	"github.com/photoprism/photoprism/internal/form/batch"
)

// NewPhotosForm returns a new batch edit form instance
// initialized with values from the selected photos.
func NewPhotosForm(photos search.PhotoResults) *batch.PhotosForm {
	// Create a new batch edit form and initialize it
	// with the values from the selected photos.
	frm := &batch.PhotosForm{
		Albums: batch.Items{Items: []batch.Item{}, Mixed: false, Action: batch.ActionNone},
		Labels: batch.Items{Items: []batch.Item{}, Mixed: false, Action: batch.ActionNone},
	}

	// Populate Albums and Labels from selected photos (no raw SQL; use preload helpers)
	total := len(photos)
	if total > 0 {
		type albumAgg struct {
			title string
			cnt   int
		}
		type labelAgg struct {
			name string
			cnt  int
		}
		albumCount := map[string]albumAgg{}
		labelCount := map[string]labelAgg{}

		for _, sp := range photos {
			if sp.PhotoUID == "" {
				continue
			}
			p, err := query.PhotoPreloadByUID(sp.PhotoUID)
			if err != nil || !p.HasID() {
				continue
			}
			// Albums on this photo
			for _, a := range p.Albums {
				if a.AlbumUID == "" || a.Deleted() {
					continue
				}
				v := albumCount[a.AlbumUID]
				v.title = a.AlbumTitle
				v.cnt++
				albumCount[a.AlbumUID] = v
			}
			// Labels on this photo (only visible ones: uncertainty < 100)
			for _, pl := range p.Labels {
				if pl.Uncertainty >= 100 || pl.Label == nil || !pl.Label.HasID() {
					continue
				}
				uid := pl.Label.LabelUID
				if uid == "" {
					continue
				}
				v := labelCount[uid]
				v.name = pl.Label.LabelName
				v.cnt++
				labelCount[uid] = v
			}
		}

		// Build Albums items
		frm.Albums.Items = make([]batch.Item, 0, len(albumCount))
		anyAlbumMixed := false
		for uid, agg := range albumCount {
			mixed := agg.cnt > 0 && agg.cnt < total
			if mixed {
				anyAlbumMixed = true
			}
			frm.Albums.Items = append(frm.Albums.Items, batch.Item{Value: uid, Title: agg.title, Mixed: mixed, Action: batch.ActionNone})
		}
		// Sort shared-first (Mixed=false), then by Title alphabetically
		sort.Slice(frm.Albums.Items, func(i, j int) bool {
			if frm.Albums.Items[i].Mixed != frm.Albums.Items[j].Mixed {
				return !frm.Albums.Items[i].Mixed && frm.Albums.Items[j].Mixed
			}
			return frm.Albums.Items[i].Title < frm.Albums.Items[j].Title
		})
		frm.Albums.Mixed = anyAlbumMixed
		frm.Albums.Action = batch.ActionNone

		// Build Labels items
		frm.Labels.Items = make([]batch.Item, 0, len(labelCount))
		anyLabelMixed := false
		for uid, agg := range labelCount {
			mixed := agg.cnt > 0 && agg.cnt < total
			if mixed {
				anyLabelMixed = true
			}
			frm.Labels.Items = append(frm.Labels.Items, batch.Item{Value: uid, Title: agg.name, Mixed: mixed, Action: batch.ActionNone})
		}
		// Sort shared-first (Mixed=false), then by Title alphabetically
		sort.Slice(frm.Labels.Items, func(i, j int) bool {
			if frm.Labels.Items[i].Mixed != frm.Labels.Items[j].Mixed {
				return !frm.Labels.Items[i].Mixed && frm.Labels.Items[j].Mixed
			}
			return frm.Labels.Items[i].Title < frm.Labels.Items[j].Title
		})
		frm.Labels.Mixed = anyLabelMixed
		frm.Labels.Action = batch.ActionNone
	}

	// TODO: Verify that all required PhotosForm values are present and
	//       properly initialized or use in the frontend form component.
	for i, photo := range photos {
		if i == 0 {
			frm.PhotoType.Value = photo.PhotoType
			frm.PhotoType.Action = batch.ActionNone
		} else if photo.PhotoType != frm.PhotoType.Value {
			frm.PhotoType.Mixed = true
			frm.PhotoType.Value = ""
		}

		if i == 0 {
			frm.PhotoTitle.Value = photo.PhotoTitle
			frm.PhotoTitle.Action = batch.ActionNone
		} else if photo.PhotoTitle != frm.PhotoTitle.Value {
			frm.PhotoTitle.Mixed = true
			frm.PhotoTitle.Value = ""
		}

		if i == 0 {
			frm.PhotoCaption.Value = photo.PhotoCaption
			frm.PhotoCaption.Action = batch.ActionNone
		} else if photo.PhotoCaption != frm.PhotoCaption.Value {
			frm.PhotoCaption.Mixed = true
			frm.PhotoCaption.Value = ""
		}

		if i == 0 {
			frm.TakenAt.Value = photo.TakenAt
			frm.TakenAt.Action = batch.ActionNone
		} else if photo.TakenAt != frm.TakenAt.Value {
			frm.TakenAt.Mixed = true
		}

		if i == 0 {
			frm.TakenAtLocal.Value = photo.TakenAtLocal
			frm.TakenAtLocal.Action = batch.ActionNone
		} else if photo.TakenAtLocal != frm.TakenAtLocal.Value {
			frm.TakenAtLocal.Mixed = true
		}

		if i == 0 {
			frm.PhotoDay.Value = photo.PhotoDay
			frm.PhotoDay.Action = batch.ActionNone
		} else if photo.PhotoDay != frm.PhotoDay.Value {
			frm.PhotoDay.Mixed = true
			frm.PhotoDay.Value = -2
		}

		if i == 0 {
			frm.PhotoMonth.Value = photo.PhotoMonth
			frm.PhotoMonth.Action = batch.ActionNone
		} else if photo.PhotoMonth != frm.PhotoMonth.Value {
			frm.PhotoMonth.Mixed = true
			frm.PhotoMonth.Value = -2
		}

		if i == 0 {
			frm.PhotoYear.Value = photo.PhotoYear
			frm.PhotoYear.Action = batch.ActionNone
		} else if photo.PhotoYear != frm.PhotoYear.Value {
			frm.PhotoYear.Mixed = true
			frm.PhotoYear.Value = -2
		}

		if i == 0 {
			frm.TimeZone.Value = photo.TimeZone
			frm.TimeZone.Action = batch.ActionNone
		} else if photo.TimeZone != frm.TimeZone.Value {
			frm.TimeZone.Mixed = true
			frm.TimeZone.Value = ""
		}

		if i == 0 {
			frm.PhotoCountry.Value = photo.PhotoCountry
			frm.PhotoCountry.Action = batch.ActionNone
		} else if photo.PhotoCountry != frm.PhotoCountry.Value {
			frm.PhotoCountry.Mixed = true
			frm.PhotoCountry.Value = ""
		}

		if i == 0 {
			frm.PhotoAltitude.Value = photo.PhotoAltitude
			frm.PhotoAltitude.Action = batch.ActionNone
		} else if photo.PhotoAltitude != frm.PhotoAltitude.Value {
			frm.PhotoAltitude.Mixed = true
			frm.PhotoAltitude.Value = 0
		}

		if i == 0 {
			frm.PhotoLat.Value = photo.PhotoLat
			frm.PhotoLat.Action = batch.ActionNone
		} else if photo.PhotoLat != frm.PhotoLat.Value {
			frm.PhotoLat.Mixed = true
			frm.PhotoLat.Value = 0.0
		}

		if i == 0 {
			frm.PhotoLng.Value = photo.PhotoLng
			frm.PhotoLng.Action = batch.ActionNone
		} else if photo.PhotoLng != frm.PhotoLng.Value {
			frm.PhotoLng.Mixed = true
			frm.PhotoLng.Value = 0.0
		}

		if i == 0 {
			frm.PhotoIso.Value = photo.PhotoIso
			frm.PhotoIso.Action = batch.ActionNone
		} else if photo.PhotoIso != frm.PhotoIso.Value {
			frm.PhotoIso.Mixed = true
			frm.PhotoIso.Value = 0
		}

		if i == 0 {
			frm.PhotoFocalLength.Value = photo.PhotoFocalLength
			frm.PhotoFocalLength.Action = batch.ActionNone
		} else if photo.PhotoFocalLength != frm.PhotoFocalLength.Value {
			frm.PhotoFocalLength.Mixed = true
			frm.PhotoFocalLength.Value = 0
		}

		if i == 0 {
			frm.PhotoFNumber.Value = photo.PhotoFNumber
			frm.PhotoFNumber.Action = batch.ActionNone
		} else if photo.PhotoFNumber != frm.PhotoFNumber.Value {
			frm.PhotoFNumber.Mixed = true
			frm.PhotoFNumber.Value = 0
		}

		if i == 0 {
			frm.PhotoExposure.Value = photo.PhotoExposure
			frm.PhotoExposure.Action = batch.ActionNone
		} else if photo.PhotoExposure != frm.PhotoExposure.Value {
			frm.PhotoExposure.Mixed = true
			frm.PhotoExposure.Value = ""
		}

		if i == 0 {
			frm.PhotoFavorite.Value = photo.PhotoFavorite
			frm.PhotoFavorite.Action = batch.ActionNone
		} else if photo.PhotoFavorite != frm.PhotoFavorite.Value {
			frm.PhotoFavorite.Mixed = true
			frm.PhotoFavorite.Value = false
		}

		if i == 0 {
			frm.PhotoPrivate.Value = photo.PhotoPrivate
			frm.PhotoPrivate.Action = batch.ActionNone
		} else if photo.PhotoPrivate != frm.PhotoPrivate.Value {
			frm.PhotoPrivate.Mixed = true
			frm.PhotoPrivate.Value = false
		}

		if i == 0 {
			frm.PhotoScan.Value = photo.PhotoScan
			frm.PhotoScan.Action = batch.ActionNone
		} else if photo.PhotoScan != frm.PhotoScan.Value {
			frm.PhotoScan.Mixed = true
			frm.PhotoScan.Value = false
		}

		if i == 0 {
			frm.PhotoPanorama.Value = photo.PhotoPanorama
			frm.PhotoPanorama.Action = batch.ActionNone
		} else if photo.PhotoPanorama != frm.PhotoPanorama.Value {
			frm.PhotoPanorama.Mixed = true
			frm.PhotoPanorama.Value = false
		}

		if i == 0 {
			frm.CameraID.Value = int(photo.CameraID)
			frm.CameraID.Action = batch.ActionNone
		} else if photo.CameraID != uint(frm.CameraID.Value) {
			frm.CameraID.Mixed = true
			frm.CameraID.Value = -2
		}

		if i == 0 {
			frm.LensID.Value = int(photo.LensID)
			frm.LensID.Action = batch.ActionNone
		} else if photo.LensID != uint(frm.LensID.Value) {
			frm.LensID.Mixed = true
			frm.LensID.Value = -2
		}

		if i == 0 {
			frm.DetailsKeywords.Value = photo.DetailsKeywords
			frm.DetailsKeywords.Action = batch.ActionNone
		} else if photo.DetailsKeywords != frm.DetailsKeywords.Value {
			frm.DetailsKeywords.Mixed = true
			frm.DetailsKeywords.Value = ""
		}

		if i == 0 {
			frm.DetailsSubject.Value = photo.DetailsSubject
			frm.DetailsSubject.Action = batch.ActionNone
		} else if photo.DetailsSubject != frm.DetailsSubject.Value {
			frm.DetailsSubject.Mixed = true
			frm.DetailsSubject.Value = ""
		}

		if i == 0 {
			frm.DetailsArtist.Value = photo.DetailsArtist
			frm.DetailsArtist.Action = batch.ActionNone
		} else if photo.DetailsArtist != frm.DetailsArtist.Value {
			frm.DetailsArtist.Mixed = true
			frm.DetailsArtist.Value = ""
		}

		if i == 0 {
			frm.DetailsCopyright.Value = photo.DetailsCopyright
			frm.DetailsCopyright.Action = batch.ActionNone
		} else if photo.DetailsCopyright != frm.DetailsCopyright.Value {
			frm.DetailsCopyright.Mixed = true
			frm.DetailsCopyright.Value = ""
		}

		if i == 0 {
			frm.DetailsLicense.Value = photo.DetailsLicense
			frm.DetailsLicense.Action = batch.ActionNone
		} else if photo.DetailsLicense != frm.DetailsLicense.Value {
			frm.DetailsLicense.Mixed = true
			frm.DetailsLicense.Value = ""
		}
	}

	// Return initialized PhotosForm.
	return frm
}
