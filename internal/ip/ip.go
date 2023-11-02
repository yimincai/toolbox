/*
 * Copyright (c) 2023. yimincai(Neil) <bravc29229@gmail.com>.
 */

package ip

import (
	"fmt"
	"io"
	"net/http"

	"github.com/fatih/color"
)

func GetPublicIP() {
	resp, err := http.Get("http://ipecho.net/plain")
	if err != nil {
		color.Red("error: %v", err.Error())
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			color.Red("error: %v", err.Error())
		}

		fmt.Println("Your public ip is:", string(body))
	}
}
