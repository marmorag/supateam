package cmd

import (
	"fmt"
	_ "github.com/marmorag/supateam/docs"
	"github.com/marmorag/supateam/internal"
	"github.com/marmorag/supateam/internal/repository"
	"github.com/marmorag/supateam/internal/seeder"
	"github.com/spf13/cobra"
)

var seederToApply string

func executeSeedCommand() {
	_ = internal.GetConfig()
	defer repository.CloseConnection()

	if s := seeder.Mapping[seeder.Name(seederToApply)]; s != nil {
		err := s.Seed()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Apply seed to the database",
	Run: func(cmd *cobra.Command, args []string) {
		executeSeedCommand()
	},
}

func init() {
	seedCmd.Flags().StringVarP(&seederToApply, "seeder", "s", "", "Which seeder to apply")
	seedCmd.MarkFlagRequired("seeder")

	rootCmd.AddCommand(seedCmd)
}
