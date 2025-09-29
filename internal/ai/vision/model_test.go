package vision

import (
	"testing"

	"github.com/photoprism/photoprism/internal/ai/tensorflow"
	"github.com/photoprism/photoprism/internal/ai/vision/ollama"
)

func TestModelGetOptionsDefaultsOllamaLabels(t *testing.T) {
	model := &Model{
		Type:     ModelTypeLabels,
		Provider: ollama.ProviderName,
	}

	model.ApplyProviderDefaults()

	opts := model.GetOptions()
	if opts == nil {
		t.Fatalf("expected options for labels model")
	}

	if opts.Temperature != DefaultTemperature {
		t.Errorf("unexpected temperature: got %v want %v", opts.Temperature, DefaultTemperature)
	}

	if opts.TopP != 0.9 {
		t.Errorf("unexpected top_p: got %v want 0.9", opts.TopP)
	}

	if len(opts.Stop) != 1 || opts.Stop[0] != "\n\n" {
		t.Fatalf("expected default stop sequence, got %#v", opts.Stop)
	}

	if opts != model.GetOptions() {
		t.Errorf("expected cached options pointer")
	}
}

func TestModelGetOptionsRespectsCustomValues(t *testing.T) {
	model := &Model{
		Type:     ModelTypeLabels,
		Provider: ollama.ProviderName,
		Options: &ApiRequestOptions{
			Temperature: 5,
			TopP:        0.95,
			Stop:        []string{"CUSTOM"},
		},
	}

	model.ApplyProviderDefaults()

	opts := model.GetOptions()
	if opts.Temperature != MaxTemperature {
		t.Errorf("temperature clamp failed: got %v want %v", opts.Temperature, MaxTemperature)
	}
	if opts.TopP != 0.95 {
		t.Errorf("top_p override lost: got %v", opts.TopP)
	}
	if len(opts.Stop) != 1 || opts.Stop[0] != "CUSTOM" {
		t.Errorf("stop override lost: %#v", opts.Stop)
	}
}

func TestModelGetOptionsFillsMissingFields(t *testing.T) {
	model := &Model{
		Type:     ModelTypeLabels,
		Provider: ollama.ProviderName,
		Options:  &ApiRequestOptions{},
	}

	model.ApplyProviderDefaults()

	opts := model.GetOptions()
	if opts.TopP != 0.9 {
		t.Errorf("expected default top_p, got %v", opts.TopP)
	}
	if len(opts.Stop) != 1 || opts.Stop[0] != "\n\n" {
		t.Errorf("expected default stop sequence, got %#v", opts.Stop)
	}
}

func TestModel_IsDefault(t *testing.T) {
	nasnetCopy := *NasnetModel
	nasnetCopy.Default = false

	cases := []struct {
		name  string
		model *Model
		want  bool
	}{
		{
			name:  "DefaultFlag",
			model: &Model{Default: true},
			want:  true,
		},
		{
			name:  "NasnetCopy",
			model: &nasnetCopy,
			want:  true,
		},
		{
			name: "CustomTensorFlow",
			model: &Model{
				Type:       ModelTypeLabels,
				Name:       "custom",
				TensorFlow: &tensorflow.ModelInfo{},
			},
			want: false,
		},
		{
			name: "RemoteService",
			model: &Model{
				Type:     ModelTypeCaption,
				Name:     "custom-caption",
				Provider: ollama.ProviderName,
			},
			want: false,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			if got := tc.model.IsDefault(); got != tc.want {
				t.Fatalf("IsDefault() = %v, want %v", got, tc.want)
			}
		})
	}
}
