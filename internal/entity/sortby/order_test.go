package sortby

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderExpr(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		assert.Equal(t, "", OrderExpr("", false, MySQL))
		assert.Equal(t, "photos.edited_at", OrderExpr("photos.edited_at", false, MySQL))
		assert.Equal(t, "photos.edited_at ASC", OrderExpr("photos.edited_at ASC", false, MySQL))
		assert.Equal(t, "photos.edited_at DESC, files.media_id", OrderExpr("photos.edited_at DESC, files.media_id", false, MySQL))
		assert.Equal(t, "photos.edited_at DESC, files.media_id ASC", OrderExpr("photos.edited_at DESC, files.media_id ASC", false, MySQL))
		assert.Equal(t, "photo_count DESC NULLS LAST, albums.album_title, albums.album_uid DESC", OrderExpr("photo_count DESC NULLS LAST, albums.album_title, albums.album_uid DESC", false, MySQL))
	})
	t.Run("Reverse", func(t *testing.T) {
		assert.Equal(t, "", OrderExpr("", true, MySQL))
		assert.Equal(t, "photos.edited_at", OrderExpr("photos.edited_at", true, MySQL))
		assert.Equal(t, "photos.edited_at DESC", OrderExpr("photos.edited_at ASC", true, MySQL))
		assert.Equal(t, "photos.edited_at ASC, files.media_id", OrderExpr("photos.edited_at DESC, files.media_id", true, MySQL))
		assert.Equal(t, "photos.edited_at ASC, files.media_id DESC", OrderExpr("photos.edited_at DESC, files.media_id ASC", true, MySQL))
		assert.Equal(t, "photo_count ASC NULLS FIRST, albums.album_title, albums.album_uid ASC", OrderExpr("photo_count DESC NULLS LAST, albums.album_title, albums.album_uid DESC", true, MySQL))
	})
	t.Run("DefaultPostgreSQL", func(t *testing.T) {
		assert.Equal(t, "", OrderExpr("", false, Postgres))
		assert.Equal(t, "photos.photo_title COLLATE \"caseinsensitive\"", OrderExpr("photos.photo_title", false, Postgres))
		assert.Equal(t, "photos.photo_title COLLATE \"caseinsensitive\" ASC", OrderExpr("photos.photo_title ASC", false, Postgres))
		assert.Equal(t, "photos.photo_title COLLATE \"caseinsensitive\" DESC, albums.album_category COLLATE \"caseinsensitive\"", OrderExpr("photos.photo_title DESC, albums.album_category", false, Postgres))
		assert.Equal(t, "photos.photo_title COLLATE \"caseinsensitive\" DESC, albums.album_category COLLATE \"caseinsensitive\" ASC", OrderExpr("photos.photo_title DESC, albums.album_category ASC", false, Postgres))
		assert.Equal(t, "photo_count DESC NULLS LAST, albums.album_title COLLATE \"caseinsensitive\", albums.album_uid DESC", OrderExpr("photo_count DESC NULLS LAST, albums.album_title, albums.album_uid DESC", false, Postgres))
	})
	t.Run("ReversePostgreSQL", func(t *testing.T) {
		assert.Equal(t, "", OrderExpr("", true, Postgres))
		assert.Equal(t, "photos.photo_title COLLATE \"caseinsensitive\"", OrderExpr("photos.photo_title", true, Postgres))
		assert.Equal(t, "photos.photo_title COLLATE \"caseinsensitive\" DESC", OrderExpr("photos.photo_title ASC", true, Postgres))
		assert.Equal(t, "photos.photo_title COLLATE \"caseinsensitive\" ASC, albums.album_category COLLATE \"caseinsensitive\"", OrderExpr("photos.photo_title DESC, albums.album_category", true, Postgres))
		assert.Equal(t, "photos.photo_title COLLATE \"caseinsensitive\" ASC, albums.album_category COLLATE \"caseinsensitive\" DESC", OrderExpr("photos.photo_title DESC, albums.album_category ASC", true, Postgres))
		assert.Equal(t, "photo_count ASC NULLS FIRST, albums.album_title COLLATE \"caseinsensitive\", albums.album_uid ASC", OrderExpr("photo_count DESC NULLS LAST, albums.album_title, albums.album_uid DESC", true, Postgres))
	})
}

func TestOrderAsc(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		assert.Equal(t, DirAsc, OrderAsc(false))
	})
	t.Run("Reverse", func(t *testing.T) {
		assert.Equal(t, DirDesc, OrderAsc(true))
	})
}

func TestOrderDesc(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		assert.Equal(t, DirDesc, OrderDesc(false))
	})
	t.Run("Reverse", func(t *testing.T) {
		assert.Equal(t, DirAsc, OrderDesc(true))
	})
}
