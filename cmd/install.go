/*
Copyright © 2023 yimincai <bravc29229@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/yimincai/toolbox/installer"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install applications.",
	Long:  `Using install scripts to install applications.`,
	Run: func(cmd *cobra.Command, args []string) {

		promptSystem := promptui.Select{
			Label: "Select your system",
			Items: []string{"Ubuntu", "MacOS", "Windows", "Cancel"},
		}

		_, system, err := promptSystem.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		if system == "Cancel" {
			os.Exit(0)
		}

		switch system {
		case "Ubuntu":
			promptApplication := promptui.Select{
				Label: "Select the application you want to install",
				Items: []string{"Docker", "Cancel"},
			}

			_, application, err := promptApplication.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			switch application {
			case "Docker":
				installer.UbuntuInstallDocker()
			}

		default:
			fmt.Println("Not support yet.")

		}

	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
