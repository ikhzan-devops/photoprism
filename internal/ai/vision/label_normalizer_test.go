package vision

import "testing"

func TestCanonicalLabelFor(t *testing.T) {
	meta, ok := canonicalLabelFor("Sea Lion")
	if !ok {
		t.Fatalf("expected canonical entry for sea lion")
	}

	if meta.Name != "Sea Lion" {
		t.Fatalf("expected canonical name Sea Lion, got %q", meta.Name)
	}

	metaLower, ok := canonicalLabelFor("sea lion")
	if !ok || metaLower.Name != "Sea Lion" {
		t.Fatalf("expected lookup to be case-insensitive, got %v %q", ok, metaLower.Name)
	}
}

func TestCanonicalLabelForUnknown(t *testing.T) {
	if _, ok := canonicalLabelFor("unknown-label-xyz"); ok {
		t.Fatalf("expected no canonical entry")
	}
}

func TestPriorityFromTopicality(t *testing.T) {
	cases := []struct {
		top float32
		exp int
	}{
		{0.95, 5},
		{0.80, 3},
		{0.65, 2},
		{0.50, 1},
		{0.35, 0},
		{0.20, -1},
		{0.05, -2},
	}

	for _, tc := range cases {
		if got := priorityFromTopicality(tc.top); got != tc.exp {
			t.Fatalf("topicality %v expected priority %d, got %d", tc.top, tc.exp, got)
		}
	}
}

func TestNormalizeLabelResultCanonical(t *testing.T) {
	label := LabelResult{Name: "sea lion", Confidence: 0.8, Topicality: 0.7}
	normalizeLabelResult(&label)

	if label.Name != "Sea Lion" {
		t.Fatalf("expected canonical name, got %q", label.Name)
	}

	if label.Priority != priorityFromTopicality(0.7) {
		t.Fatalf("expected priority derived from topicality, got %d", label.Priority)
	}

	if len(label.Categories) == 0 {
		t.Fatalf("expected categories to be set")
	}
}

func TestNormalizeLabelResultFallback(t *testing.T) {
	label := LabelResult{Name: "kittens", Topicality: 0.25}
	normalizeLabelResult(&label)

	if label.Name == "" {
		t.Fatalf("expected non-empty name")
	}

	if label.Priority == 0 {
		t.Fatalf("expected priority to be derived from topicality")
	}
}
