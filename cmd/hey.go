package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ekkinox/hey/config"
	"github.com/ekkinox/hey/openai"

	"github.com/briandowns/spinner"
	execute "github.com/commander-cli/cmd"
	"github.com/fatih/color"
	"github.com/guumaster/logsymbols"
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
	Long:  "Hey is an AI powered CLI helper: for example `hey list all files in my home directory`.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		info := color.New(color.Bold, color.FgBlue).PrintlnFunc()
		success := color.New(color.Bold, color.FgGreen).PrintlnFunc()
		error := color.New(color.Bold, color.FgRed).PrintlnFunc()

		genSpinner := spinner.New(spinner.CharSets[21], 100*time.Millisecond)
		genSpinner.Color("blue")
		genSpinner.Suffix = " generating ..."

		client, err := openai.InitClient(cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		genSpinner.Start()

		genCmd, err := client.Send(strings.Join(args, " "))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		genSpinner.Stop()

		fmt.Print(logsymbols.Info)
		info("Generated command: ")
		color.Red(genCmd)

		prompt := promptui.Prompt{
			Label:     "Apply",
			IsConfirm: true,
		}

		result, err := prompt.Run()

		if err != nil {
			fmt.Println("Cancelled.")
			os.Exit(0)
		}

		if result == "y" {

			execSpinner := spinner.New(spinner.CharSets[21], 100*time.Millisecond)
			execSpinner.Color("blue")
			execSpinner.Suffix = " executing ..."

			c := execute.NewCommand(genCmd)

			genSpinner.Start()

			err := c.Execute()
			if err != nil {
				panic(err.Error())
			}

			genSpinner.Stop()

			if c.ExitCode() == 0 {
				fmt.Print(logsymbols.Success)
				success("Execution success, output: ")
				fmt.Println(c.Stdout())
			} else if c.ExitCode() == 1 {
				fmt.Print(logsymbols.Error)
				error("Execution error, output: ")
				fmt.Println(c.Stderr())
			}
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
