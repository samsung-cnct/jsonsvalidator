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

var schemas []string
var schemaDir string

// withCmd represents the with command
var withCmd = &cobra.Command{
	Use:     "with",
	Example: "with [-d /path/to/schemas] [-s schema1.json -s schemaN.json -s ...]",
	Short:   "prepratory argument to specify schema(s) against which to validate.",
	Long: `Use the <with> argument to specify the schema file(s) to validate your
config against.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return CheckRequiredFlags(cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("with called")
	},
}

func init() {
	validateCmd.AddCommand(withCmd)
	validateCmd.PersistentFlags().StringArrayVarP(&schemas, "schemas", "s", nil, "schema(s) to use for validation.")
	validateCmd.PersistentFlags().StringVarP(&schemaDir, "schema-dir", "d", "", "directory where schemas to use for validation reside.")
	fmt.Println("got here 2")
}
