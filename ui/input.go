package ui

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type UiInput struct {
	runMode    RunMode
	promptMode PromptMode
	args       string
	pipe       string
}

func NewUIInput() (*UiInput, error) {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var exec, chat bool
	flagSet.BoolVar(&exec, "e", false, "exec prompt mode")
	flagSet.BoolVar(&chat, "c", false, "chat prompt mode")
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return nil, err
	}

	args := flagSet.Args()

	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("Error getting stat:", err)
		return nil, err
	}

	pipe := ""
	if !(stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0) {
		reader := bufio.NewReader(os.Stdin)
		var builder strings.Builder

		for {
			r, _, err := reader.ReadRune()
			if err != nil && err == io.EOF {
				break
			}
			_, err = builder.WriteRune(r)
			if err != nil {
				fmt.Println("Error getting input:", err)
				return nil, err
			}
		}

		pipe = strings.TrimSpace(builder.String())
	}

	runMode := ReplMode
	if len(args) > 0 {
		runMode = CliMode
	}

	promptMode := DefaultPromptMode
	if exec && !chat {
		promptMode = ExecPromptMode
	} else if !exec && chat {
		promptMode = ChatPromptMode
	}

	return &UiInput{
		runMode:    runMode,
		promptMode: promptMode,
		args:       strings.Join(args, " "),
		pipe:       pipe,
	}, nil
}

func (i *UiInput) GetRunMode() RunMode {
	return i.runMode
}

func (i *UiInput) GetPromptMode() PromptMode {
	return i.promptMode
}

func (i *UiInput) GetArgs() string {
	return i.args
}

func (i *UiInput) GetPipe() string {
	return i.pipe
}
