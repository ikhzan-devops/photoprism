package batch

import (
	"github.com/photoprism/photoprism/internal/entity"
)

// PhotoForm represents photo batch edit form values.
type PhotoForm struct {
	PhotoType        String  `json:"Type"`
	PhotoTitle       String  `json:"Title"`
	PhotoCaption     String  `json:"Caption"`
	TakenAt          Time    `json:"TakenAt"`
	TakenAtLocal     Time    `json:"TakenAtLocal"`
	PhotoDay         Int     `json:"Day"`
	PhotoMonth       Int     `json:"Month"`
	PhotoYear        Int     `json:"Year"`
	TimeZone         String  `json:"TimeZone"`
	PhotoCountry     String  `json:"Country"`
	PhotoAltitude    Int     `json:"Altitude"`
	PhotoLat         Float64 `json:"Lat"`
	PhotoLng         Float64 `json:"Lng"`
	PhotoIso         Int     `json:"Iso"`
	PhotoFocalLength Int     `json:"FocalLength"`
	PhotoFNumber     Float32 `json:"FNumber"`
	PhotoExposure    String  `json:"Exposure"`
	PhotoFavorite    Bool    `json:"Favorite"`
	PhotoPrivate     Bool    `json:"Private"`
	PhotoScan        Bool    `json:"Scan"`
	PhotoPanorama    Bool    `json:"Panorama"`
	CameraID         UInt    `json:"CameraID"`
	LensID           UInt    `json:"LensID"`

	DetailsSubject   String `json:"Subject"`
	DetailsArtist    String `json:"Artist"`
	DetailsCopyright String `json:"Copyright"`
	DetailsLicense   String `json:"License"`
}

