/*
 * Copyright (c) 2023. yimincai(Neil) <bravc29229@gmail.com>.
 */

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yimincai/toolbox/internal/ip"
)

// ipCmd represents the ip command
var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Get my public ip.",
	Long:  `Get my public ip.`,
	Run: func(cmd *cobra.Command, args []string) {
		ip.GetPublicIP()
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)
}
