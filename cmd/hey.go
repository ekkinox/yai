package cmd

import (
	"errors"
	"fmt"
	"github.com/ekkinox/hey/run"
	"os"
	"strings"
	"time"

	"github.com/ekkinox/hey/config"
	"github.com/ekkinox/hey/openai"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var cfg config.Config

func init() {
	cobra.OnInitialize(Initialize)
}

var heyCmd = &cobra.Command{
	Use:   "hey",
	Short: "AI powered CLI helper.",
	Long:  "Hey is an AI powered CLI helper.",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			for {
				prompt := promptui.Prompt{
					Label: "How can I help you? (q to quit)",
					Validate: func(input string) error {
						if strings.Trim(input, " ") == "" {
							return errors.New("Please provide an input.")
						}

						return nil
					},
				}

				input, err := prompt.Run()

				if input == "quit" || input == "q" || err != nil {
					color.HiRed("[quit]")
					os.Exit(0)
				}

				err = Process(input)
				if err != nil {
					os.Exit(1)
				}
			}
		} else {
			err := Process(strings.Join(args, " "))
			if err != nil {
				os.Exit(1)
			}
			os.Exit(0)
		}

		genSpinner := spinner.New(spinner.CharSets[21], 100*time.Millisecond)
		genSpinner.Suffix = " let me think ..."
		genSpinner.Start()

		client := openai.InitClient(cfg)

		output, err := client.Send(strings.Join(args, " "))
		if err != nil {
			color.HiRed("Error.", err)
			os.Exit(1)
		}

		genSpinner.Stop()

		if !output.Executable {
			fmt.Println(output.Content)
			os.Exit(0)
		}

		fmt.Print("I am about to execute: ")
		color.HiYellow("`" + output.Content + "`")

		prompt := promptui.Prompt{
			Label:     "Confirm",
			IsConfirm: true,
		}

		result, err := prompt.Run()

		if err != nil {
			color.HiRed("[cancelled]")
			os.Exit(0)
		}

		if result == "y" {

			err = run.RunInteractive(output.Content)
			if err != nil {
				color.HiRed("[failure]")
				os.Exit(1)
			}

			color.HiGreen("[success]")
			os.Exit(0)
		}
	},
}

func Initialize() {
	cfg = config.InitConfig()
}

func Execute() {
	if err := heyCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Process(input string) error {
	genSpinner := spinner.New(spinner.CharSets[21], 100*time.Millisecond)
	genSpinner.Suffix = " let me think ..."
	genSpinner.Start()

	client := openai.InitClient(cfg)

	output, err := client.Send(input)
	if err != nil {
		color.HiRed("Error.", err)
		return err
	}

	genSpinner.Stop()

	if !output.Executable {
		fmt.Println(output.Content)
		fmt.Println(" ")
		return nil
	}

	fmt.Print("I am about to execute: ")
	color.HiYellow("`" + output.Content + "`")

	prompt := promptui.Prompt{
		Label:     "Confirm",
		IsConfirm: true,
	}

	result, err := prompt.Run()

	if err != nil {
		color.HiRed("[cancelled]")
		fmt.Println(" ")
		return nil
	}

	if result == "y" {

		err = run.RunInteractive(output.Content)
		if err != nil {
			color.HiRed("[failure]")
			fmt.Println(" ")
			return err
		}

		color.HiGreen("[success]")
		fmt.Println(" ")
		return nil
	}

	return nil
}