func NewPhotoForm(photos entity.Photos) *PhotoForm {
	frm := &PhotoForm{}

	for _, photo := range photos {
		if photo.PhotoType != "" && frm.PhotoType.Value == "" {
			frm.PhotoType.Value = photo.PhotoType
			frm.PhotoType.Matches = true
		} else if photo.PhotoType != frm.PhotoType.Value {
			frm.PhotoType.Matches = false
		}

		if photo.PhotoTitle != "" && frm.PhotoTitle.Value == "" {
			frm.PhotoTitle.Value = photo.PhotoTitle
			frm.PhotoTitle.Matches = true
		} else if photo.PhotoTitle != frm.PhotoTitle.Value {
			frm.PhotoTitle.Matches = false
		}

		if photo.PhotoCaption != "" && frm.PhotoCaption.Value == "" {
			frm.PhotoCaption.Value = photo.PhotoCaption
			frm.PhotoTitle.Matches = true
		} else if photo.PhotoCaption != frm.PhotoCaption.Value {
			frm.PhotoCaption.Matches = false
		}

		if !photo.TakenAt.IsZero() && frm.TakenAt.Value.IsZero() {
			frm.TakenAt.Value = photo.TakenAt
			frm.TakenAt.Matches = true
		} else if photo.TakenAt != frm.TakenAt.Value {
			frm.TakenAt.Matches = false
		}

		if !photo.TakenAtLocal.IsZero() && frm.TakenAtLocal.Value.IsZero() {
			frm.TakenAtLocal.Value = photo.TakenAtLocal
			frm.TakenAtLocal.Matches = true
		} else if photo.TakenAtLocal != frm.TakenAtLocal.Value {
			frm.TakenAtLocal.Matches = false
		}

		if photo.PhotoDay > 0 && frm.PhotoDay.Value == 0 {
			frm.PhotoDay.Value = photo.PhotoDay
			frm.PhotoDay.Matches = true
		} else if photo.PhotoDay != frm.PhotoDay.Value {
			frm.PhotoDay.Matches = false
		}

		if photo.PhotoMonth > 0 && frm.PhotoMonth.Value == 0 {
			frm.PhotoMonth.Value = photo.PhotoMonth
			frm.PhotoMonth.Matches = true
		} else if photo.PhotoMonth != frm.PhotoMonth.Value {
			frm.PhotoMonth.Matches = false
		}

		if photo.PhotoYear > 0 && frm.PhotoYear.Value == 0 {
			frm.PhotoYear.Value = photo.PhotoYear
			frm.PhotoYear.Matches = true
		} else if photo.PhotoYear != frm.PhotoYear.Value {
			frm.PhotoYear.Matches = false
		}

		if photo.TimeZone != "" && frm.TimeZone.Value == "" {
			frm.TimeZone.Value = photo.TimeZone
			frm.TimeZone.Matches = true
		} else if photo.TimeZone != frm.TimeZone.Value {
			frm.TimeZone.Matches = false
		}

		if photo.PhotoCountry != "" && frm.PhotoCountry.Value == "" {
			frm.PhotoCountry.Value = photo.PhotoCountry
			frm.PhotoCountry.Matches = true
		} else if photo.PhotoCountry != frm.PhotoCountry.Value {
			frm.PhotoCountry.Matches = false
		}

		if photo.PhotoAltitude != 0 && frm.PhotoAltitude.Value == 0 {
			frm.PhotoAltitude.Value = photo.PhotoAltitude
			frm.PhotoAltitude.Matches = true
		} else if photo.PhotoAltitude != frm.PhotoAltitude.Value {
			frm.PhotoAltitude.Matches = false
		}

		if photo.PhotoLat != 0.0 && frm.PhotoLat.Value == 0.0 {
			frm.PhotoLat.Value = photo.PhotoLat
			frm.PhotoLat.Matches = true
		} else if photo.PhotoLat != frm.PhotoLat.Value {
			frm.PhotoLat.Matches = false
		}

		if photo.PhotoLng != 0.0 && frm.PhotoLng.Value == 0.0 {
			frm.PhotoLng.Value = photo.PhotoLng
			frm.PhotoLng.Matches = true
		} else if photo.PhotoLng != frm.PhotoLng.Value {
			frm.PhotoLng.Matches = false
		}

		if photo.PhotoIso != 0 && frm.PhotoIso.Value == 0 {
			frm.PhotoIso.Value = photo.PhotoIso
			frm.PhotoIso.Matches = true
		} else if photo.PhotoIso != frm.PhotoIso.Value {
			frm.PhotoIso.Matches = false
		}

		if photo.PhotoFocalLength != 0 && frm.PhotoFocalLength.Value == 0 {
			frm.PhotoFocalLength.Value = photo.PhotoFocalLength
			frm.PhotoFocalLength.Matches = true
		} else if photo.PhotoFocalLength != frm.PhotoFocalLength.Value {
			frm.PhotoFocalLength.Matches = false
		}

		if photo.PhotoFNumber != 0.0 && frm.PhotoFNumber.Value == 0.0 {
			frm.PhotoFNumber.Value = photo.PhotoFNumber
			frm.PhotoFNumber.Matches = true
		} else if photo.PhotoFNumber != frm.PhotoFNumber.Value {
			frm.PhotoFNumber.Matches = false
		}

		if photo.PhotoExposure != "" && frm.PhotoExposure.Value == "" {
			frm.PhotoExposure.Value = photo.PhotoExposure
			frm.PhotoExposure.Matches = true
		} else if photo.PhotoExposure != frm.PhotoExposure.Value {
			frm.PhotoExposure.Matches = false
		}

		if photo.PhotoFavorite && !frm.PhotoFavorite.Value {
			frm.PhotoFavorite.Value = photo.PhotoFavorite
			frm.PhotoFavorite.Matches = true
		} else if photo.PhotoFavorite != frm.PhotoFavorite.Value {
			frm.PhotoFavorite.Matches = false
		}

		if photo.PhotoPrivate && !frm.PhotoPrivate.Value {
			frm.PhotoPrivate.Value = photo.PhotoPrivate
			frm.PhotoPrivate.Matches = true
		} else if photo.PhotoPrivate != frm.PhotoPrivate.Value {
			frm.PhotoPrivate.Matches = false
		}

		if photo.PhotoScan && !frm.PhotoScan.Value {
			frm.PhotoScan.Value = photo.PhotoScan
			frm.PhotoScan.Matches = true
		} else if photo.PhotoScan != frm.PhotoScan.Value {
			frm.PhotoScan.Matches = false
		}

		if photo.PhotoPanorama && !frm.PhotoPanorama.Value {
			frm.PhotoPanorama.Value = photo.PhotoPanorama
			frm.PhotoPanorama.Matches = true
		} else if photo.PhotoPanorama != frm.PhotoPanorama.Value {
			frm.PhotoPanorama.Matches = false
		}

		if photo.CameraID != 0 && frm.CameraID.Value == 0 {
			frm.CameraID.Value = photo.CameraID
			frm.CameraID.Matches = true
		} else if photo.CameraID != frm.CameraID.Value {
			frm.CameraID.Matches = false
		}

		if photo.LensID != 0 && frm.LensID.Value == 0 {
			frm.LensID.Value = photo.LensID
			frm.LensID.Matches = true
		} else if photo.LensID != frm.LensID.Value {
			frm.LensID.Matches = false
		}

		if photo.Details != nil {
			if photo.Details.Subject != "" && frm.DetailsSubject.Value == "" {
				frm.DetailsSubject.Value = photo.Details.Subject
				frm.DetailsSubject.Matches = true
			} else if photo.Details.Subject != frm.DetailsSubject.Value {
				frm.DetailsSubject.Matches = false
			}

			if photo.Details.Artist != "" && frm.DetailsArtist.Value == "" {
				frm.DetailsArtist.Value = photo.Details.Artist
				frm.DetailsArtist.Matches = true
			} else if photo.Details.Artist != frm.DetailsArtist.Value {
				frm.DetailsArtist.Matches = false
			}

			if photo.Details.Copyright != "" && frm.DetailsCopyright.Value == "" {
				frm.DetailsCopyright.Value = photo.Details.Copyright
				frm.DetailsCopyright.Matches = true
			} else if photo.Details.Copyright != frm.DetailsCopyright.Value {
				frm.DetailsCopyright.Matches = false
			}

			if photo.Details.License != "" && frm.DetailsLicense.Value == "" {
				frm.DetailsLicense.Value = photo.Details.License
				frm.DetailsLicense.Matches = true
			} else if photo.Details.License != frm.DetailsLicense.Value {
				frm.DetailsLicense.Matches = false
			}
		}

	}

	return frm
}
