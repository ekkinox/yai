package detect

import (
	"os"
	"runtime"
	"strings"

	"github.com/ekkinox/hey/run"

	"github.com/mitchellh/go-homedir"
)

const OS_linux = "linux"
const OS_darwin = "darwin"
const OS_windows = "windows"
const OS_other = "other"

func DetectOperatingSystem() string {
	switch runtime.GOOS {
	case "linux":
		return OS_linux
	case "darwin":
		return OS_darwin
	case "windows":
		return OS_windows
	default:
		return OS_other
	}
}

func DetectDistribution() string {
	dist, err := run.Run("lsb_release", "-sd")
	if err != nil {
		return ""
	}

	return strings.Trim(strings.Trim(dist, "\n"), "\"")
}

func DetectShell() string {
	shell, err := run.Run("echo", os.Getenv("SHELL"))
	if err != nil {
		return ""
	}

	split := strings.Split(strings.Trim(strings.Trim(shell, "\n"), "\""), "/")

	return split[len(split)-1]

}

func DetectHomeDir() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		return ""
	}

	return homeDir
}
