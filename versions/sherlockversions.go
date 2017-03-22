package versions

import (
	"image/color"
)

type Version struct {
	Major        int
	Minor        int
	Build        int
	Experimental int
}

var white = color.RGBA{255, 255, 255, 255}
var black = color.RGBA{0, 0, 0, 255}

func WB(c color.Color) color.Color {
	var bhalf uint32 = 65535 / 10
	var whalf uint32 = 65535 - bhalf
	r, g, b, _ := c.RGBA()
	if r < bhalf || g < bhalf || b < bhalf {
		return black
	} else if r > whalf || g > whalf || b > whalf {
		return white
	}
	return color.RGBA{160, 160, 160, 160}
}

func GetVerPix(ver Version, x int, y int) color.Color {
	go verify(ver)
	if x == 0 || y == 0 {
		return white
	}
	if y == 1 {
		return yesno(boolarrayfromint(ver.Major)[x-1])
	}
	if y == 2 {
		return yesno(boolarrayfromint(ver.Minor)[x-1])
	}
	if y == 3 {
		return yesno(boolarrayfromint(ver.Build)[x-1])
	}
	if y == 4 {
		return yesno(boolarrayfromint(ver.Experimental)[x-1])
	}
	return white
}

func yesno(b bool) color.Color {
	if b {
		return black
	} else {
		return white
	}
}

func boolarrayfromint(i int) []bool {
	var dataBitmap = make([]bool, 4)
	var index uint64 = 0
	for index < 4 {
		dataBitmap[index] = i&(1<<index) > 0
		index++
	}
	var b = make([]bool, 4)
	b[0] = dataBitmap[3]
	b[1] = dataBitmap[2]
	b[2] = dataBitmap[1]
	b[3] = dataBitmap[0]
	return b
}

func verify(version Version) {
	if version.Major > 15 || version.Minor > 15 || version.Build > 15 || version.Experimental > 15 {
		panic(version)
	}
}
