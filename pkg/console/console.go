package console

import (
	"github.com/fatih/color"
)

// Green println prints a string with a new line
func Green(str string) {
	color.Green("%v\n", str)
}

// Yellow println prints a string with a new line
func Yellow(str string) {
	color.Yellow("%v\n", str)
}

// Blue println prints a string with a new line
func Blue(str string) {
	color.Blue("%v\n", str)
}

// Red println prints a string with a new line
func Red(str string) {
	color.Red("%v\n", str)
}
