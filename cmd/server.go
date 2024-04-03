package cmd

import (
	"fmt"

	"github.com/d4mr/adcast/server"
	cobra "github.com/spf13/cobra"
)

var mediaDir string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(AsciiArt)
		server.StartServer(mediaDir)
	},
}

func init() {
	serverCmd.Flags().StringVarP(&mediaDir, "media-directory", "m", "", "Directory to serve media from")
	serverCmd.MarkFlagRequired("media-directory")

	rootCmd.AddCommand(serverCmd)
}
