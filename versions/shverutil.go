package versions

import (
	"image"
	"image/color"
)

func GetSHVersion(img image.Image) Version {
	colorarray := [25]color.Color{}
	var ai int
	ai = 0
	for i := img.Bounds().Dy() - 5; i < img.Bounds().Dy(); i++ {
		for I := img.Bounds().Dx() - 5; I < img.Bounds().Dx(); I++ {
			colorarray[ai] = img.At(I, i)
			ai++
		}
	}
	fail := Version{0, 0, 0, 0}
	header := colorarray[:5]
	Major := colorarray[5:10]
	Minor := colorarray[10:15]
	Build := colorarray[15:20]
	Exper := colorarray[20:25]
	for _, c := range header {
		if (c != white) && (WB(c) != white) {
			return fail
		}
	}
	return Version{
		convert(Major),
		convert(Minor),
		convert(Build),
		convert(Exper),
	}
}

func convert(c []color.Color) int {
	if len(c) > 5 || len(c) < 5 {
		panic(c)
	}
	var w int = 0
	if c[4] == black || WB(c[4]) == black {
		w = w + 1
	}
	if c[3] == black || WB(c[3]) == black {
		w = w + 2
	}
	if c[2] == black || WB(c[2]) == black {
		w = w + 4
	}
	if c[1] == black || WB(c[1]) == black {
		w = w + 8
	}
	return w
}

func Upperleftwhite(img image.Image) bool {
	colorarray := [25]color.Color{}
	var ai int
	ai = 0
	for i := 0; i < 5; i++ {
		for I := 0; I < 5; I++ {
			colorarray[ai] = img.At(I, i)
			ai++
		}
	}

	for _, c := range colorarray {
		if (c != white) && (WB(c) != white) {
			return false
		}
	}
	return true
}
