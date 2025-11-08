package system

import (
	"os/exec"
	"runtime"
)

// OpenExplorer opens the platform specific file manager at the provided path.
func OpenExplorer(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		for _, candidate := range []string{"xdg-open", "nautilus", "dolphin", "thunar"} {
			if _, err := exec.LookPath(candidate); err == nil {
				cmd = exec.Command(candidate, path)
				break
			}
		}
		if cmd == nil {
			return ErrNoFileManager
		}
	default:
		return ErrUnsupportedOS
	}

	return cmd.Start()
}
