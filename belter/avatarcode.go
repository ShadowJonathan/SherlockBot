package Belt

import (
	"../versions"
	"image"
	"image/color"
	"image/png"
	"os"
	"net/http"
)


func HE(err error) {
	if err != nil {
		panic(err)
	}
}


func EncodefileBasic(file, dest string) {
	data, err := os.Open(file)
	HE(err)
	img, _, err := image.Decode(data)
	HE(err)
	img = EncodeBasic(img)
	out, err := os.Create(dest)
	png.Encode(out, img)
}

func EncodeBasic(img image.Image) image.Image {
	e := &basicencodedimage{}
	e.Image = img
	return e
}

type basicencodedimage struct {
	image.Image
}

func (m *basicencodedimage) At(x, y int) color.Color {
	if x <= 5 && y <= 5 {
		return color.RGBA{255, 255, 255, 255}
	}
	return m.Image.At(x, y)
}

func EncodefileVersion(file, dest string, ver versions.Version) {
	data, err := os.Open(file)
	HE(err)
	img, _, err := image.Decode(data)
	HE(err)
	img = EncodeVersion(img, ver)
	out, err := os.Create(dest)
	png.Encode(out, img)
}

func EncodeVersion(img image.Image, version versions.Version) image.Image {
	e := &versionencodedimage{}
	e.Ver = version
	e.Image = img
	return e
}

type versionencodedimage struct {
	image.Image
	Ver versions.Version
}

func (m *versionencodedimage) At(x, y int) color.Color {
	if x >= m.Bounds().Dx()-5 && y >= m.Bounds().Dy()-5 {
		rx := x - (m.Bounds().Dx() - 5)
		ry := y - (m.Bounds().Dy() - 5)
		return versions.GetVerPix(m.Ver,rx,ry)
	}
	return m.Image.At(x, y)
}

func DecodeFile(file string) versions.Version {
	data, err := os.Open(file)
	HE(err)
	img, _, err := image.Decode(data)
	HE(err)
	return Decode(img)
}

func Decode(img image.Image) versions.Version {
	ver := versions.GetSHVersion(img)
	return ver
}

func DecodeUrl(url string) (bool,versions.Version) {
	resp, err := http.Get(url)
	HE(err)
	body := resp.Body
	defer body.Close()
	img, _,err := image.Decode(body)
	HE(err)
	ver := Decode(img)
	nilver := versions.Version{}
	if ver != nilver {
		return true, ver
	} else {
		return false, ver

	}
}