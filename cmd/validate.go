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

	"github.com/blang/semver"
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

// SemverFormatChecker struct to extend gojsonschema FormatCheckers
type SemVerFormatChecker struct{}

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

// FileContentsNormalizer injests a file, reads the contents, & returns
// the content in JSON. It can take YAML or JSON data. If it's JSON,
// it's just returned. If it's YAML, it's validated, JSONized, and
// then returned. If the YAML is not valid then the application
// exits with an error.
func FileContentsNormalizer(configFile string) (jsonContents []byte, err error) {
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
func FileExists(name string) (exists bool, err error) {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false, err
		}
	}

	return true, err
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

// IsFormat for SemVerFormatChecker - custom format checker for semantics version format
// extending gojsonschema.FormatChecker
// https://github.com/xeipuuv/gojsonschema
func (f SemVerFormatChecker) IsFormat(input string) (validSemVer bool) {
	_, ok := semver.Make(input)

	if ok != nil {
		return false
	}

	return true
}

// JSONDataRespValidate will is the main function that performs validation. However,
// this function always returns a data structure `result` of type struct containing
// data for the validation.
func JSONDataRespValidate(schemaFile string, configFile string) (jsonOutput []byte, err error) {
	result := NewResult()
	result.Config = (configFile)
	result.Schema = (schemaFile)
	fexists, err := FileExists(configFile)

	if fexists {
		jsonData, err = FileContentsNormalizer(configFile)

		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaFile)
	documentLoader := gojsonschema.NewBytesLoader(jsonData)

	validated, err := gojsonschema.Validate(schemaLoader, documentLoader)

	if err != nil {
		result.Exception = append(result.Exception, err.Error())
	} else {
		if validated.Valid() {
			result.IsValid = true
		} else {
			for _, desc := range validated.Errors() {
				e := errors.New(desc.String())
				result.Exception = append(result.Exception, e.Error())
			}
		}
	}

	return json.Marshal(result)
}

// JSONStrRespValidate calls JSONDataRespValidate() and marshalls the JSON response to
// a string returning the string to the caller.
func JSONStrRespValidate(schemaFile string, configFile string) (jsonOutput string, err error) {
	if jsonResponse, err := Validate(schemaFile, configFile); err == nil {
		return string(jsonResponse), err
	} else {
		result := NewResult()
		result.Config = (configFile)
		result.Schema = (schemaFile)
		result.IsValid = false
		result.Exception = append(result.Exception, err.Error())

		errResult, err := json.Marshal(result)

		return string(errResult), err
	}
}

// Validate is simply a method alias which points to JSONDataRespValidate making
// `Validate` by default return a go data structure to the caller. If one wants
// to get a JSON string returned explicitly then one must call `JSONStrRespValidate()`
func Validate(schemaFile string, configFile string) (jsonOutput []byte, err error) {
	return JSONDataRespValidate(schemaFile, configFile)
}

func RegisterCustomFormatters() {
	// extend the checker to handle CIDRs
	gojsonschema.FormatCheckers.Add("cidr", CIDRFormatChecker{})

	// extend the checker to handle symver
	gojsonschema.FormatCheckers.Add("semver", SemVerFormatChecker{})
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
		RegisterCustomFormatters()

		if jsonstr, err := JSONStrRespValidate(schemaFile, configFile); err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(jsonstr)
		}

		return err
	},
}

func init() {
	RootCmd.AddCommand(validateCmd)

	validateCmd.PersistentFlags().StringVarP(&schemaFile, "schema", "s", "", "schema file to validate against.")
	validateCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file to be validated.")
}
