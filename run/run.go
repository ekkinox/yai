package run

import (
	"os"
	"os/exec"
)

func Run(cmd string, arg ...string) (string, error) {
	out, err := exec.Command(cmd, arg...).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func RunInteractive(cmd string) error {
	run := exec.Command("bash", "-c", cmd)
	run.Stdin = os.Stdin
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr

	if err := run.Run(); err != nil {
		return err
	}

	return nil
}
