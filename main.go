package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	midi "gitlab.com/gomidi/midi/v2"
)

const VERSION = "0.1.0-alpha"

func main() {
	rootCmd := &cobra.Command{
		Use:   "Seq",
		Short: "A sequencer for your cli",
		Long:  "A sequencer for your cli",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			outport, err := midi.OutPort(0)
			if err != nil {
				panic("Did not get midi outport")
			}
			defer outport.Close()
			p := RunProgram(outport)
			_, err = p.Run()
			if err != nil {
				log.Fatal("Program Failure")
			} else {
				return
			}
		},
	}

	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "Version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("seq v%s\n", VERSION)
		},
	}

	cmdListOutports := &cobra.Command{
		Use:   "list",
		Short: "List available midi outports",
		Run: func(cmd *cobra.Command, args []string) {
			outports := midi.GetOutPorts()
			for i, outport := range outports {
				fmt.Printf("%d) %s\n", i+1, outport)
			}
		},
	}

	rootCmd.AddCommand(cmdListOutports)
	rootCmd.AddCommand(cmdVersion)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("Program failed")
	}
}
