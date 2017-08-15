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

package cmd

import (
	"errors"
	"regexp"

	jsv "jsonsvalidator/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	configFile string
	schemaFile string
)

// CheckRequiredFlags enforce Cobra to fail with an error
// when required CLI args are missing.
// TODO this will need to be revisited since flags
// may or may not require arguments. In the case of this
// application, both args (currently) require arguments.
// Unfortunately, this is an oversight of Cobra as 20170801
func CheckRequiredFlags(flags *pflag.FlagSet) error {
	requiredError := false
	flagName := ""

	flags.VisitAll(func(flag *pflag.Flag) {
		requiredAnnotation := flag.Annotations[cobra.BashCompOneRequiredFlag]

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

// RequiredFlagHasArgs Check if required flags have args
func RequiredFlagHasArgs(flag string, arg string) error {
	re := regexp.MustCompile(arg)
	match := re.FindStringSubmatch(flag)

	if len(arg) == 0 || len(match) > 0 {
		return errors.New("flag `" + flag + "` requires argument")
	}

	return nil
}

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Set config file to be validated.",
	Long: `Use the "validate" verb preparatory to specifying the flags used to` +
		`specify a JSON schema and config, -f and -c.`,
	Example: "validate  --schema <schema> --config <instance/config file>",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err = CheckRequiredFlags(cmd.Flags()); err != nil {
			return err
		}

		if err = RequiredFlagHasArgs("schema", schemaFile); err != nil {
			return err
		}

		if err = RequiredFlagHasArgs("config", configFile); err != nil {
			return err
		}

		return err
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = jsv.DoValidate(schemaFile, configFile)
		return err
	},
}

func init() {
	RootCmd.AddCommand(validateCmd)

	validateCmd.PersistentFlags().StringVarP(&schemaFile, "schema", "s", "", "schema file to validate against.")
	validateCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file to be validated.")
}
