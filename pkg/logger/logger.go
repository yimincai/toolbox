/*
 * Copyright (c) 2023. yimincai(Neil) <bravc29229@gmail.com>.
 */

package logger

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

func Banner(str string) {
	myFigure := figure.NewColorFigure(str, "", "green", true)
	myFigure.Print()
	fmt.Println()
}

func Console(str string) {
	color.Green("=====\t%v\t=====\n", str)
}

func Notice(str string) {
	color.Yellow("%v\n", str)
}
