package keepassxc_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MarkusFreitag/keepassxc-go/pkg/keepassxc"
)

func TestSocketPath(t *testing.T) {
	require.Regexp(
		t,
		regexp.MustCompile(`/run/user/\d+/org.keepassxc.KeePassXC.BrowserServer`),
		keepassxc.SocketPath(),
	)
}
