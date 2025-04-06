package vision

import (
	"errors"
	"sort"

	"github.com/photoprism/photoprism/internal/ai/classify"
)

// Labels returns suitable labels for the specified image thumbnail.
func Labels(thumbnails []string) (result classify.Labels, err error) {
	if len(thumbnails) == 0 {
		return result, errors.New("missing thumbnail filenames")
	}

	if Config == nil {
		return result, errors.New("missing configuration")
	} else if len(Config.Labels) == 0 {
		return result, errors.New("missing labels model configuration")
	}

	config := Config.Labels[0]
	model := config.ClassifyModel()

	if model == nil {
		return result, errors.New("missing labels model")
	}

	for i := range thumbnails {
		labels, modelErr := model.File(thumbnails[i], Config.Thresholds.Confidence)

		if modelErr != nil {
			return result, modelErr
		}

		for j := range labels {
			found := false

			for k := range result {
				if labels[j].Name == result[k].Name {
					found = true

					if labels[j].Uncertainty < result[k].Uncertainty {
						result[k].Uncertainty = labels[j].Uncertainty
					}

					if labels[j].Priority > result[k].Priority {
						result[k].Priority = labels[j].Priority
					}
				}
			}

			if !found {
				result = append(result, labels...)
			}
		}
	}

	sort.Sort(result)

	return result, nil
}
