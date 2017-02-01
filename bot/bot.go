package main

import (
	"io/ioutil"
	"strings"

	"fmt"

	"strconv"

	"../belter"
)

func main() {
	token, err := ioutil.ReadFile("../token")
	if err != nil {
		fmt.Println("Error reading token file: " + err.Error())
		err := ioutil.WriteFile("../token", token, 9000)
		if err != nil {
			fmt.Println("Error writing sample token file: " + err.Error())
		}
		return
	}
	restart, upgrade := Belt.Initialize(strings.TrimSpace(string(token)))
	ioutil.WriteFile("../retcmd.botboot", compilebotboot(restart, upgrade), 0777)
}

func compilebotboot(restart, upgrade bool) []byte {
	var BB []string
	BB[0] = strconv.FormatBool(restart)
	BB = append(BB, strconv.FormatBool(upgrade))
	S := strings.Join(BB, " ")
	return []byte(S)
}
