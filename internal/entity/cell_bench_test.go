package entity

import (
	"math/rand/v2"
	"strings"
	"testing"
)

func randRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

const characterRunes = " abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const sha1Runes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(len int) string {
	sb := strings.Builder{}
	sb.Grow(len)
	for i := 0; i < len; {
		sb.WriteByte(characterRunes[rand.IntN(53)])
		i++
	}

	return sb.String()
}

func randomSHA1() string {
	sb := strings.Builder{}
	sb.Grow(32)
	for i := 0; i < 32; {
		sb.WriteByte(characterRunes[rand.IntN(36)])
		i++
	}
	return sb.String()
}

func BenchmarkCreateDeleteCellAndPlace(b *testing.B) {
	for interations := 0; interations < b.N; interations++ {
		lat := randRange(-90, 90)
		lng := randRange(-180, 180)
		cell := NewCell(lat, lng)
		place := &Place{
			ID:            randomString(12),
			PlaceLabel:    randomString(20),
			PlaceDistrict: randomString(30),
			PlaceCity:     randomString(30),
			PlaceState:    randomString(30),
			PlaceCountry:  randomString(2),
			PlaceKeywords: randomString(10),
			PlaceFavorite: false,
		}

		if cell.Place = FirstOrCreatePlace(place); cell.Place == nil {
			b.Fatal("unable to find/create place")
		}

		cell.PlaceID = cell.Place.ID

		if FirstOrCreateCell(cell) == nil {
			b.Fatal("unable to find/create cell")
		}
		if err := cell.Delete(); err != nil {
			b.Fatal(err)
		}

		if err := place.Delete(); err != nil {
			b.Fatal(err)
		}
	}
}
