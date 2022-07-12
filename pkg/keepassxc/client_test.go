package keepassxc_test

import (
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/require"

	"github.com/MarkusFreitag/keepassxc-go/pkg/keepassxc"
)

func TestSocketPath(t *testing.T) {
	monkey.Patch(os.Getuid, func() int {
		return 1000
	})
	defer monkey.Unpatch(os.Getuid)

	require.Equal(
		t,
		"/run/user/1000/org.keepassxc.KeePassXC.BrowserServer",
		keepassxc.SocketPath(),
	)
}
