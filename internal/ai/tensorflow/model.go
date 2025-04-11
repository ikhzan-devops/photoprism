package tensorflow

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	tf "github.com/wamuir/graft/tensorflow"

	"github.com/photoprism/photoprism/pkg/clean"
)

// SavedModel loads a saved TensorFlow model from the specified path.
func SavedModel(modelPath string, tags []string) (model *tf.SavedModel, err error) {
	log.Infof("tensorflow: loading %s", clean.Log(filepath.Base(modelPath)))

	if len(tags) == 0 {
		tags = []string{"serve"}
	}

	return tf.LoadSavedModel(modelPath, tags, nil)
}

// GuessInputAndOutput tries to inspect a loaded saved model to build the
// ModelInfo struct
func GuessInputAndOutput(model *tf.SavedModel) (input *PhotoInput, output *ModelOutput, err error) {
	modelOps := model.Graph.Operations()

	for i := range modelOps {
		if strings.HasPrefix(modelOps[i].Type(), "Placeholder") && modelOps[i].NumOutputs() == 1 && modelOps[i].Output(0).Shape().NumDimensions() == 4 {
			shape := modelOps[i].Output(0).Shape()
			input = &PhotoInput{
				Name:     modelOps[i].Name(),
				Height:   shape.Size(1),
				Width:    shape.Size(2),
				Channels: shape.Size(3),
			}
		} else if (modelOps[i].Type() == "Softmax" || strings.HasPrefix(modelOps[i].Type(), "StatefulPartitionedCall")) &&
			modelOps[i].NumOutputs() == 1 && modelOps[i].Output(0).Shape().NumDimensions() == 2 {
			output = &ModelOutput{
				Name:       modelOps[i].Name(),
				NumOutputs: modelOps[i].Output(0).Shape().Size(1),
			}
		}

		if input != nil && output != nil {
			return
		}
	}

	return nil, nil, fmt.Errorf("Could not guess the inputs and outputs")
}

func GetInputAndOutputFromSavedModel(model *tf.SavedModel) (*PhotoInput, *ModelOutput, error) {
	if model == nil {
		return nil, nil, fmt.Errorf("GetInputAndOutputFromSavedModel: nil input")
	}

	for _, v := range model.Signatures {
		inputs := v.Inputs
		outputs := v.Outputs

		if len(inputs) == 1 && len(outputs) == 1 {
			_, inputTensor := GetOne(inputs)
			outputVarName, outputTensor := GetOne(outputs)

			if inputTensor != nil && outputTensor != nil {
				if inputTensor.Shape.NumDimensions() == 4 &&
					outputTensor.Shape.NumDimensions() == 2 {
					var inputIdx, outputIdx = 0, 0
					var err error

					inputName, inputIndex, found := strings.Cut(inputTensor.Name, ":")
					if found {
						inputIdx, err = strconv.Atoi(inputIndex)
						if err != nil {
							return nil, nil, fmt.Errorf("Could not parse index %s: %w", inputIndex, err)
						}
					}

					outputName, outputIndex, found := strings.Cut(outputTensor.Name, ":")
					if found {
						outputIdx, err = strconv.Atoi(outputIndex)
						if err != nil {
							return nil, nil, fmt.Errorf("Could not parse index: %s: %w", outputIndex, err)
						}
					}

					return &PhotoInput{
							Name:        inputName,
							OutputIndex: inputIdx,
							Height:      inputTensor.Shape.Size(1),
							Width:       inputTensor.Shape.Size(2),
							Channels:    inputTensor.Shape.Size(3),
						}, &ModelOutput{
							Name:          outputName,
							OutputIndex:   outputIdx,
							NumOutputs:    outputTensor.Shape.Size(1),
							OutputsLogits: strings.Contains(Deref(outputVarName, ""), "logits"),
						}, nil

				}
			}

		}
	}

	return nil, nil, fmt.Errorf("GetInputAndOutputFromSignature: could not find valid signatures")
}
