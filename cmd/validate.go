// Copyright © 2017 Samsung CNCT - Jim Conner <snafu.x@gmail.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"os"

	"github.com/spf13/cobra"
)

/*
// FileExists check if a file exists on the system
func FileExists(name string) bool {
    if _, err := os.Stat(name); err != nil {
    if os.IsNotExist(err) {
                return false
            }
    }
    return true
}

func validArgs(cmd *cobra.Command, args []string) error {
	if err := cli.RequiresMinArgs(1)(cmd, args); err != nil {
		return err
	}

    return FileExists(args[0])
    },
*/

// place holder for the configFile
var configFile string

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Set config file to be validated.",
	Long: `Use "validate" to specify the file that needs to be ` +
		`validated via a JSON schema. Remember, this is ` +
		`THE CONFIG THAT NEEDS TO BE VALIDATED.`,
	Example: "validate  <instance/config file>",
	// the following causes Run to never get hit if
	// it succeeds. So the following erroneous invocation
	// fails: ./jsonvalidator validate with
	// however, ./jsonvalidation validate successfully fails.
	//Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TEST")
		fmt.Printf("args are %v\n", args)

		if len(args) == 0 {
			os.Stderr.WriteString("'validate' requires one argument\n")
			os.Exit(-1)
		}

		if args[1] == "with" {
			fmt.Println("'validate' requires at least on arg")
		}
	},
}

func init() {
	RootCmd.AddCommand(validateCmd)

	fmt.Println("got here 1.")
}
