package tensorflow

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
)

func loadLabelsFromPath(path string) (labels []string, err error) {
	log.Infof("tensorflow: loading model labels from %s", path)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Labels are separated by newlines
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}

	err = scanner.Err()

	return labels, err
}

// LoadLabels loads the labels of classification models from the specified path and returns them.
func LoadLabels(modelPath string, expectedLabels int) (labels []string, err error) {

	dir := os.DirFS(modelPath)
	matches, err := fs.Glob(dir, "labels*.txt")
	if err != nil {
		return nil, err
	}

	for i := range matches {
		labels, err := loadLabelsFromPath(filepath.Join(modelPath, matches[i]))
		if err != nil {
			return nil, err
		}

		switch expectedLabels - len(labels) {
		case 0:
			log.Infof("Found a valid labels file: %s", matches[i])
			return labels, nil
		case 1:
			log.Infof("Found a valid labels file %s but we have to add bias", matches[i])

			return append([]string{"background"}, labels...), nil
		default:
			log.Infof("File not valid. Expected %d labels and have %d",
				expectedLabels, len(labels))
		}
	}
	return nil, os.ErrNotExist
}
