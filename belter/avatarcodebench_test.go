package Belt

import (
	"image"
	"os"
	"testing"
)

func BenchmarkDecodeFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DecodeFile("SH.png")
	}
}

func BenchmarkDecode2(b *testing.B) {
	data, err := os.Open("SH.png")
	HE(err)
	img, _, err := image.Decode(data)
	HE(err)
	for i := 0; i < b.N; i++ {
		Decode(img)
	}
}
