package main

import (
	"fmt"
	"os"

	"github.com/godcong/excavator"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:        "ex",
	Aliases:    nil,
	SuggestFor: nil,
	Short:      "ex parsing dictionary",
	Long:       `ex parsing dictionary and insert to db`,
}

func main() {
	host := rootCmd.PersistentFlags().StringP("host", "H", "http://localhost", "set the root url path")
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if err := excavator.CommonlyTop(*host, excavator.CommonlyBase); err != nil {
			panic(err)
		}
	}
	rootCmd.SuggestionsMinimumDistance = 1
	Execute()
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
