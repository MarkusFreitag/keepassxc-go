//go:build linux
// +build linux

package keepassxc

import (
	"errors"
	"fmt"
	"os"
	"path"
)

const SocketName = "org.keepassxc.KeePassXC.BrowserServer"

var lookupPaths = []string{
	os.Getenv("XDG_RUNTIME_DIR"),
	os.Getenv("TMPDIR"),
	path.Join(os.Getenv("HOME"), "/snap/keepassxc/common/"),
	fmt.Sprintf("/run/user/%d/", os.Getuid()),
}

func SocketPath() (string, error) {
	var err error
	var filename string

	for _, base := range lookupPaths {
		filename = path.Join(base, SocketName)
		_, err = os.Stat(filename)
		switch {
		case errors.Is(err, nil):
			break
		case errors.Is(err, os.ErrNotExist):
		default:
			return "", fmt.Errorf("keepassxc socket lookup error: %s", err)
		}
	}

	if err != nil {
		return "", fmt.Errorf("keepassxc socket not found")
	}

	return filename, nil
}
