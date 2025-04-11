package tensorflow

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	pb "github.com/wamuir/graft/tensorflow/core/protobuf/for_core_protos_go_proto"
	"google.golang.org/protobuf/proto"
)

// Input description for a photo input for a model
type PhotoInput struct {
	Name        string `yaml:"Name,omitempty" json:"name,omitempty"`
	OutputIndex int    `yaml:"Index,omitempty" json:"index,omitempty"`
	Height      int64  `yaml:"Height,omitempty" json:"height,omitempty"`
	Width       int64  `yaml:"Width,omitempty" json:"width,omitempty"`
	Channels    int64  `yaml:"Channels,omitempty" json:"channels,omitempty"`
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

// Merge other input with this.
func (p *PhotoInput) Merge(other *PhotoInput) {
	if p.Name == "" {
		p.Name = other.Name
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

	if p.Channels == 0 {
		p.Channels = other.Channels
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
	for _, v := range sig {
		inputs := v.GetInputs()
		outputs := v.GetOutputs()

		if len(inputs) == 1 && len(outputs) == 1 {
			_, inputTensor := GetOne(inputs)
			outputVarName, outputTensor := GetOne(outputs)

			if inputTensor != nil && (*inputTensor).GetTensorShape() != nil &&
				outputTensor != nil && (*outputTensor).GetTensorShape() != nil {
				inputDims := (*inputTensor).GetTensorShape().Dim
				outputDims := (*outputTensor).GetTensorShape().Dim

				if len(inputDims) == 4 && len(outputDims) == 2 {
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
							Channels:    inputDims[3].GetSize(),
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
