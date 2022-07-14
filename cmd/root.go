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

var (
	profileName string
	client      *keepassxc.Client
)

var rootCmd = &cobra.Command{
	Use:   "keepassxc-go",
	Short: "interact with keepassxc via unix-socket",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		socketPath := keepassxc.SocketPath()
		if _, err := os.Stat(socketPath); os.IsNotExist(err) {
			return fmt.Errorf("keepassxc socket not found '%s'", socketPath)
		}

		store, err := keystore.Load()
		if err != nil {
			return err
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
				return errors.New("keystore has multiple profiles, please specify the one to use")
			}
			for _, profile := range store.Profiles {
				if profile.Name == profileName {
					key = profile.NaclKey()
				}
			}
			if key == nil {
				return fmt.Errorf("could not find profile '%s'", profileName)
			}
		}

		client = keepassxc.NewClient(socketPath, profileName, key)
		if err := client.Connect(); err != nil {
			return err
		}

		if err := client.ChangePublicKeys(); err != nil {
			return err
		}

		if key == nil {
			if err := client.Associate(); err != nil {
				return err
			}
			name, key := client.GetAssociatedProfile()
			err = store.Add(&keystore.Profile{Name: name, Key: key})
			if err != nil {
				return err
			}
			err = store.Save()
			if err != nil {
				return err
			}
		} else {
			if err := client.TestAssociate(); err != nil {
				return err
			}
		}

		return nil
	},
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
