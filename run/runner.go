package run

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func RunCommand(cmd string, arg ...string) (string, error) {
	out, err := exec.Command(cmd, arg...).Output()
	if err != nil {
		message := fmt.Sprintf("error: %v", err)
		log.Println(message)

		return message, err
	}

	return string(out), nil
}

func PrepareInteractiveCommand(input string) *exec.Cmd {
	return exec.Command(
		"bash",
		"-c",
		fmt.Sprintf("echo \"\n=> output:\n\";%s; echo \"\n\";", strings.TrimRight(input, ";")),
	)
}
