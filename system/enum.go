package system

type OperatingSystem int

const (
	UnknownOperatingSystem OperatingSystem = iota
	LinuxOperatingSystem
	MacOperatingSystem
	WindowsOperatingSystem
)

func (o OperatingSystem) String() string {
	switch o {
	case LinuxOperatingSystem:
		return "linux"
	case MacOperatingSystem:
		return "macOS"
	case WindowsOperatingSystem:
		return "windows"
	default:
		return "unknown"
	}
}
