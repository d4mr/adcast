package cmd

import (
	"fmt"
	"os"

	cobra "github.com/spf13/cobra"
)

var AsciiArt = `
                $$\                                $$\     
                $$ |                               $$ |    
 $$$$$$\   $$$$$$$ | $$$$$$$\ $$$$$$\   $$$$$$$\ $$$$$$\   
 \____$$\ $$  __$$ |$$  _____|\____$$\ $$  _____|\_$$  _|  
 $$$$$$$ |$$ /  $$ |$$ /      $$$$$$$ |\$$$$$$\    $$ |    
$$  __$$ |$$ |  $$ |$$ |     $$  __$$ | \____$$\   $$ |$$\ 
\$$$$$$$ |\$$$$$$$ |\$$$$$$$\\$$$$$$$ |$$$$$$$  |  \$$$$  |
 \_______| \_______| \_______|\_______|\_______/    \____/ 	
`

var rootCmd = &cobra.Command{
	Use:   "adcast",
	Short: "Adcast is a FFMPEG wrapper for ad embedding",
	Long: `Adcast allows for Ad insertion into video files, without re-encoding through HLS or through JIT re-encoding.
It supports following optimations for fast stitching via FFMPEG:
	- media pre-segmentation at user specified breakpoints
	- maintaining ad media cache in the same format as the content

Made with love, by d4mr`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(AsciiArt)
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
