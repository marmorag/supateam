package cmd

import (
	"fmt"
	_ "github.com/marmorag/supateam/docs"
	"github.com/marmorag/supateam/internal/seeder"
	"github.com/spf13/cobra"
)

func executeCipherCommand(filepath string) {
	ciphered, err := seeder.WriteSecureData(filepath)
	if err != nil {
		panic(err)
	}

	fmt.Println(ciphered)
}

var cipherCmd = &cobra.Command{
	Use:   "cipher",
	Short: "cipher file",
	Run: func(cmd *cobra.Command, args []string) {
		file := cmd.Flag("file").Value.String()
		executeCipherCommand(file)
	},
}

func init() {
	cipherCmd.Flags().String("file", "", "File path to cipher")
	rootCmd.AddCommand(cipherCmd)
}
