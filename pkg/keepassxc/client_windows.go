//go:build windows
// +build windows

package keepassxc

import (
	"fmt"
	"net"
	"os"

	"github.com/Microsoft/go-winio"
)

// SocketName is a standard KeepassXC socket name.
const SocketName = "org.keepassxc.KeePassXC.BrowserServer"

func SocketPath() (string, error) {
	return fmt.Sprintf(`\\.\pipe\%s_%s`, SocketName, os.Getenv("USERNAME")), nil
}

func Connect(socketPath string) (net.Conn, error) {
	return winio.DialPipe(socketPath, nil)
}
