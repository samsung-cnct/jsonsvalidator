// Copyright © 2017 Samsung CNCT
//
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

package jsv

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version variable contains the version of this application
	Version string

	// Build variable contains the build version of this application
	Build string
)

// versionCmd represents the version argument which will give
// the caller of the command the runtime version, build, OS, and Arch
// for the command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current version.",
	Long:  `Display the version of this application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version: ", Version)
		fmt.Println("Git commit hash: ", Build)
		fmt.Println("OS: ", runtime.GOOS)
		fmt.Println("Arch: ", runtime.GOARCH)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
