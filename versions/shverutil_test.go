package versions

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func TestUtilWerks(t *testing.T) {
	imgf, _ := os.Open("SH.png")
	img, err := png.Decode(imgf)
	if err != nil {
		panic(err)
	}
	img = etest(img)
	fmt.Println()
	fmt.Println(GetSHVersion(img))
	out, err := os.Create("result.png")
	png.Encode(out, img)
}

func etest(img image.Image) image.Image {
	e := &encodetest{}
	e.Image = img
	return e
}

type encodetest struct {
	image.Image
}

func (e *encodetest) At(x, y int) color.Color {
	if x >= e.Image.Bounds().Dx()-5 && y >= e.Image.Bounds().Dy()-5 {
		if x - (e.Image.Bounds().Dx()-5) == 3 && y - (e.Image.Bounds().Dy()-5) == 4 {
			return color.RGBA{0,0,0,255}
		}
		return color.RGBA{255, 255, 255, 255}
	}
	return e.Image.At(x, y)
}
