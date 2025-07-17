package tensorflow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	pb "github.com/wamuir/graft/tensorflow/core/protobuf/for_core_protos_go_proto"
	"google.golang.org/protobuf/proto"
)

// The number of channels expected. This is a fixed value because there
// seems to be an standard for input images defined as "What decodeImage
// returns"
const ExpectedChannels = 3

// Interval of allowed values
type Interval struct {
	Start  float32  `yaml:"Start,omitempty" json:"start,omitempty"`
	End    float32  `yaml:"End,omitempty" json:"end,omitempty"`
	Mean   *float32 `yaml:"Mean,omitempty" json:"mean,omitempty"`
	StdDev *float32 `yaml:"StdDev,omitempty" json:"stdDev,omitempty"`
}

// The size/mean of the interval
func (i Interval) Size() float32 {
	return i.End - i.Start
}

// The offset of the interval
func (i Interval) Offset() float32 {
	if i.StdDev == nil {
		return i.Start
	} else {
		return *i.StdDev
	}
}

// The standard interval returned by decodeImage is [0, 1]
func StandardInterval() *Interval {
	return &Interval{
		Start: 0.0,
		End:   1.0,
	}
}

// How should we resize the images
// JSON and YAML functions are given to make it
// user friendly from the configuration files
type ResizeOperation int

const (
	UndefinedResizeOperation ResizeOperation = iota
	ResizeBreakAspectRatio
	CenterCrop
	Padding
)

func (o ResizeOperation) String() string {
	switch o {
	case UndefinedResizeOperation:
		return "Undefined"
	case ResizeBreakAspectRatio:
		return "ResizeBreakAspectRatio"
	case CenterCrop:
		return "CenterCrop"
	case Padding:
		return "Padding"
	default:
		return "Unknown"
	}
}

func NewResizeOperation(s string) (ResizeOperation, error) {
	switch s {
	case "Undefined":
		return UndefinedResizeOperation, nil
	case "ResizeBreakAspectRatio":
		return ResizeBreakAspectRatio, nil
	case "CenterCrop":
		return CenterCrop, nil
	case "Padding":
		return Padding, nil
	default:
		return UndefinedResizeOperation, fmt.Errorf("Invalid operation %s", s)
	}
}

func (o ResizeOperation) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o *ResizeOperation) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, err := NewResizeOperation(s)
	if err != nil {
		return err
	}
	*o = val

	return nil
}

func (o ResizeOperation) MarshalYAML() (any, error) {
	return o.String(), nil
}

func (o *ResizeOperation) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	val, err := NewResizeOperation(s)
	if err != nil {
		return err
	}
	*o = val
	return nil
}

// How should we order the input vectors
// JSON and YAML functions are given to make it
// user friendly from the configuration files
type ColorChannelOrder int

const (
	UndefinedOrder ColorChannelOrder = 0
	RGB                              = 123
	RBG                              = 132
	GRB                              = 213
	GBR                              = 231
	BRG                              = 312
	BGR                              = 321
)

func (o ColorChannelOrder) Indices() (r, g, b int) {
	i := int(o)

	if i == 0 {
		i = 123
	}

	for idx := 0; i > 0 && idx < 3; idx += 1 {
		remainder := i % 10
		i /= 10

		switch remainder {
		case 1:
			r = 2 - idx
		case 2:
			g = 2 - idx
		case 3:
			b = 2 - idx
		}
	}

	return
}

func (o ColorChannelOrder) String() string {
	value := int(o)

	if value == 0 {
		value = 123
	}

	convert := func(remainder int) string {
		switch remainder {
		case 1:
			return "R"
		case 2:
			return "G"
		case 3:
			return "B"
		default:
			return "?"
		}
	}

	result := ""
	for value > 0 {
		remainder := value % 10
		value /= 10

		result = convert(remainder) + result
	}

	return result
}

func NewColorChannelOrder(val string) (ColorChannelOrder, error) {
	if len(val) != 3 {
		return UndefinedOrder, fmt.Errorf("Invalid length, expected 3")
	}

	convert := func(c rune) int {
		switch c {
		case 'R':
			return 1
		case 'G':
			return 2
		case 'B':
			return 3
		default:
			return 0
		}
	}

	result := 0
	for _, c := range val {
		index := convert(c)
		if index == 0 {
			return UndefinedOrder, fmt.Errorf("Invalid val %c", c)
		}
		result = result*10 + index
	}
	return ColorChannelOrder(result), nil
}

