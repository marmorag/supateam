package cmd

import (
	"fmt"
	_ "github.com/marmorag/supateam/docs"
	"github.com/marmorag/supateam/internal"
	"github.com/marmorag/supateam/internal/repository"
	"github.com/marmorag/supateam/internal/seeder"
	"github.com/spf13/cobra"
)

func executeSeedCommand() {
	_ = internal.GetConfig()
	defer repository.CloseConnection()

	s := seeder.Seeder{}
	err := s.Seed()

	if err != nil {
		fmt.Println(err.Error())
	}
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Sedd base application",
	Run: func(cmd *cobra.Command, args []string) {
		executeSeedCommand()
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}
