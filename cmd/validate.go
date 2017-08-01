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
	"os"

	"io/ioutil"

	//p "github.com/kr/pretty"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
)

var (
	configFile string
	schemaFile string
	jsonData   []byte
)

// Result is the structure containing the validation results from JSON schema validation
type Result struct {
	IsValid   bool    `json:"is_valid"`
	Exception []error `json:"exception"`
	Config    string  `json:"config"`
	Schema    string  `json:"schema"`
}

// NewResult initializes Result with default values.
func NewResult() Result {
	result := Result{}
	result.IsValid = false
	result.Config = ""
	result.Exception = nil
	result.Schema = ""

	return result
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

//func validate(schemaFile string, config []byte) error {
func validate(schemaFile string, config []byte) {
	result := NewResult()
	result.Config = (configFile)
	result.Schema = (schemaFile)

	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaFile)
	documentLoader := gojsonschema.NewBytesLoader(config)

	validated, err := gojsonschema.Validate(schemaLoader, documentLoader)

	if err != nil {
		//fmt.Printf("err is %s", err)
		result.Exception = append(result.Exception, err)
	} else {
		if validated.Valid() {
			result.IsValid = true
		} else {
			for _, desc := range validated.Errors() {
				e := errors.New(desc.String())
				result.Exception = append(result.Exception, e)
			}
		}
	}

	um, _ := json.Marshal(result)
	fmt.Println(string(um))
	//return err
}

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Set config file to be validated.",
	Long: `Use the "validate" verb preparatory to specifying the flags used to` +
		`specify a JSON schema and config, -f and -c.`,
	Example: "validate  -s <schema> -c <instance/config file>",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if FileExists(configFile) {
			jsonData, err = nrmlzFileContents(configFile)
		}

		/*
			if err = validate(schemaFile, jsonData); err != nil {
				return err
			}
		*/
		validate(schemaFile, jsonData)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(validateCmd)

	validateCmd.PersistentFlags().StringVarP(&schemaFile, "schema", "s", "", "schema file to validate against.")
	validateCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file to be validated.")
}
