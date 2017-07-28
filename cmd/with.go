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
	Example: "with <filename(s),director(y|ies)>",
	Short:   "prepratory argument to specify schema(s) against which to validate.",
	Long: `Use the <with> argument to specify the schema file(s) to validate your
config against.`,
	//Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("with called")
	},
}

func init() {
	validateCmd.AddCommand(withCmd)
	fmt.Println("got here 2")
}
