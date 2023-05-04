package tools

import (
	"os/user"

	"github.com/fatih/color"
)

func GetUsername() string {
	user, err := user.Current()
	if err != nil {
		color.Red("error:%v", err.Error())
	}

	return user.Username
}
