package tensorflow

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

// Resize Operation Tests
var allOperations = []ResizeOperation{
	UndefinedResizeOperation,
	ResizeBreakAspectRatio,
	CenterCrop,
	Padding,
}

func TestResizeOperations(t *testing.T) {
	for i := range allOperations {
		text := allOperations[i].String()

		op, err := NewResizeOperation(text)
		if err != nil {
			t.Fatalf("Invalid operation %s: %v", text, err)
		}

		assert.Equal(t, op, allOperations[i])
	}
}

const exampleOperationJSON = `"CenterCrop"`

func TestResizeOperationJSON(t *testing.T) {
	var op ResizeOperation

	err := json.Unmarshal(
		[]byte(exampleOperationJSON), &op)

	if err != nil {
		t.Fatal("Could not unmarshal the example operation")
	}

	for i := range allOperations {
		serialized, err := json.Marshal(allOperations[i])
		if err != nil {
			t.Fatalf("Could not marshal %v: %v",
				allOperations[i], err)
		}

		err = json.Unmarshal(serialized, &op)
		if err != nil {
			t.Fatalf("Could not unmarshal %s: %v",
				string(serialized), err)
		}

		assert.Equal(t, op, allOperations[i])
	}
}

const exampleOperationYAML = "CenterCrop"

func TestResizeOperationYAML(t *testing.T) {
	var op ResizeOperation

	err := yaml.Unmarshal(
		[]byte(exampleOperationYAML), &op)

	if err != nil {
		t.Fatal("Could not unmarshal the example operation")
	}

	for i := range allOperations {
		serialized, err := yaml.Marshal(allOperations[i])
		if err != nil {
			t.Fatalf("Could not marshal %v: %v",
				allOperations[i], err)
		}

		err = yaml.Unmarshal(serialized, &op)
		if err != nil {
			t.Fatalf("Could not unmarshal %s: %v",
				string(serialized), err)
		}

		assert.Equal(t, op, allOperations[i])
	}
}

// Resize Operation Tests
var allColorChannelOrders = []ColorChannelOrder{
	RGB,
	RBG,
	GRB,
	GBR,
	BRG,
	BGR,
}

func TestColorChannelOrders(t *testing.T) {
	for i := range allColorChannelOrders {
		text := allColorChannelOrders[i].String()

		order, err := NewColorChannelOrder(text)
		if err != nil {
			t.Fatalf("Invalid order %s: %v", text, err)
		}

		assert.Equal(t, order, allColorChannelOrders[i])
	}
}

const exampleOrderJSON = `"RGB"`

func TestColorChannelOrderJSON(t *testing.T) {
	var order ColorChannelOrder

	err := json.Unmarshal(
		[]byte(exampleOrderJSON), &order)

	if err != nil {
		t.Fatal("Could not unmarshal the example operation")
	}

	for i := range allColorChannelOrders {
		serialized, err := json.Marshal(allColorChannelOrders[i])
		if err != nil {
			t.Fatalf("Could not marshal %v: %v",
				allColorChannelOrders[i], err)
		}

		err = json.Unmarshal(serialized, &order)
		if err != nil {
			t.Fatalf("Could not unmarshal %s: %v",
				string(serialized), err)
		}

		assert.Equal(t, order, allColorChannelOrders[i])
	}
}

const exampleOrderYAML = "RGB"

func TestColorChannelOrderYAML(t *testing.T) {
	var order ColorChannelOrder

	err := yaml.Unmarshal(
		[]byte(exampleOrderYAML), &order)

	if err != nil {
		t.Fatal("Could not unmarshal the example operation")
	}

	for i := range allColorChannelOrders {
		serialized, err := yaml.Marshal(allColorChannelOrders[i])
		if err != nil {
			t.Fatalf("Could not marshal %v: %v",
				allColorChannelOrders[i], err)
		}

		err = yaml.Unmarshal(serialized, &order)
		if err != nil {
			t.Fatalf("Could not unmarshal %s: %v",
				string(serialized), err)
		}

		assert.Equal(t, order, allColorChannelOrders[i])
	}
}

func TestOrderIndices(t *testing.T) {
	r, g, b := UndefinedOrder.Indices()

	assert.Equal(t, r, 0)
	assert.Equal(t, g, 1)
	assert.Equal(t, b, 2)

	powerFx := func(i int) int {
		switch i {
		case 0:
			return 100
		case 1:
			return 10
		case 2:
			return 1
		default:
			return -1
		}
	}

	for i := range allColorChannelOrders {
		r, g, b = allColorChannelOrders[i].Indices()
		assert.Equal(t, powerFx(r)+2*powerFx(g)+3*powerFx(b), int(allColorChannelOrders[i]))
	}
}