func (o ColorChannelOrder) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o *ColorChannelOrder) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, err := NewColorChannelOrder(s)
	if err != nil {
		return err
	}
	*o = val

	return nil
}

func (o ColorChannelOrder) MarshalYAML() (any, error) {
	return o.String(), nil
}

func (o *ColorChannelOrder) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	val, err := NewColorChannelOrder(s)
	if err != nil {
		return err
	}
	*o = val
	return nil
}

// Input description for a photo input for a model
type PhotoInput struct {
	Name              string            `yaml:"Name,omitempty" json:"name,omitempty"`
	Intervals         []Interval        `yaml:"Intervals,omitempty" json:"intervals,omitempty"`
	ResizeOperation   ResizeOperation   `yaml:"ResizeOperation,omitempty" json:"resizeOperation,omitemty"`
	ColorChannelOrder ColorChannelOrder `yaml:"ColorChannelOrder,omitempty" json:"inputOrder,omitempty"`
	OutputIndex       int               `yaml:"Index,omitempty" json:"index,omitempty"`
	Height            int64             `yaml:"Height,omitempty" json:"height,omitempty"`
	Width             int64             `yaml:"Width,omitempty" json:"width,omitempty"`
}

// When dimensions are not defined, it means the model accepts any size of
// photo
func (p PhotoInput) IsDynamic() bool {
	return p.Height == -1 && p.Width == -1
}

// Get the resolution
func (p PhotoInput) Resolution() int {
	return int(p.Height)
}

// Set the resolution: same height and width
func (p *PhotoInput) SetResolution(resolution int) {
	p.Height = int64(resolution)
	p.Width = int64(resolution)
}

// Get the interval or the default one.
// If just one interval has been fixed, then we assume
// it is the same for every channel. If no intervals
// have been defined, the default [0, 1] is returned
func (p PhotoInput) GetInterval(channel int) *Interval {
	if len(p.Intervals) <= channel {
		if len(p.Intervals) == 1 {
			return &p.Intervals[0]
		}
		return StandardInterval()
	} else {
		return &p.Intervals[channel]
	}
}

// Merge other input with this.
func (p *PhotoInput) Merge(other *PhotoInput) {
	if p.Name == "" {
		p.Name = other.Name
	}

	if p.Intervals == nil && other.Intervals != nil {
		p.Intervals = other.Intervals
	}

	if p.OutputIndex == 0 {
		p.OutputIndex = other.OutputIndex
	}

	if p.Height == 0 {
		p.Height = other.Height
	}

	if p.Width == 0 {
		p.Width = other.Width
	}

	if p.ResizeOperation == UndefinedResizeOperation {
		p.ResizeOperation = other.ResizeOperation
	}

	if p.ColorChannelOrder == UndefinedOrder {
		p.ColorChannelOrder = other.ColorChannelOrder
	}
}

// The output expected for a model
type ModelOutput struct {
	Name          string `yaml:"Name,omitempty" json:"name,omitempty"`
	OutputIndex   int    `yaml:"Index,omitempty" json:"index,omitempty"`
	NumOutputs    int64  `yaml:"Outputs,omitempty" json:"outputs,omitempty"`
	OutputsLogits bool   `yaml:"Logits,omitempty" json:"logits,omitempty"`
}

// Merge other output with this
func (m *ModelOutput) Merge(other *ModelOutput) {
	if m.Name == "" {
		m.Name = other.Name
	}

	if m.OutputIndex == 0 {
		m.OutputIndex = other.OutputIndex
	}

	if m.NumOutputs == 0 {
		m.NumOutputs = other.NumOutputs
	}

	if !m.OutputsLogits {
		m.OutputsLogits = other.OutputsLogits
	}
}

// The meta information for the model
type ModelInfo struct {
	TFVersion string       `yaml:"-" json:"-"`
	Tags      []string     `yaml:"Tags" json:"tags"`
	Input     *PhotoInput  `yaml:"Input" json:"input"`
	Output    *ModelOutput `yaml:"Output" json:"output"`
}

