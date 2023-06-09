/*
Copyright © 2023 yimincai <bravc29229@gmail.com>
*/
package tools

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

func PrintBanner(str string) {
	myFigure := figure.NewColorFigure(str, "", "green", true)
	myFigure.Print()
	fmt.Println()
}

func PrintConsole(str string) {
	color.Green("=====\t%v\t=====\n", str)
}

func PrintNotice(str string) {
	color.Yellow("%v\n", str)
}
