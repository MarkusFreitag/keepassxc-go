/*
Copyright © 2021 Markus Freitag <fmarkus@mailbox.org>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	plaintext bool
	showAll   bool
)

var getLoginsCmd = &cobra.Command{
	Use:   "get-logins URL",
	Short: "query info for the specified url",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := client.GetLogins(args[0])
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			return fmt.Errorf("could not find entries for '%s'", args[0])
		}
		for _, entry := range entries {
			pass := entry.Password.String()
			if plaintext {
				pass = entry.Password.Plaintext()
			}
			fmt.Printf("%s %s %s", entry.Name, entry.Login, pass)
			if entry.Expired {
				fmt.Print(" EXPIRED")
			}
			fmt.Print("\n")

			if !showAll {
				break
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getLoginsCmd)

	getLoginsCmd.Flags().BoolVar(&plaintext, "plaintext", false, "print out the password - BE CAREFUL")
	getLoginsCmd.Flags().BoolVar(&showAll, "all", false, "show all matches otherwise only the first will be printed")
}
