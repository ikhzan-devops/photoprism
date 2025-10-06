package face

import "testing"

func TestParseEngine(t *testing.T) {
	cases := map[string]EngineName{
		"":         EngineAuto,
		"AUTO":     EngineAuto,
		"pigo":     EnginePigo,
		"  PIGO  ": EnginePigo,
		"onnx":     EngineONNX,
		"OnNx":     EngineONNX,
		"unknown":  EngineAuto,
	}

	for input, expected := range cases {
		if got := ParseEngine(input); got != expected {
			t.Fatalf("ParseEngine(%q) = %q, expected %q", input, got, expected)
		}
	}
}
