package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

const VERSION = "0.1.0-alpha"

func main() {
	rootCmd := &cobra.Command{
		Use:   "Seq",
		Short: "A sequencer for your cli",
		Long:  "A sequencer for your cli",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			println("seq")
		},
	}

	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "Version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("seq v%s\n", VERSION)
		},
	}

	rootCmd.AddCommand(cmdVersion)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("Program failed")
	}
}
