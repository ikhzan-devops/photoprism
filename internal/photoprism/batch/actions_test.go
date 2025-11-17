package batch

import (
	"testing"

	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/entity/query"
)

// TestApplyAlbums exercises batch action logic.
func TestApplyAlbums(t *testing.T) {
	t.Run("AddPhotoToExistingAlbumByUID", func(t *testing.T) {
		photo := entity.PhotoFixtures.Get("Photo01")
		albumUID := entity.AlbumFixtures.Get("christmas2030").AlbumUID

		albums := Items{
			Items: []Item{
				{Action: ActionAdd, Value: albumUID},
			},
		}

		if err := ApplyAlbums(photo.PhotoUID, albums); err != nil {
			t.Fatal(err)
		}

		// Verify photo was added to album by checking photos_albums table
		var photoAlbum entity.PhotoAlbum
		if err := entity.Db().Where("album_uid = ? AND photo_uid = ? AND hidden = ?", albumUID, photo.PhotoUID, false).First(&photoAlbum).Error; err != nil {
			t.Fatal(err)
		}
	})
	t.Run("AddPhotoToNewAlbumByTitle", func(t *testing.T) {
		photo := entity.PhotoFixtures.Get("Photo02")
		albumTitle := "Test Album for Actions"

		albums := Items{
			Items: []Item{
				{Action: ActionAdd, Title: albumTitle},
			},
		}

		if err := ApplyAlbums(photo.PhotoUID, albums); err != nil {
			t.Fatal(err)
		}

		// Verify album was created and photo added
		var album entity.Album
		if err := entity.Db().Where("album_title = ?", albumTitle).First(&album).Error; err != nil {
			t.Fatal(err)
		}

		if album.AlbumTitle != albumTitle {
			t.Errorf("expected album title %s, got %s", albumTitle, album.AlbumTitle)
		}

		var photoAlbum entity.PhotoAlbum
		if err := entity.Db().Where("album_uid = ? AND photo_uid = ? AND hidden = ?", album.AlbumUID, photo.PhotoUID, false).First(&photoAlbum).Error; err != nil {
			t.Fatal(err)
		}
	})
	t.Run("RemovePhotoFromAlbum", func(t *testing.T) {
		// First add photo to album
		photo := entity.PhotoFixtures.Get("Photo03")
		album := entity.AlbumFixtures.Get("holiday-2030")

		// Create photo-album relation manually
		var existing entity.PhotoAlbum
		err := entity.Db().Where("album_uid = ? AND photo_uid = ?", album.AlbumUID, photo.PhotoUID).First(&existing).Error
		if err != nil {
			photoAlbumEntry := entity.PhotoAlbum{
				PhotoUID: photo.PhotoUID,
				AlbumUID: album.AlbumUID,
				Hidden:   false,
			}

			if err := entity.Db().Create(&photoAlbumEntry).Error; err != nil {
				t.Fatal(err)
			}
		} else if existing.Hidden {
			existing.Hidden = false
			if err := entity.Db().Save(&existing).Error; err != nil {
				t.Fatal(err)
			}
		}

		// Verify photo is in album
		var checkEntry entity.PhotoAlbum
		if err := entity.Db().Where("album_uid = ? AND photo_uid = ? AND hidden = ?", album.AlbumUID, photo.PhotoUID, false).First(&checkEntry).Error; err != nil {
			t.Fatal(err)
		}

		// Now remove it
		albums := Items{
			Items: []Item{
				{Action: ActionRemove, Value: album.AlbumUID},
			},
		}

		if err := ApplyAlbums(photo.PhotoUID, albums); err != nil {
			t.Fatal(err)
		}

		// Verify photo was removed (should be marked as hidden)
		var removedEntry entity.PhotoAlbum
		if err := entity.Db().Where("album_uid = ? AND photo_uid = ?", album.AlbumUID, photo.PhotoUID).First(&removedEntry).Error; err != nil {
			t.Fatal(err)
		}

		if !removedEntry.Hidden {
			t.Error("expected photo to be marked as hidden in album")
		}
	})
	// Error cases
	t.Run("AddPhotoToNonExistingAlbumByUID", func(t *testing.T) {
		photo := entity.PhotoFixtures.Get("Photo04")
		nonExistingAlbumUID := "at9lxuqxpoaaaaaa" // Invalid/non-existing UID

		albums := Items{
			Items: []Item{
				{Action: ActionAdd, Value: nonExistingAlbumUID},
			},
		}

		err := ApplyAlbums(photo.PhotoUID, albums)
		if err == nil {
			t.Error("expected error when adding photo to non-existing album, but got none")
		}
	})
	t.Run("AddPhotoToAlbumWithInvalidUID", func(t *testing.T) {
		photo := entity.PhotoFixtures.Get("Photo04")
		invalidUID := "invalid-uid-format" // Invalid UID format

		albums := Items{
			Items: []Item{
				{Action: ActionAdd, Value: invalidUID},
			},
		}

		err := ApplyAlbums(photo.PhotoUID, albums)
		if err == nil {
			t.Error("expected error when adding photo to album with invalid UID, but got none")
		}
	})
	t.Run("RemovePhotoFromNonExistingAlbum", func(t *testing.T) {
		photo := entity.PhotoFixtures.Get("Photo05")
		nonExistingAlbumUID := "at9lxuqxpobbbbbb" // Non-existing UID

		albums := Items{
			Items: []Item{
				{Action: ActionRemove, Value: nonExistingAlbumUID},
			},
		}

		err := ApplyAlbums(photo.PhotoUID, albums)
		if err == nil {
			t.Error("expected error when removing photo from non-existing album, but got none")
		}
	})
	t.Run("InvalidActionOnAlbum", func(t *testing.T) {
		photo := entity.PhotoFixtures.Get("Photo06")
		albumUID := entity.AlbumFixtures.Get("christmas2030").AlbumUID

		albums := Items{
			Items: []Item{
				{Action: "invalid-action", Value: albumUID}, // Invalid action
			},
		}

		err := ApplyAlbums(photo.PhotoUID, albums)
		if err == nil {
			t.Error("expected error for invalid action, but got none")
		}
	})
	t.Run("EmptyAlbumItems", func(t *testing.T) {
		photo := entity.PhotoFixtures.Get("Photo07")

		albums := Items{
			Items: []Item{}, // Empty items
		}

		// This should not error, but should be a no-op
		err := ApplyAlbums(photo.PhotoUID, albums)
		if err != nil {
			t.Errorf("expected no error for empty album items, but got: %v", err)
		}
	})
	t.Run("AddPhotoToAlbumWithEmptyValueAndTitle", func(t *testing.T) {
		photo := entity.PhotoFixtures.Get("Photo08")

		albums := Items{
			Items: []Item{
				{Action: ActionAdd, Value: "", Title: ""}, // Both empty
			},
		}

		err := ApplyAlbums(photo.PhotoUID, albums)
		if err == nil {
			t.Error("expected error when both Value and Title are empty, but got none")
		}
	})
	t.Run("InvalidPhotoUID", func(t *testing.T) {
		invalidPhotoUID := "invalid-photo-uid"
		albumUID := entity.AlbumFixtures.Get("christmas2030").AlbumUID

		albums := Items{
			Items: []Item{
				{Action: ActionAdd, Value: albumUID},
			},
		}

		if err := ApplyAlbums(invalidPhotoUID, albums); err == nil {
			t.Error("expected error for invalid photo UID, but got none")
		}
	})
}

