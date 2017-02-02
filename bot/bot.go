package main

import (
	"io/ioutil"
	"strings"

	"fmt"

	"strconv"

	"../belter"
)

func main() {
	for {
		token, err := ioutil.ReadFile("../token")
		if err != nil {
			fmt.Println("Error reading token file: " + err.Error())
			err := ioutil.WriteFile("../token", token, 9000)
			if err != nil {
				fmt.Println("Error writing sample token file: " + err.Error())
			}
			ioutil.WriteFile("../retcmd.botboot", compilebotboot(false), 0777)
			return
		}
		
		restart, upgrade := Belt.Initialize(strings.TrimSpace(string(token)))
		
		if !restart && !upgrade {
			ioutil.WriteFile("../retcmd.botboot", compilebotboot(upgrade), 0777)
			return
		}
		if restart && !upgrade {
		} else if upgrade {
			ioutil.WriteFile("../retcmd.botboot", compilebotboot(upgrade), 0777)
			return
		}
	}
}

func compilebotboot(upgrade bool) []byte {
	return []byte(strconv.FormatBool(upgrade))
}
