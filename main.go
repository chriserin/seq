package main

import (
	"fmt"
	"log"

	"github.com/chriserin/seq/internal/seqmidi"
	"github.com/chriserin/seq/internal/timing"
	"github.com/spf13/cobra"
	midi "gitlab.com/gomidi/midi/v2"
)

const VERSION = "0.1.0-alpha"

var cliOptions struct {
	gridTemplate string
	instrument   string
	transmitter  bool
	receiver     bool
	outport      bool
	theme        string
	midiout      string
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "seq",
		Short: "A sequencer for your cli",
		Long:  "A sequencer for your cli",
		Args:  cobra.MaximumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"seq"}, cobra.ShellCompDirectiveFilterFileExt
		},
		Run: func(cmd *cobra.Command, args []string) {

			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Recovered from panic: %v\n", r)
				}
			}()

			var err error
			midiConnection, err := seqmidi.InitMidiConnection(cliOptions.outport, cliOptions.midiout)
			if err != nil {
				fmt.Println("Midi Failure", err)
				return
			}
			defer midiConnection.Close()
			midiLoopMode := timing.MlmStandAlone
			if cliOptions.transmitter {
				midiLoopMode = timing.MlmTransmitter
			} else if cliOptions.receiver {
				midiLoopMode = timing.MlmReceiver
			}
			var filename string
			if len(args) > 0 {
				filename = args[0]
			}
			p := RunProgram(filename, midiConnection, cliOptions.gridTemplate, cliOptions.instrument, midiLoopMode, cliOptions.theme)
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
	rootCmd.Flags().StringVar(&cliOptions.gridTemplate, "template", "Drums", "Choose a template (default: Drums)")
	rootCmd.Flags().StringVar(&cliOptions.instrument, "instrument", "Standard", "Choose an instrument for CC integration (default: Standard)")
	rootCmd.Flags().BoolVar(&cliOptions.transmitter, "transmitter", false, "Seq will run in transmitter mode")
	rootCmd.Flags().BoolVar(&cliOptions.receiver, "receiver", false, "Seq will run in receiver mode")
	rootCmd.Flags().BoolVar(&cliOptions.outport, "outport", false, "Seq will create an outport to send midi")
	rootCmd.Flags().StringVar(&cliOptions.theme, "theme", "miles", "Choose an theme for the sequencer visual representation")
	rootCmd.Flags().StringVar(&cliOptions.midiout, "midiout", "", "Choose a midi out port")

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("Program failed")
	}
}
