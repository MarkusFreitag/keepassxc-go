//go:build freebsd
// +build freebsd

package keepassxc

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"strconv"
)

const SocketName = "org.keepassxc.KeePassXC.BrowserServer"

var lookupPaths = []string{
	os.Getenv("XDG_RUNTIME_DIR"),
	os.Getenv("TMPDIR"),
	"/tmp",
	path.Join("/var/run/user/", strconv.Itoa(os.Getuid())),
}

func SocketPath() (string, error) {
	var filename string

	for _, base := range lookupPaths {
		filename = path.Join(base, SocketName)

		if _, err := os.Stat(filename); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", fmt.Errorf("keepassxc socket lookup error: %s", err)
		}

		break
	}

	if filename == "" {
		return "", fmt.Errorf("keepassxc socket not found")
	}

	return filename, nil
}

func connect(socketPath string) (net.Conn, error) {
	return net.DialUnix("unix", nil, &net.UnixAddr{Name: socketPath, Net: "unix"})
}
