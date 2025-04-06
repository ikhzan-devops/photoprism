package classify

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"math"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
	tf "github.com/wamuir/graft/tensorflow"

	"github.com/photoprism/photoprism/pkg/clean"
)

// Model represents a TensorFlow classification model.
type Model struct {
	model      *tf.SavedModel
	modelPath  string
	assetsPath string
	resolution int
	modelTags  []string
	labels     []string
	disabled   bool
	mutex      sync.Mutex
}

// NewModel returns new TensorFlow classification model instance.
func NewModel(assetsPath, modelPath string, resolution int, modelTags []string, disabled bool) *Model {
	return &Model{assetsPath: assetsPath, modelPath: modelPath, resolution: resolution, modelTags: modelTags, disabled: disabled}
}

// NewNasnet returns new Nasnet TensorFlow classification model instance.
func NewNasnet(assetsPath string, disabled bool) *Model {
	return NewModel(assetsPath, "nasnet", 224, []string{"photoprism"}, disabled)
}

// Init initialises tensorflow models if not disabled
func (m *Model) Init() (err error) {
	if m.disabled {
		return nil
	}

	return m.loadModel()
}

// File returns matching labels for a jpeg media file.
func (m *Model) File(filename string, confidenceThreshold int) (result Labels, err error) {
	if m.disabled {
		return result, nil
	}

	imageBuffer, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return m.Labels(imageBuffer, confidenceThreshold)
}

// Labels returns matching labels for a jpeg media string.
func (m *Model) Labels(img []byte, confidenceThreshold int) (result Labels, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("classify: %s (inference panic)\nstack: %s", r, debug.Stack())
		}
	}()

	if m.disabled {
		return result, nil
	}

	if loadErr := m.loadModel(); loadErr != nil {
		return nil, loadErr
	}

	// Create tensor from image.
	tensor, err := m.createTensor(img)

	if err != nil {
		return nil, err
	}

	// Run inference.
	output, err := m.model.Session.Run(
		map[tf.Output]*tf.Tensor{
			m.model.Graph.Operation("input_1").Output(0): tensor,
		},
		[]tf.Output{
			m.model.Graph.Operation("predictions/Softmax").Output(0),
		},
		nil)

	if err != nil {
		return result, fmt.Errorf("classify: %s (run inference)", err.Error())
	}

	if len(output) < 1 {
		return result, fmt.Errorf("classify: inference failed, no output")
	}

	// Return best labels
	result = m.bestLabels(output[0].Value().([][]float32)[0], confidenceThreshold)

	if len(result) > 0 {
		log.Tracef("classify: image classified as %+v", result)
	}

	return result, nil
}

func (m *Model) loadLabels(path string) error {
	modelLabels := path + "/labels.txt"

	log.Infof("classify: loading labels from labels.txt")

	// Load labels
	f, err := os.Open(modelLabels)

	if err != nil {
		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Labels are separated by newlines
	for scanner.Scan() {
		m.labels = append(m.labels, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// ModelLoaded tests if the TensorFlow model is loaded.
func (m *Model) ModelLoaded() bool {
	return m.model != nil
}

func (m *Model) loadModel() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.ModelLoaded() {
		return nil
	}

	modelPath := path.Join(m.assetsPath, m.modelPath)

	log.Infof("classify: loading %s", clean.Log(filepath.Base(modelPath)))

	// Load model
	model, err := tf.LoadSavedModel(modelPath, m.modelTags, nil)

	if err != nil {
		return err
	}

	m.model = model

	return m.loadLabels(modelPath)
}

// bestLabels returns the best 5 labels (if enough high probability labels) from the prediction of the model
func (m *Model) bestLabels(probabilities []float32, confidenceThreshold int) Labels {
	var result Labels

	for i, p := range probabilities {
		if i >= len(m.labels) {
			// break if probabilities and labels does not match
			break
		}

		// discard labels with low probabilities
		if p < 0.1 {
			continue
		}

		labelText := strings.ToLower(m.labels[i])

		rule, _ := Rules.Find(labelText)

		// discard labels that don't met the threshold
		if p < rule.Threshold {
			continue
		}

		// Get rule label name instead of t.labels name if it exists
		if rule.Label != "" {
			labelText = rule.Label
		}

		labelText = strings.TrimSpace(labelText)

		confidence := int(math.Round(float64(p * 100)))

		if confidence >= confidenceThreshold {
			result = append(result, Label{Name: labelText, Source: SrcImage, Uncertainty: 100 - confidence, Priority: rule.Priority, Categories: rule.Categories})
		}
	}

	// Sort by probability
	sort.Sort(result)

	// Return the best labels only.
	if l := len(result); l < 5 {
		return result[:l]
	} else {
		return result[:5]
	}
}

// createTensor converts bytes jpeg image in a tensor object required as tensorflow model input
func (m *Model) createTensor(image []byte) (*tf.Tensor, error) {
	img, err := imaging.Decode(bytes.NewReader(image), imaging.AutoOrientation(true))

	if err != nil {
		return nil, err
	}

	width, height := m.resolution, m.resolution

	img = imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)

	return imageToTensor(img, width, height)
}

func imageToTensor(img image.Image, imageHeight, imageWidth int) (tfTensor *tf.Tensor, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("classify: %s (panic)\nstack: %s", r, debug.Stack())
		}
	}()

	if imageHeight <= 0 || imageWidth <= 0 {
		return tfTensor, fmt.Errorf("classify: image width and height must be > 0")
	}

	var tfImage [1][][][3]float32

	for j := 0; j < imageHeight; j++ {
		tfImage[0] = append(tfImage[0], make([][3]float32, imageWidth))
	}

	for i := 0; i < imageWidth; i++ {
		for j := 0; j < imageHeight; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			tfImage[0][j][i][0] = convertValue(r)
			tfImage[0][j][i][1] = convertValue(g)
			tfImage[0][j][i][2] = convertValue(b)
		}
	}

	return tf.NewTensor(tfImage)
}

func convertValue(value uint32) float32 {
	return (float32(value>>8) - float32(127.5)) / float32(127.5)
}
