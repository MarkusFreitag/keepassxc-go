/*
Copyright © 2021 Markus Freitag <fmarkus@mailbox.org>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/kevinburke/nacl"
	"github.com/spf13/cobra"

	"github.com/MarkusFreitag/keepassxc-go/pkg/keepassxc"
	"github.com/MarkusFreitag/keepassxc-go/pkg/keystore"
)

var profileName string

var rootCmd = &cobra.Command{
	Use:   "keepassxc-go",
	Short: "interact with keepassxc via unix-socket",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		msg := err.Error()
		if strings.HasPrefix(msg, "error") {
			rootCmd.SilenceUsage = true
		}

		fmt.Fprintln(os.Stderr, msg)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SilenceErrors = true

	rootCmd.PersistentFlags().StringVarP(&profileName, "profile", "p", "", "Only necessary if keystore contains multiple profiles")
}

func initializeClient() (*keepassxc.Client, error) {
	socketPath, err := keepassxc.SocketPath()
	if err != nil {
		return nil, err
	}

	store, err := keystore.Load()
	if err != nil {
		return nil, err
	}

	var key nacl.Key
	switch len(store.Profiles) {
	case 0:
		break
	case 1:
		key = store.Profiles[0].NaclKey()
		profileName = store.Profiles[0].Name
	default:
		if profileName == "" {
			return nil, errors.New("keystore has multiple profiles, please specify the one to use")
		}
		for _, profile := range store.Profiles {
			if profile.Name == profileName {
				key = profile.NaclKey()
			}
		}
		if key == nil {
			return nil, fmt.Errorf("could not find profile '%s'", profileName)
		}
	}

	client := keepassxc.NewClient(socketPath, profileName, key)
	if err := client.Connect(); err != nil {
		return nil, err
	}

	if err := client.ChangePublicKeys(); err != nil {
		return nil, err
	}

	if key == nil {
		if err := client.Associate(); err != nil {
			return nil, err
		}
		name, key := client.GetAssociatedProfile()
		err = store.Add(&keystore.Profile{Name: name, Key: key})
		if err != nil {
			return nil, err
		}
		err = store.Save()
		if err != nil {
			return nil, err
		}
	} else {
		if err := client.TestAssociate(); err != nil {
			return nil, err
		}
	}

	return client, nil
}
