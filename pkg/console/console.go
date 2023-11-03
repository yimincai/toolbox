package console

import (
	"github.com/fatih/color"
)

func Green(str string) {
	color.Green("%v\n", str)
}

func Yellow(str string) {
	color.Yellow("%v\n", str)
}

func Blue(str string) {
	color.Blue("%v\n", str)
}

func Red(str string) {
	color.Red("%v\n", str)
}
