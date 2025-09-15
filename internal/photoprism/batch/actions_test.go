package batch

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/photoprism/photoprism/internal/entity"
)

func TestMain(m *testing.M) {
	log = logrus.StandardLogger()
	log.SetLevel(logrus.TraceLevel)

	db := entity.InitTestDb(
		os.Getenv("PHOTOPRISM_TEST_DRIVER"),
		os.Getenv("PHOTOPRISM_TEST_DSN"))

	defer db.Close()

	code := m.Run()

	os.Exit(code)
}

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
}

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
		err := entity.Db().Where("photo_id = ? AND label_id = ?", photo.ID, label.ID).First(&deletedLabel).Error
		if err == nil {
			t.Error("expected label to be deleted, but it was found")
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
}
