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
	Name        string
	OutputIndex int
	Height      int64
	Width       int64
	Channels    int64
}

// When dimensions are not defined, it means the model accepts any size of
// photo
func (f PhotoInput) IsDynamic() bool {
	return f.Height == -1 && f.Width == -1
}

// Get the resolution
func (f PhotoInput) Resolution() int {
	return int(f.Height)
}

// Set the resolution: same height and width
func (f *PhotoInput) SetResolution(resolution int) {
	f.Height = int64(resolution)
	f.Width = int64(resolution)
}

// The output expected for a model
type ModelOutput struct {
	Name          string
	OutputIndex   int
	NumOutputs    int64
	OutputsLogits bool
}

// The meta information for the model
type ModelInfo struct {
	TFVersion string
	Tags      []string
	Input     *PhotoInput
	Output    *ModelOutput
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
			log.Printf("Could not get the inputs and outputs from signatures. (TF Version %s): %w", newModel.TFVersion, err)
		}

		models = append(models, newModel)
	}

	return models, nil
}
