/*
Copyright © 2024 gucchisk <gucchi_sk@yahoo.co.jp>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ratelimitreq",
	Short: "Request with rate limit",
	Long: `Request with rate limit`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
