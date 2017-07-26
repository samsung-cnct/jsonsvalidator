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

	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var configFile string
var validateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "Set config file to be validated.",
	Example: "validate -c <instance/config file>",
	Long: `"validate" is the prepratory argyument which requires the ` +
		`'-c <config>' flag to specify the file that needs to be ` +
		`validated via JSON schema. Remember, this is ` +
		`THE CONFIG THAT NEEDS TO BE VALIDATED.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("inside pre-run")
		return CheckRequiredFlags(cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			fmt.Println("-c <config> is required. No other arguments are recognized.")
		}

		if len(args) == 0 {
			fmt.Println("-c <config> is required")
			cmd.Usage()
		}

		fmt.Printf("got here. args are %v with length %d: \n", args, len(args))
	},
}

func init() {
	RootCmd.AddCommand(validateCmd)
	validateCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file to be validated.")

	fmt.Println("got here 1.")
}
