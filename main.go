/*
 * Copyright (c) 2023. yimincai(Neil) <bravc29229@gmail.com>.
 */

package main

import (
	"github.com/yimincai/toolbox/cmd"
	"github.com/yimincai/toolbox/pkg/logger"
)

func main() {
	logger.Banner("NEIL TOOL BOX")
	cmd.Execute()
}
