package batch

// PhotosForm represents photo batch edit form values.
type PhotosForm struct {
	PhotoType        String  `json:"Type,omitempty"`
	PhotoTitle       String  `json:"Title,omitempty"`
	PhotoCaption     String  `json:"Caption,omitempty"`
	TakenAt          Time    `json:"TakenAt,omitempty"`
	TakenAtLocal     Time    `json:"TakenAtLocal,omitempty"`
	PhotoDay         Int     `json:"Day,omitempty"`
	PhotoMonth       Int     `json:"Month,omitempty"`
	PhotoYear        Int     `json:"Year,omitempty"`
	TimeZone         String  `json:"TimeZone,omitempty"`
	PhotoCountry     String  `json:"Country,omitempty"`
	PhotoAltitude    Int     `json:"Altitude,omitempty"`
	PhotoLat         Float64 `json:"Lat,omitempty"`
	PhotoLng         Float64 `json:"Lng,omitempty"`
	PhotoIso         Int     `json:"Iso,omitempty"`
	PhotoFocalLength Int     `json:"FocalLength,omitempty"`
	PhotoFNumber     Float32 `json:"FNumber,omitempty"`
	PhotoExposure    String  `json:"Exposure,omitempty"`
	PhotoFavorite    Bool    `json:"Favorite,omitempty"`
	PhotoPrivate     Bool    `json:"Private,omitempty"`
	PhotoScan        Bool    `json:"Scan,omitempty"`
	PhotoPanorama    Bool    `json:"Panorama,omitempty"`
	CameraID         Int     `json:"CameraID,omitempty"`
	LensID           Int     `json:"LensID,omitempty"`
	Albums           Items   `json:"Albums,omitempty"`
	Labels           Items   `json:"Labels,omitempty"`

	DetailsKeywords  String `json:"DetailsKeywords,omitempty"`
	DetailsSubject   String `json:"DetailsSubject,omitempty"`
	DetailsArtist    String `json:"DetailsArtist,omitempty"`
	DetailsCopyright String `json:"DetailsCopyright,omitempty"`
	DetailsLicense   String `json:"DetailsLicense,omitempty"`
}
