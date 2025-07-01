package places

type SearchResult struct {
	Id          string    `json:"id"`
	Name        string    `json:"name,omitempty"`
	City        string    `json:"city,omitempty"`
	Country     string    `json:"country"`
	Lat         float64   `json:"lat"`
	Lng         float64   `json:"lng"`
	Boundingbox []float64 `json:"boundingbox,omitempty"`
	Importance  float64   `json:"importance,omitempty"`
	Licence     string    `json:"licence,omitempty"`
}

type SearchResults = []SearchResult
