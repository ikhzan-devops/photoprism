package entity

import (
	"testing"
)

func BenchmarkCreateDeleteAlbum(b *testing.B) {
	for interations := 0; interations < b.N; interations++ {
		album := NewAlbum("BenchMarkAlbum", AlbumManual)
		if err := album.Create(); err != nil {
			b.Fatal(err)
		}
		if err := album.DeletePermanently(); err != nil {
			b.Fatal(err)
		}
	}
}
