/*
 * Copyright (c) 2023. yimincai(Neil) <bravc29229@gmail.com>.
 */

package logger

import (
	"fmt"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

// Banner prints a banner
func Banner(str string) {
	myFigure := figure.NewColorFigure(str, "", "green", true)
	myFigure.Print()
	fmt.Println()
}

// Console prints a string with a new line
func Console(str string) {
	color.Green("=====\t%v\t=====\n", str)
}

// Green prints a string with a new line with time
func Green(str string) {
	now := time.Now()
	color.Green("%v %v\n", now.Format("2006-01-02 15:04:05"), str)
}

// Yellow prints a string with a new line with time
func Yellow(str string) {
	now := time.Now()
	color.Yellow("%v %v\n", now.Format("2006-01-02 15:04:05"), str)
}

// Blue prints a string with a new line with time
func Blue(str string) {
	now := time.Now()
	color.Blue("%v %v\n", now.Format("2006-01-02 15:04:05"), str)
}

// Red prints a string with a new line with time
func Red(str string) {
	now := time.Now()
	color.Red("%v %v\n", now.Format("2006-01-02 15:04:05"), str)
}