// Merge other model info. In case of having information
// for a field, the current model will keep its current value
func (m *ModelInfo) Merge(other *ModelInfo) {
	if m.TFVersion == "" {
		m.TFVersion = other.TFVersion
	}

	if len(m.Tags) == 0 {
		m.Tags = other.Tags
	}

	if m.Input == nil {
		m.Input = other.Input
	} else if other.Input != nil {
		m.Input.Merge(other.Input)
	}

	if m.Output == nil {
		m.Output = other.Output
	} else if other.Output != nil {
		m.Output.Merge(other.Output)
	}
}

// We consider a model complete if we know its inputs and outputs
func (m ModelInfo) IsComplete() bool {
	return m.Input != nil && m.Output != nil
}

// GetInputAndOutputFromSignature gets the signatures from a MetaGraphDef and
// uses them to build PhotoInput and ModelOutput structs, that will complete
// ModelInfo struct
func GetInputAndOutputFromMetaSignature(meta *pb.MetaGraphDef) (*PhotoInput, *ModelOutput, error) {
	if meta == nil {
		return nil, nil, fmt.Errorf("GetInputAndOutputFromSignature: nil input")
	}

	sig := meta.GetSignatureDef()
	for k, v := range sig {
		inputs := v.GetInputs()
		outputs := v.GetOutputs()

		if len(inputs) == 1 && len(outputs) == 1 {
			_, inputTensor := GetOne(inputs)
			outputVarName, outputTensor := GetOne(outputs)

			if inputTensor != nil && (*inputTensor).GetTensorShape() != nil &&
				outputTensor != nil && (*outputTensor).GetTensorShape() != nil {
				inputDims := (*inputTensor).GetTensorShape().Dim
				outputDims := (*outputTensor).GetTensorShape().Dim

				if inputDims[3].GetSize() != ExpectedChannels {
					log.Warnf("tensorflow: skipping signature %v because channels are expected to be %d, have %d",
						k, ExpectedChannels, inputDims[3].GetSize())
				}

				if len(inputDims) == 4 &&
					inputDims[3].GetSize() == ExpectedChannels &&
					len(outputDims) == 2 {
					var err error
					var inputIdx, outputIdx = 0, 0

					inputName, inputIndex, found := strings.Cut((*inputTensor).GetName(), ":")
					if found {

						inputIdx, err = strconv.Atoi(inputIndex)
						if err != nil {
							return nil, nil, fmt.Errorf("Could not parse index %s: %w", inputIndex, err)
						}
					}

					outputName, outputIndex, found := strings.Cut((*outputTensor).GetName(), ":")
					if found {

						outputIdx, err = strconv.Atoi(outputIndex)
						if err != nil {
							return nil, nil, fmt.Errorf("Could not parse index: %s: %w", outputIndex, err)
						}
					}

					return &PhotoInput{
							Name:        inputName,
							OutputIndex: inputIdx,
							Height:      inputDims[1].GetSize(),
							Width:       inputDims[2].GetSize(),
						}, &ModelOutput{
							Name:          outputName,
							OutputIndex:   outputIdx,
							NumOutputs:    outputDims[1].GetSize(),
							OutputsLogits: strings.Contains(Deref(outputVarName, ""), "logits"),
						}, nil

				}
			}

		}

	}

	return nil, nil, fmt.Errorf("GetInputAndOutputFromMetaSignature: Could not find a valid signature")
}

func GetModelInfo(path string) ([]ModelInfo, error) {

	path = filepath.Join(path, "saved_model.pb")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Could not read the file %s: %w", path, err)
	}

	model := new(pb.SavedModel)

	err = proto.Unmarshal(data, model)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal the file %s: %w", path, err)
	}

	models := make([]ModelInfo, 0)
	metas := model.GetMetaGraphs()
	for i := range metas {
		def := metas[i].GetMetaInfoDef()
		input, output, err := GetInputAndOutputFromMetaSignature(metas[i])
		newModel := ModelInfo{
			TFVersion: def.GetTensorflowVersion(),
			Tags:      def.GetTags(),
			Input:     input,
			Output:    output,
		}

		if err != nil {
			log.Errorf("Could not get the inputs and outputs from signatures. (TF Version %s): %w", newModel.TFVersion, err)
		}

		models = append(models, newModel)
	}

	return models, nil
}
