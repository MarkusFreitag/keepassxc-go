//go:build darwin
// +build darwin

package keepassxc

func SocketPath() (string, error) {
	tmpDir, ok := os.LookupEnv("TMPDIR")
	if !ok {
		return "", errors.New("$TMPDIR not set, can not find socket")
	}

	path := filepath.Join(tmpDir, "/org.keepassxc.KeePassXC.BrowserServer")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("keepassxc socket not found '%s'", path)
	}
	return path, nil
}
