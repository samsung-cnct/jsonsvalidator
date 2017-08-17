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
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"

	"io/ioutil"

	"github.com/blang/semver"
	"github.com/ghodss/yaml"
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

// SemVerFormatChecker struct to extend gojsonschema FormatCheckers
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
	goodSemVer := false
	_, ok := semver.Make(input)

	if ok == nil {
		goodSemVer = true
	}

	return goodSemVer
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
	}

	result := NewResult()
	result.Config = (configFile)
	result.Schema = (schemaFile)
	result.IsValid = false
	result.Exception = append(result.Exception, err.Error())

	errResult, err := json.Marshal(result)

	return string(errResult), err
}

// Validate is simply a method alias which points to JSONDataRespValidate making
// `Validate` by default return a go data structure to the caller. If one wants
// to get a JSON string returned explicitly then one must call `JSONStrRespValidate()`
func Validate(schemaFile string, configFile string) (jsonOutput []byte, err error) {
	return JSONDataRespValidate(schemaFile, configFile)
}

// RegisterCustomFormatters initializes any custom formatters for jsongoschema to
// validate our special JSON types.
func RegisterCustomFormatters() {
	// extend the checker to handle CIDRs
	gojsonschema.FormatCheckers.Add("cidr", CIDRFormatChecker{})

	// extend the checker to handle symver
	gojsonschema.FormatCheckers.Add("semver", SemVerFormatChecker{})
}

// DoValidate is the entry point into validating JSON documents
func DoValidate(schemaFile, configFile string) (err error) {
	RegisterCustomFormatters()

	if jsonstr, err := JSONStrRespValidate(schemaFile, configFile); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(jsonstr)
	}

	return err
}
