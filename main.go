package main

import (
	"fmt"
	"log"

	"github.com/chriserin/seq/internal/seqmidi"
	"github.com/spf13/cobra"
	midi "gitlab.com/gomidi/midi/v2"
)

const VERSION = "0.1.0-alpha"

var gridTemplate string
var instrument string
var transmitter bool
var receiver bool
var output bool
var theme string

func main() {
	rootCmd := &cobra.Command{
		Use:   "Seq",
		Short: "A sequencer for your cli",
		Long:  "A sequencer for your cli",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Recovered from panic: %v\n", r)
				}
			}()

			var err error
			midiConnection, err := seqmidi.InitMidiConnection(output)
			if err != nil {
				fmt.Println("Midi Failure", err)
				return
			}
			defer midiConnection.Close()
			midiLoopMode := MLM_STAND_ALONE
			if transmitter {
				midiLoopMode = MLM_TRANSMITTER
			} else if receiver {
				midiLoopMode = MLM_RECEIVER
			}
			var filename string
			if len(args) > 0 {
				filename = args[0]
			}
			p := RunProgram(filename, midiConnection, gridTemplate, instrument, midiLoopMode, theme)
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
	rootCmd.Flags().StringVar(&gridTemplate, "template", "Drums", "Choose a template (default: Drums)")
	rootCmd.Flags().StringVar(&instrument, "instrument", "Standard", "Choose an instrument for CC integration (default: Standard)")
	rootCmd.Flags().BoolVar(&transmitter, "transmitter", false, "Seq will run in transmitter mode")
	rootCmd.Flags().BoolVar(&receiver, "receiver", false, "Seq will run in receiver mode")
	rootCmd.Flags().BoolVar(&output, "output", false, "Seq will create an output to send midi")
	rootCmd.Flags().StringVar(&theme, "theme", "miles", "Choose an theme for the sequencer visual representation")

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("Program failed")
	}
}
