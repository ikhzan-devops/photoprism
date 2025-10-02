package search

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/entity/sortby"
	"github.com/photoprism/photoprism/internal/form"
)

func TestLabels(t *testing.T) {
	t.Run("search with query", func(t *testing.T) {
		query := form.NewLabelSearch("q:C")
		query.Count = 1005
		query.Order = "slug"
		result, err := Labels(query)

		if err != nil {
			t.Fatal(err)
		}

		t.Logf("results: %+v", result)

		assert.LessOrEqual(t, 2, len(result))

		for _, r := range result {
			assert.IsType(t, Label{}, r)
			assert.NotEmpty(t, r.ID)
			assert.NotEmpty(t, r.LabelName)
			assert.NotEmpty(t, r.LabelSlug)
			assert.NotEmpty(t, r.CustomSlug)

			if fix, ok := entity.LabelFixtures[r.LabelSlug]; ok {
				assert.Equal(t, fix.LabelName, r.LabelName)
				assert.Equal(t, fix.LabelSlug, r.LabelSlug)
				assert.Equal(t, fix.CustomSlug, r.CustomSlug)
			}
		}
	})
	t.Run("search for cow", func(t *testing.T) {
		query := form.NewLabelSearch("Q:cow")
		query.Count = 1005
		query.Order = "slug"
		result, err := Labels(query)

		if err != nil {
			t.Fatal(err)
		}

		t.Logf("results: %+v", result)

		assert.LessOrEqual(t, 1, len(result))

		for _, r := range result {
			assert.IsType(t, Label{}, r)
			assert.NotEmpty(t, r.ID)
			assert.NotEmpty(t, r.LabelName)
			assert.NotEmpty(t, r.LabelSlug)
			assert.NotEmpty(t, r.CustomSlug)

			if fix, ok := entity.LabelFixtures[r.LabelSlug]; ok {
				assert.Equal(t, fix.LabelName, r.LabelName)
				assert.Equal(t, fix.LabelSlug, r.LabelSlug)
				assert.Equal(t, fix.CustomSlug, r.CustomSlug)
			}
		}
	})
	t.Run("search for favorites", func(t *testing.T) {
		query := form.NewLabelSearch("Favorite:true")
		query.Count = 15
		result, err := Labels(query)

		if err != nil {
			t.Fatal(err)
		}

		assert.LessOrEqual(t, 2, len(result))

		for _, r := range result {
			assert.True(t, r.LabelFavorite)
			assert.IsType(t, Label{}, r)
			assert.NotEmpty(t, r.ID)
			assert.NotEmpty(t, r.LabelName)
			assert.NotEmpty(t, r.LabelSlug)
			assert.NotEmpty(t, r.CustomSlug)

			if fix, ok := entity.LabelFixtures[r.LabelSlug]; ok {
				assert.Equal(t, fix.LabelName, r.LabelName)
				assert.Equal(t, fix.LabelSlug, r.LabelSlug)
				assert.Equal(t, fix.CustomSlug, r.CustomSlug)
			}
		}
	})
	t.Run("order count", func(t *testing.T) {
		query := form.NewLabelSearch("")
		query.All = true
		query.Order = sortby.Count
		result, err := Labels(query)

		if err != nil {
			t.Fatal(err)
		}

		if len(result) < 2 {
			t.Fatalf("expected multiple labels")
		}

		if result[0].PhotoCount < result[1].PhotoCount {
			t.Fatalf("expected descending photo count")
		}
	})
	t.Run("order slug", func(t *testing.T) {
		query := form.NewLabelSearch("")
		query.All = true
		query.Order = sortby.Slug
		result, err := Labels(query)

		if err != nil {
			t.Fatal(err)
		}

		if len(result) < 2 {
			t.Fatalf("expected multiple labels")
		}

		if result[0].CustomSlug > result[1].CustomSlug {
			t.Fatalf("expected slug ascending")
		}
	})
	t.Run("default filter excludes low priority", func(t *testing.T) {
		query := form.NewLabelSearch("")
		result, err := Labels(query)

		if err != nil {
			t.Fatal(err)
		}

		for _, label := range result {
			if !label.LabelFavorite {
				if label.LabelPriority < 0 || label.PhotoCount <= 1 {
					t.Fatalf("label %s should have been filtered", label.LabelSlug)
				}
			}
		}
	})
	t.Run("search with empty query", func(t *testing.T) {
		query := form.NewLabelSearch("")
		result, err := Labels(query)

		if err != nil {
			t.Fatal(err)
		}

		t.Log(result)
		assert.LessOrEqual(t, 3, len(result))
	})
	t.Run("search with invalid query string", func(t *testing.T) {
		query := form.NewLabelSearch("xxx:bla")
		result, err := Labels(query)

		assert.Error(t, err, "unknown filter")
		assert.Empty(t, result)
	})
	t.Run("search for ID", func(t *testing.T) {
		f := form.SearchLabels{
			Query:    "",
			UID:      "ls6sg6b1wowuy3c4",
			Slug:     "",
			Name:     "",
			All:      false,
			Favorite: false,
			Count:    0,
			Offset:   0,
			Order:    "",
		}

		result, err := Labels(f)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "cake", result[0].LabelSlug)
	})
	t.Run("search for label landscape", func(t *testing.T) {
		f := form.SearchLabels{
			Query: "landscape",
		}

		result, err := Labels(f)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "flower", result[0].LabelSlug)
	})
}
