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
	"errors"
	"fmt"
	"os"

	//p "github.com/kr/pretty"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var cfgFile string
var runtimeCommandName = os.Args[0]
var example1 = fmt.Sprintf("%s validate -c <config> with -s <schema>...", runtimeCommandName)
var example2 = fmt.Sprintf("%s validate -c <config> with -d <schema_dir>", runtimeCommandName)

// CheckRequiredFlags used to assist in determining if required flags not
// specified by CLI.
func CheckRequiredFlags(flags *pflag.FlagSet) error {
	requiredError := false
	flagName := ""

	/*
		p.Print(flags)
		fmt.Println("\n--------------------------------")
	*/

	flags.VisitAll(func(flag *pflag.Flag) {

		/*
			fmt.Println("cobra.BashComOneRequiredFlag:")
			p.Print(cobra.BashCompOneRequiredFlag)
			fmt.Println("\n--------------------------------")
		*/

		requiredAnnotation := flag.Annotations[cobra.BashCompOneRequiredFlag]

		//fmt.Printf("CHECKED REQUIRED FLAGS! Annotation is: %v\n", flag.Annotations)

		if len(requiredAnnotation) == 0 {
			return
		}

		flagRequired := requiredAnnotation[0] == "true"

		if flagRequired && !flag.Changed {
			requiredError = true
			flagName = flag.Name
		}
	})

	if requiredError {
		return errors.New("Required flag `" + flagName + "` has not been set")
	}

	return nil
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     runtimeCommandName,
	Short:   "Validate JSON config against the JSON schema validator (spec v4).",
	Example: example1 + "\n" + example2,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
