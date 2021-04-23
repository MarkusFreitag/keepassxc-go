package main

import (
	"fmt"

	"github.com/MarkusFreitag/golang-keepassxc/keepassxc"
)

func main() {
	socket := "/run/user/1000/org.keepassxc.KeePassXC.BrowserServer"

	client := keepassxc.NewClient(socket, "golang", nil)
	if err := client.Connect(); err != nil {
		fmt.Printf("connect err: %s\n", err.Error())
		return
	}

	if err := client.ChangePublicKeys(); err != nil {
		fmt.Printf("change-public-keys err: %s\n", err.Error())
		return
	}

	if err := client.Associate(); err != nil {
		fmt.Printf("associate err: %s\n", err.Error())
		return
	}

	if err := client.TestAssociate(); err != nil {
		fmt.Printf("test-associate err: %s\n", err.Error())
		return
	}

	entries, err := client.GetLogins("https://github.com")
	if err != nil {
		fmt.Printf("get-logins err: %s\n", err.Error())
		return
	}
	for _, entry := range entries {
		fmt.Printf("%s %s %s", entry.Name, entry.Login, entry.Password)
		if entry.Expired {
			fmt.Print(" expired")
		}
		fmt.Print("\n")
	}
}
