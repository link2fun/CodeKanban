package utils

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
)

// BuildLaunchURL resolves the URL that should be opened in the browser when the
// service starts. It prefers Domain and falls back to ServeAt if needed.
func BuildLaunchURL(cfg *AppConfig) string {
	if cfg == nil {
		return ""
	}

	domain := strings.TrimSpace(cfg.Domain)
	if domain == "" {
		port := strings.TrimSpace(cfg.ServeAt)
		port = strings.TrimPrefix(port, ":")
		if port != "" {
			domain = fmt.Sprintf("127.0.0.1:%s", port)
		}
	}

	if domain == "" {
		return ""
	}

	if !strings.Contains(domain, "://") {
		domain = "http://" + domain
	}

	u, err := url.Parse(domain)
	if err != nil {
		return ""
	}

	path := strings.TrimSpace(cfg.WebUrl)
	if path == "" {
		path = "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	u.Path = path
	return u.String()
}

// OpenBrowser attempts to launch the system browser for the provided URL.
func OpenBrowser(target string) error {
	if target == "" {
		return nil
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", target)
	case "darwin":
		cmd = exec.Command("open", target)
	default:
		cmd = exec.Command("xdg-open", target)
	}

	return cmd.Start()
}
