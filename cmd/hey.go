package cmd

import (
	"fmt"
	"os"
	"os/exec"
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
var cfgFile string

func init() {
	cobra.OnInitialize(Initialize)
	heyCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/hey.yaml)")
}

var heyCmd = &cobra.Command{
	Use:   "hey",
	Short: "AI powered CLI helper",
	Long:  "Hey is an AI powered CLI helper: you can try for example `hey list me all files in my home directory, including hidden ones, and sort them`.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		genSpinner := spinner.New(spinner.CharSets[21], 100*time.Millisecond)
		genSpinner.Suffix = " let me think ..."

		client, err := openai.InitClient(cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		genSpinner.Start()

		genValid, genCmd, err := client.Send(strings.Join(args, " "))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		genSpinner.Stop()

		if !genValid {
			color.Red(genCmd)
			color.Red("Command execution cancelled.")
			os.Exit(0)
		}

		fmt.Print("I am about to execute: ")
		color.Yellow("`" + genCmd + "`")

		prompt := promptui.Prompt{
			Label:     "Confirm",
			IsConfirm: true,
		}

		result, err := prompt.Run()

		if err != nil {
			color.Red("Command execution cancelled.")
			os.Exit(0)
		}

		if result == "y" {

			cmd := exec.Command("bash", "-c", genCmd)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				color.Red("Command execution failure: ", err)
				os.Exit(1)
			}

			color.Green("Command execution success.")
			os.Exit(0)
		}
	},
}

func Initialize() {
	cfg = config.InitConfig(cfgFile)
}

func Execute() {
	if err := heyCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
