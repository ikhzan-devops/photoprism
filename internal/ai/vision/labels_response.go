package vision

type LabelsResponse struct {
	Id     string        `yaml:"Id,omitempty" json:"id,omitempty"`
	Model  *Model        `yaml:"Model,omitempty" json:"model"`
	Result []LabelResult `yaml:"Result,omitempty" json:"result"`
}

func NewLabelsResponse(id string, model *Model, results []LabelResult) LabelsResponse {
	if model == nil {
		model = NasnetModel
	}

	if results == nil {
		results = []LabelResult{}
	}

	return LabelsResponse{
		Id:     id,
		Model:  model,
		Result: results,
	}
}

type LabelResult struct {
	Id         string  `yaml:"Id,omitempty" json:"id,omitempty"`
	Name       string  `yaml:"Name,omitempty" json:"name"`
	Category   string  `yaml:"Category,omitempty" json:"category"`
	Confidence float64 `yaml:"Confidence,omitempty" json:"confidence,omitempty"`
	Topicality float64 `yaml:"Topicality,omitempty" json:"topicality,omitempty"`
}
