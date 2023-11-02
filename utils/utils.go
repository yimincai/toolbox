/*
 * Copyright (c) 2023. yimincai(Neil) <bravc29229@gmail.com>.
 */

package utils

import (
	"github.com/fatih/color"
	"os/user"
)

func GetUsername() string {
	u, err := user.Current()
	if err != nil {
		color.Red("error:%v", err.Error())
	}

	return u.Username
}
