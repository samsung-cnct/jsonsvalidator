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
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"

	"io/ioutil"

	p "github.com/kr/pretty"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/xeipuuv/gojsonschema"
)

var (
	configFile string
	schemaFile string
	jsonData   []byte
)

// Result is the structure containing the validation results from JSON schema validation
type Result struct {
	IsValid   bool     `json:"is_valid"`
	Exception []string `json:"exception"`
	Config    string   `json:"config"`
	Schema    string   `json:"schema"`
}

// CIDRFormatChecker struct to extend gojsonschema FormatCheckers
type CIDRFormatChecker struct{}

// SymverFormatChecker struct to extend gojsonschema FormatCheckers
type SymverFormatChecker struct{}

// NewResult initializes Result with default values.
func NewResult() Result {
	result := Result{}
	result.IsValid = false
	result.Config = ""
	result.Exception = nil
	result.Schema = ""

	return result
}

// CheckRequiredFlags enforce Cobra to fail with an error
// when required CLI args are missing.
// TODO this will need to be revisited since flags
// may or may not require arguments. In the case of this
// application, both args (currently) require arguments.
// Unfortunately, this is an oversight of Cobra as 20170801
func CheckRequiredFlags(flags *pflag.FlagSet) error {
	requiredError := false
	flagName := ""
	//p.Println(flags)
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

// nrmlzFileContents injests a file, reads the contents, & returns
// the content in JSON. It can take YAML or JSON data. If it's JSON,
// it's just returned. If it's YAML, it's validated, JSONized, and
// then returned. If the YAML is not valid then the application
// exits with an error.
func nrmlzFileContents(configFile string) (jsonContents []byte, err error) {
	/*
	   read contents:
	   is
	     valid json? return it
	   else
	     is validyaml? yaml2json it
	   else
	     return nil, err
	*/
	var fileContents []byte

	fileContents, err = ioutil.ReadFile(configFile)

	if err != nil {
		return nil, err
	}

	if isJSON(fileContents) {
		return fileContents, nil
	}

	if jsonContents, err = yaml.YAMLToJSON(fileContents); err != nil {
		if err != nil {
			return nil, err
		}
	}

	return jsonContents, nil
}

func isJSON(b []byte) bool {
	var js interface{}
	return json.Unmarshal(b, &js) == nil
}

// FileExists check if a file exists on the system
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// IsFormat for CIDRFormatChecker - custom format checker for CIDRs
// extending gojsonschema.FormatChecker
// https://github.com/xeipuuv/gojsonschema
func (f CIDRFormatChecker) IsFormat(input string) (validCIDR bool) {
	validCIDR = false

	_, _, err := net.ParseCIDR(input)

	if err == nil {
		validCIDR = true
	}

	return validCIDR
}

// IsFormat for SymverFormatChecker - custom format checker for symver
// extending gojsonschema.FormatChecker
// https://github.com/xeipuuv/gojsonschema
func (f SymverFormatChecker) IsFormat(input string) (validSymver bool) {
	expressions := make([]*regexp.Regexp, 0)
	fmt.Println("GOT HERE")
	expressions = append(expressions, regexp.MustCompile(`^v?(0|[1-9]\d*)\.`),
		regexp.MustCompile(`(0|[1-9]\d*)\.`),
		regexp.MustCompile(`(0|[1-9]\d*)`),
		regexp.MustCompile(`(-(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?`),
		regexp.MustCompile(`(\+[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*)?$`))

	for _, re := range expressions {
		match := re.FindStringSubmatch(input)
		fmt.Println("match was...")
		p.Println(re)

		if match != nil {
			return true
		}
	}

	return false
}

func validate(schemaFile string, config []byte) {
	result := NewResult()
	result.Config = (configFile)
	result.Schema = (schemaFile)

	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaFile)
	documentLoader := gojsonschema.NewBytesLoader(config)

	//p.Println(schemaLoader)
	validated, err := gojsonschema.Validate(schemaLoader, documentLoader)

	if err != nil {
		result.Exception = append(result.Exception, err.Error())
	} else {
		if validated.Valid() {
			result.IsValid = true
		} else {
			for _, desc := range validated.Errors() {
				/*
					p.Printf("description: %v\n", desc.Description())
					p.Printf("context: %v\n", desc.Type())
					p.Printf("type: %v\n", desc.Context())
					p.Printf("field: %v\n", desc.Field())
					p.Printf("details: %v\n", desc.Details())
				*/
				e := errors.New(desc.String())
				result.Exception = append(result.Exception, e.Error())
			}
		}
	}

	jsonOutput, _ := json.Marshal(result)
	fmt.Println(string(jsonOutput))
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
		// extend the checker to handle CIDRs
		gojsonschema.FormatCheckers.Add("cidr", CIDRFormatChecker{})

		// extend the checker to handle symver
		gojsonschema.FormatCheckers.Add("symver", SymverFormatChecker{})

		if FileExists(configFile) {
			jsonData, err = nrmlzFileContents(configFile)

			if err != nil {
				return err
			}
		}

		validate(schemaFile, jsonData)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(validateCmd)

	validateCmd.PersistentFlags().StringVarP(&schemaFile, "schema", "s", "", "schema file to validate against.")
	validateCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file to be validated.")
}
