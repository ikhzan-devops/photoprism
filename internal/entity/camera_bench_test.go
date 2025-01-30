package entity

import (
	"testing"
)

func BenchmarkCreateDeleteCamera(b *testing.B) {
	for interations := 0; interations < b.N; interations++ {
		camera := NewCamera("Palasonic", "Palasonic Dumix")

		if err := camera.Create(); err != nil {
			b.Fatal(err)
		}
		cameraCache.Delete(camera.CameraSlug)
		if err := UnscopedDb().Delete(camera).Error; err != nil {
			b.Fatal(err)
		}
	}
}
