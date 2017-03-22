package Belt

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	"../versions"
	"github.com/bwmarrin/discordgo"
)

func swap(broadcast bool, version bool) {
	var img image.Image
	if broadcast {
		f, err := os.Create("unedited.png")
		defer f.Close()
		HE(err)
		img = getav(sh.OwnID)
		png.Encode(f, img)
		img = EncodeBasic(img)
		if version {
			img = EncodeVersion(img, sh.version)
		}
		changeav(img)
	} else if !broadcast {
		b, err := os.Open("unedited.png")
		HE(err)
		img, _, err := image.Decode(b)
		HE(err)
		changeav(img)
	}
}

func getav(ID string) image.Image {
	var us *discordgo.User
	for _, g := range sh.dg.State.Guilds {
		u := GetUser(ID, g)
		if u.ID == ID {
			us = u
			break
		}
	}
	if us.Avatar == "" {
		return nil
	}
	resp, err := http.Get(discordgo.EndpointUserAvatar(us.ID, us.Avatar))
	HE(err)
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	HE(err)
	return img
}

func changeav(img image.Image) {
	w := new(bytes.Buffer)
	jpeg.Encode(w, img, &jpeg.Options{Quality: 100})
	i, err := ioutil.ReadAll(w)

	HE(err)
	base64 := base64.StdEncoding.EncodeToString(i)

	avatar := fmt.Sprintf("data:%s;base64,%s", http.DetectContentType(i), base64)

	_, err = sh.dg.UserUpdate("", "", sh.OwnName, avatar, "")
	HE(err)
}

func sweep() map[string]versions.Version {
	var ids = make(map[string]versions.Version)

	ch := make(chan map[string]versions.Version, len(sh.dg.State.Guilds))

	for _, g := range sh.dg.State.Guilds {

		go func(ids map[string]versions.Version, g *discordgo.Guild, ch chan map[string]versions.Version) {

			c := make(chan map[string]versions.Version, len(g.Members))

			for _, m := range g.Members {
				go func(ids map[string]versions.Version, m *discordgo.Member, c chan map[string]versions.Version) {
					ok, ver := identify(m)
					if ok {
						ids[m.User.ID] = ver
					}
					c <- ids

				}(ids, m, c)
			}

			for i := 0; i < len(g.Members); i++ {
				for s, ver := range <-c {
					ids[s] = ver
				}
			}
			ch <- ids

		}(ids, g, ch)

	}

	for i := 0; i < len(sh.dg.State.Guilds); i++ {
		for s, ver := range <-ch {
			ids[s] = ver
		}
	}

	fmt.Println(ids)
	runtime.GC()

	return ids
}

func identify(m *discordgo.Member) (bool, versions.Version) {
	if m.User.Avatar == "" {
		return false, versions.Version{}
	}
	return DecodeUrl(discordgo.EndpointUserAvatar(m.User.ID, m.User.Avatar))

}
