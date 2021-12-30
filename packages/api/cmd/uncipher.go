package cmd

import (
	"fmt"
	_ "github.com/marmorag/supateam/docs"
	"github.com/marmorag/supateam/internal/seeder"
	"github.com/spf13/cobra"
)

func executeUncipherCommand(filepath string) {
	deciphered, err := seeder.ReadSecureData(filepath)
	if err != nil {
		panic(err)
	}

	fmt.Println(deciphered)
}

var uncipherCmd = &cobra.Command{
	Use:   "uncipher",
	Short: "uncipher test",
	Run: func(cmd *cobra.Command, args []string) {
		file := cmd.Flag("file").Value.String()
		executeUncipherCommand(file)
	},
}

func init() {
	uncipherCmd.Flags().String("file", "", "File path to cipher/decipher")
	rootCmd.AddCommand(uncipherCmd)
}
