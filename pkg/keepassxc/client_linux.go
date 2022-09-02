//go:build linux
// +build linux

package keepassxc

import (
	"fmt"
	"os"
)

func SocketPath() (string, error) {
	path := fmt.Sprintf(
		"/run/user/%d/org.keepassxc.KeePassXC.BrowserServer", os.Getuid(),
	)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("keepassxc socket not found '%s'", path)
	}
	return path, nil
}