// TestApplyLabels exercises batch action logic.
func TestApplyLabels(t *testing.T) {
	t.Run("AddExistingLabelByUID", func(t *testing.T) {
		photo := entity.PhotoFixtures.Pointer("Photo06")
		labelUID := entity.LabelFixtures.Get("landscape").LabelUID

		labels := Items{
			Items: []Item{
				{Action: ActionAdd, Value: labelUID},
			},
		}

		if err := ApplyLabels(photo, labels); err != nil {
			t.Fatal(err)
		}

		// Verify photo has the label with 100% confidence and batch source
		var photoLabel entity.PhotoLabel

		if err := entity.Db().Preload("Label").Where("photo_id = ? AND label_id = (SELECT id FROM labels WHERE label_uid = ?)", photo.ID, labelUID).First(&photoLabel).Error; err != nil {
			t.Fatal(err)
		}

		if photoLabel.Uncertainty != 0 {
			t.Errorf("expected uncertainty 0 (100%% confidence), got %d", photoLabel.Uncertainty)
		}

		if photoLabel.LabelSrc != entity.SrcBatch {
			t.Errorf("expected label source %s, got %s", entity.SrcBatch, photoLabel.LabelSrc)
		}
	})
	t.Run("AddNewLabelByTitle", func(t *testing.T) {
		photo := entity.PhotoFixtures.Pointer("Photo07")
		labelTitle := "Test Label for Actions"

		labels := Items{
			Items: []Item{
				{Action: ActionAdd, Title: labelTitle},
			},
		}

		if err := ApplyLabels(photo, labels); err != nil {
			t.Fatal(err)
		}

		// Verify label was created and added to photo
		var label entity.Label

		if err := entity.Db().Where("label_name = ?", labelTitle).First(&label).Error; err != nil {
			t.Fatal(err)
		}

		if label.LabelName != labelTitle {
			t.Errorf("expected label name %s, got %s", labelTitle, label.LabelName)
		}

		var photoLabel entity.PhotoLabel

		if err := entity.Db().Where("photo_id = ? AND label_id = ?", photo.ID, label.ID).First(&photoLabel).Error; err != nil {
			t.Fatal(err)
		}

		if photoLabel.Uncertainty != 0 {
			t.Errorf("expected uncertainty 0, got %d", photoLabel.Uncertainty)
		}

		if photoLabel.LabelSrc != entity.SrcBatch {
			t.Errorf("expected label source %s, got %s", entity.SrcBatch, photoLabel.LabelSrc)
		}
	})
	t.Run("RemoveLabelByUID", func(t *testing.T) {
		photo := entity.PhotoFixtures.Pointer("Photo08")
		label := entity.LabelFixtures.Get("flower")

		// First add the label manually
		photoLabel := entity.NewPhotoLabel(photo.ID, label.ID, 0, entity.SrcBatch)

		if err := entity.Db().Create(&photoLabel).Error; err != nil {
			t.Fatal(err)
		}

		// Verify label is on photo
		var checkPhotoLabel entity.PhotoLabel

		if err := entity.Db().Where("photo_id = ? AND label_id = ?", photo.ID, label.ID).First(&checkPhotoLabel).Error; err != nil {
			t.Fatal(err)
		}

		// Now remove it
		labels := Items{
			Items: []Item{
				{Action: ActionRemove, Value: label.LabelUID},
			},
		}

		if err := ApplyLabels(photo, labels); err != nil {
			t.Fatal(err)
		}

		// Verify label was removed (should be deleted from photos_labels)
		var deletedLabel entity.PhotoLabel

		if err := entity.Db().Where("photo_id = ? AND label_id = ?", photo.ID, label.ID).First(&deletedLabel).Error; err == nil {
			t.Error("expected label to be deleted, but it was found")
		}
	})
	t.Run("RemoveLabelWithPreloadedPhoto", func(t *testing.T) {
		photo := entity.PhotoFixtures.Pointer("Photo17")
		label := entity.LabelFixtures.Get("landscape")

		entity.Db().Where("photo_id = ? AND label_id = ?", photo.ID, label.ID).Delete(&entity.PhotoLabel{})

		if err := entity.Db().Create(entity.NewPhotoLabel(photo.ID, label.ID, 0, entity.SrcManual)).Error; err != nil {
			t.Fatal(err)
		}

		preloaded, err := query.PhotoPreloadByUID(photo.PhotoUID)

		if err != nil {
			t.Fatal(err)
		}

		if len(preloaded.Labels) == 0 {
			t.Fatalf("expected preloaded labels for %s", photo.PhotoUID)
		}

		labels := Items{
			Items: []Item{{Action: ActionRemove, Value: label.LabelUID}},
		}

		if err = ApplyLabels(&preloaded, labels); err != nil {
			t.Fatal(err)
		}

		var deleted entity.PhotoLabel

		if err = entity.Db().Where("photo_id = ? AND label_id = ?", photo.ID, label.ID).First(&deleted).Error; err == nil {
			t.Fatal("expected label relation to be removed")
		}
	})
	t.Run("RemoveAutoLabelSetsUncertaintyTo100", func(t *testing.T) {
		photo := entity.PhotoFixtures.Pointer("Photo09")
		label := entity.LabelFixtures.Get("cake")

		// Add label with auto source (not manual/batch)
		photoLabel := entity.NewPhotoLabel(photo.ID, label.ID, 15, entity.SrcImage)

		if err := entity.Db().Create(&photoLabel).Error; err != nil {
			t.Fatal(err)
		}

		labels := Items{
			Items: []Item{
				{Action: ActionRemove, Value: label.LabelUID},
			},
		}

		if err := ApplyLabels(photo, labels); err != nil {
			t.Fatal(err)
		}

		// Verify label uncertainty was set to 100% (blocked)
		var blockedLabel entity.PhotoLabel

		if err := entity.Db().Where("photo_id = ? AND label_id = ?", photo.ID, label.ID).First(&blockedLabel).Error; err != nil {
			t.Fatal(err)
		}

		if blockedLabel.Uncertainty != 100 {
			t.Errorf("expected uncertainty 100 (blocked), got %d", blockedLabel.Uncertainty)
		}

		if blockedLabel.LabelSrc != entity.SrcBatch {
			t.Errorf("expected label source %s, got %s", entity.SrcBatch, blockedLabel.LabelSrc)
		}
	})
	t.Run("UpdateExistingLabelConfidence", func(t *testing.T) {
		photo := entity.PhotoFixtures.Pointer("Photo10")
		label := entity.LabelFixtures.Get("landscape")

		// First, delete any existing photo-label relation to ensure clean test
		entity.Db().Where("photo_id = ? AND label_id = ?", photo.ID, label.ID).Delete(&entity.PhotoLabel{})

		// Add label with some uncertainty using FirstOrCreatePhotoLabel
		photoLabel := entity.FirstOrCreatePhotoLabel(entity.NewPhotoLabel(photo.ID, label.ID, 50, entity.SrcImage))

		if photoLabel == nil {
			t.Fatal("failed to create photo label")
		}

		// Verify initial state
		if photoLabel.Uncertainty != 50 {
			t.Errorf("expected uncertainty 50, got %d", photoLabel.Uncertainty)
		}

		if photoLabel.LabelSrc != entity.SrcImage {
			t.Errorf("expected label source %s, got %s", entity.SrcImage, photoLabel.LabelSrc)
		}

		// Re-add same label via batch (should update to 100% confidence)
		labels := Items{
			Items: []Item{
				{Action: ActionAdd, Value: label.LabelUID},
			},
		}

		if err := ApplyLabels(photo, labels); err != nil {
			t.Fatal(err)
		}

		// Verify label confidence was updated
		var updatedLabel entity.PhotoLabel

		if err := entity.Db().Where("photo_id = ? AND label_id = ?", photo.ID, label.ID).First(&updatedLabel).Error; err != nil {
			t.Fatal(err)
		}

		if updatedLabel.Uncertainty != 0 {
			t.Errorf("expected uncertainty 0 (100%% confidence), got %d", updatedLabel.Uncertainty)
		}

		if updatedLabel.LabelSrc != entity.SrcBatch {
			t.Errorf("expected label source %s, got %s", entity.SrcBatch, updatedLabel.LabelSrc)
		}
	})
	t.Run("InvalidPhotoReturnsError", func(t *testing.T) {
		labels := Items{
			Items: []Item{
				{Action: ActionAdd, Value: "some-uid"},
			},
		}

		err := ApplyLabels(nil, labels)

		if err == nil {
			t.Error("expected error for nil photo")
		}

		emptyPhoto := &entity.Photo{}
		err = ApplyLabels(emptyPhoto, labels)

		if err == nil {
			t.Error("expected error for empty photo")
		}
	})
	// Additional error cases
	t.Run("AddNonExistingLabelByUID", func(t *testing.T) {
		photo := entity.PhotoFixtures.Pointer("Photo11")
		nonExistingLabelUID := "lt9lxuqxpoaaaaaa" // Invalid/non-existing UID

		labels := Items{
			Items: []Item{
				{Action: ActionAdd, Value: nonExistingLabelUID},
			},
		}

		err := ApplyLabels(photo, labels)

		if err == nil {
			t.Error("expected error when adding non-existing label by UID, but got none")
		}
	})
	t.Run("AddLabelWithInvalidUID", func(t *testing.T) {
		photo := entity.PhotoFixtures.Pointer("Photo12")
		invalidUID := "invalid-label-uid-format" // Invalid UID format

		labels := Items{
			Items: []Item{
				{Action: ActionAdd, Value: invalidUID},
			},
		}

		err := ApplyLabels(photo, labels)

		if err == nil {
			t.Error("expected error when adding label with invalid UID, but got none")
		}
	})
	t.Run("RemoveNonExistingLabelByUID", func(t *testing.T) {
		photo := entity.PhotoFixtures.Pointer("Photo13")
		nonExistingLabelUID := "lt9lxuqxpobbbbbb" // Non-existing UID

		labels := Items{
			Items: []Item{
				{Action: ActionRemove, Value: nonExistingLabelUID},
			},
		}

		err := ApplyLabels(photo, labels)
		if err == nil {
			t.Error("expected error when removing non-existing label, but got none")
		}
	})
	t.Run("InvalidActionOnLabel", func(t *testing.T) {
		photo := entity.PhotoFixtures.Pointer("Photo14")
		labelUID := entity.LabelFixtures.Get("landscape").LabelUID

		labels := Items{
			Items: []Item{
				{Action: "invalid-action", Value: labelUID}, // Invalid action
			},
		}

		err := ApplyLabels(photo, labels)
		if err == nil {
			t.Error("expected error for invalid action, but got none")
		}
	})
	t.Run("EmptyLabelItems", func(t *testing.T) {
		photo := entity.PhotoFixtures.Pointer("Photo15")

		labels := Items{
			Items: []Item{}, // Empty items
		}

		// This should not error, but should be a no-op
		err := ApplyLabels(photo, labels)

		if err != nil {
			t.Errorf("expected no error for empty label items, but got: %v", err)
		}
	})
	t.Run("AddLabelWithEmptyValueAndTitle", func(t *testing.T) {
		photo := entity.PhotoFixtures.Pointer("Photo16")

		labels := Items{
			Items: []Item{
				{Action: ActionAdd, Value: "", Title: ""}, // Both empty
			},
		}

		err := ApplyLabels(photo, labels)

		if err == nil {
			t.Error("expected error when both Value and Title are empty, but got none")
		}
	})
	t.Run("RemoveLabelNotAssignedToPhoto", func(t *testing.T) {
		photo := entity.PhotoFixtures.Pointer("Photo17")
		labelUID := entity.LabelFixtures.Get("bird").LabelUID

		// Ensure the label is not assigned to this photo
		entity.Db().Where("photo_id = ? AND label_id = (SELECT id FROM labels WHERE label_uid = ?)", photo.ID, labelUID).Delete(&entity.PhotoLabel{})

		labels := Items{
			Items: []Item{
				{Action: ActionRemove, Value: labelUID},
			},
		}

		err := ApplyLabels(photo, labels)

		if err == nil {
			t.Error("expected error when removing label not assigned to photo, but got none")
		}
	})
}

// TestIndexPhotoLabels exercises batch action logic.
func TestIndexPhotoLabels(t *testing.T) {
	labels := entity.PhotoLabels{
		{LabelID: 11},
		{LabelID: 0},
		{LabelID: 22},
	}

	idx := indexPhotoLabels(labels)

	if len(idx) != 2 {
		t.Fatalf("expected 2 labels in index, got %d", len(idx))
	}

	if idx[11] == nil || idx[22] == nil {
		t.Fatal("expected indexed labels to be present")
	}
}
