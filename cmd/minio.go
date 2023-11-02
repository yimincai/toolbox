/*
 * Copyright (c) 2023. yimincai(Neil) <bravc29229@gmail.com>.
 */

package cmd

import (
	"fmt"
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/yimincai/toolbox/tools"
)

// minioCmd represents the minio command
var minioCmd = &cobra.Command{
	Use:   "minio",
	Short: "minio client tool.",
	Long:  `minio client tool, you can use it to manage your minio server including: backup, restore, etc.`,
	Run: func(cmd *cobra.Command, args []string) {

		var op string

		if len(args) == 0 {
			promptSystem := promptui.Select{
				Label: "Select your operation",
				Items: []string{"Dump", "Delete", "Restore", "Upload", "Cancel"},
			}

			var err error

			_, op, err = promptSystem.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}
		} else {
			op = firstCharToUpper(args[0])
		}

		if len(args) == 2 {
			tools.DestinationDir = args[1]
		}

		switch op {
		case "Cancel":
			os.Exit(0)
		case "Dump":
			tools.DumpMinioBucket()
		case "Delete":
			tools.DeleteMinioBucket()
		case "Restore":
			tools.MinioRestore()
		case "Upload":
			tools.UploadMinioBucket()
		default:
			fmt.Println("Input error.")
		}
	},
}

func firstCharToUpper(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size <= 1 {
		return s
	}
	lc := unicode.ToUpper(r)
	if r == lc {
		return s
	}
	return string(lc) + s[size:]
}

func init() {
	rootCmd.AddCommand(minioCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// minioCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// minioCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
