/*
 * Copyright (c) 2023. yimincai(Neil) <bravc29229@gmail.com>.
 */

package utils

import (
	"os/user"

	"github.com/fatih/color"
)

// GetUsername gets the current system username
func GetUsername() string {
	u, err := user.Current()
	if err != nil {
		color.Red("error:%v", err.Error())
	}

	return u.Username
}
