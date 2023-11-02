/*
 * Copyright (c) 2023. yimincai(Neil) <bravc29229@gmail.com>.
 */

package logger

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"time"
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

func Green(str string) {
	now := time.Now()
	color.Green("%v %v\n", now.Format("2006-01-02 15:04:05"), str)
}

func Yellow(str string) {
	now := time.Now()
	color.Yellow("%v %v\n", now.Format("2006-01-02 15:04:05"), str)
}

func Blue(str string) {
	now := time.Now()
	color.Blue("%v %v\n", now.Format("2006-01-02 15:04:05"), str)
}

func Red(str string) {
	now := time.Now()
	color.Red("%v %v\n", now.Format("2006-01-02 15:04:05"), str)
}
