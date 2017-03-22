package Belt

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"testing"

	"../versions"
)

const this = "../belter/"

func TestEncodefile2(t *testing.T) {
	EncodefileBasic(this+"sherlock.jpg", this+"Encoded.png")
}

func TestEncodefileVersion(t *testing.T) {
	EncodefileVersion(this+"sherlock.jpg", this+"SH.png", versions.Version{3, 5, 12, 15})
}

func TestDecodeFile(t *testing.T) {
	ver := DecodeFile(this + "SH.png")
	fmt.Println(ver)
}

func TestPushVersion(t *testing.T) {
	data, err := os.Open(this + "sherlock.jpg")
	HE(err)
	img, _, err := image.Decode(data)
	HE(err)
	img = EncodeBasic(img)
	img = EncodeVersion(img, versions.Version{0, 0, 2, 0})
	out, err := os.Create("Encoded.png")
	HE(err)
	png.Encode(out, img)
}

func TestDecodeUrl(t *testing.T) {
	gfi := make(chan bool)
	go Run(gfi)
	_ = <-gfi
	fmt.Println(DecodeUrl("https://discordapp.com/api/users/268000155825864705/avatars/7fcfe0f89f8494609ba6c2d984642d70.jpg"))
}

func TestVerPix(t *testing.T) {
	ver := versions.Version{0, 4, 8, 0}
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			fmt.Println(versions.GetVerPix(ver, x, y))
		}
	}
}

const here = "/Users/fokedejong-noorman/mygo/src/JohnWatson.bot/testfileserver/"

var localimage = "http://localhost:" + localport + "/SH.png"
var localport = "9003"

var serverup bool

func BenchmarkDecodeUrl(b *testing.B) {
	goforit := make(chan bool)
	if !serverup {
		go Run(goforit)
		if <-goforit == true {
			serverup = true
		} else {
			panic("AHHH")
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DecodeUrl(localimage)
	}
}

func Run(ch chan bool) {
	im := http.NewServeMux()
	im.HandleFunc("/", http.FileServer(http.Dir(here+"images")).ServeHTTP)
	fmt.Println("Starting image server...")
	ch <- true
	fmt.Println(http.Dir(here + "images"))
	log.Fatal(http.ListenAndServe(":"+localport, im))
}
