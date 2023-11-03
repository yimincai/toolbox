/*
 * Copyright (c) 2023. yimincai(Neil) <bravc29229@gmail.com>.
 */

package cmd

import (
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/yimincai/toolbox/internal/minio"
	"github.com/yimincai/toolbox/pkg/console"
	"github.com/yimincai/toolbox/pkg/logger"

	"github.com/spf13/cobra"
)

// minioCmd represents the minio command
var minioCmd = &cobra.Command{
	Use:   "minio",
	Short: "minio client tool.",
	Long: `minio client tool, you can use it to manage your minio server including: backup, restore, etc.
	For example: toolbox minio dump endpoint bucket user password [destinationDir] useSSL(true/false), toolbox minio dump 127.0.0.1:9000 test test 123456 ./backup/minio true, last two parameters are optional.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 5 {
			console.Red("Input error, please check your input.")
			console.Red("For example: toolbox minio dump endpoint bucket user password [destinationDir] useSSL(true/false), last two parameters are optional.")
			console.Blue("toolbox minio dump 127.0.0.1:9000 example user password ./backup/minio false")
			return
		}

		if len(args) != 0 {
			op := firstCharToUpper(args[0])
			minio.Endpoint = args[1]
			minio.BucketName = args[2]
			minio.User = args[3]
			minio.Password = args[4]

			if len(args) == 6 {
				minio.DestinationDir = args[5]
			}

			if len(args) == 7 {
				ssl := args[6]
				if ssl == "true" {
					minio.UseSSL = true
				} else if ssl == "false" {
					minio.UseSSL = false
				} else {
					logger.Red("Invalid value for useSSL. Please provide either true or false.")
				}
			}

			switch op {
			case "Cancel":
				os.Exit(0)
			case "Dump":
				minio.DumpBucket(10)
			case "Delete":
				minio.DeleteBucket(10)
			case "Restore":
				minio.RestoreBucket(10)
			case "Upload":
				minio.UploadBucket(10)
			default:
				logger.Red("Input error")
			}
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
}
