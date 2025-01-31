package form

import "github.com/ulule/deepcopier"

// Label represents a label edit form.
type Label struct {
	LabelName        string `json:"Name"`
	LabelPriority    int    `json:"Priority"`
	LabelFavorite    bool   `json:"Favorite"`
	LabelDescription string `json:"Description"`
	LabelNotes       string `json:"Notes"`
	Thumb            string `json:"Thumb"`
	ThumbSrc         string `json:"ThumbSrc"`
	Uncertainty      int    `json:"Uncertainty"`
}

// NewLabel creates a new form struct based on the interface values.
func NewLabel(m interface{}) (*Label, error) {
	frm := &Label{}
	err := deepcopier.Copy(m).To(frm)
	return frm, err
}
