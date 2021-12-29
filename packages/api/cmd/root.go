package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func executeRootCommand() {
	log.Println("Empty.")
}

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Base command do nothing for now",
	Run: func(cmd *cobra.Command, args []string) {
		executeRootCommand()
	},
}

func init() {
	// init things for root command
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}