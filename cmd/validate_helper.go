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
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

// fileContentsNormalizer injests a file, reads the contents, & returns
// the content in JSON. It can take YAML or JSON data. If it's JSON,
// it's just returned. If it's YAML, it's validated, JSONized, and
// then returned. If the YAML is not valid then the application
// exits with an error.
func fileContentsNormalizer(configFile string) ([]byte, error) {
	fileContents, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	if isJSON(fileContents) {
		return fileContents, nil
	}

	jsonContents, err := yaml.YAMLToJSON(fileContents)
	if err != nil {
		return nil, err
	}

	return jsonContents, nil
}

func isJSON(b []byte) bool {
	var js interface{}
	return json.Unmarshal(b, &js) == nil
}

// fileExists check if a file exists on the system
func fileExists(name string) (bool, error) {
	file, err := os.Stat(name)
	if err != nil {
		return !os.IsNotExist(err), err
	}

	if file.IsDir() {
		return false,  errors.New("invalid file specified; requested file is a directory, not a file")
	}

	return true, nil
}

// JSONDataRespValidate will is the main function that performs validation. However,
// this function always returns a data structure `result` of type struct containing
// data for the validation.
func JSONDataRespValidate(schemaFile string, configFile string) (jsonOutput []byte, err error) {
	if _, err := fileExists(configFile); err != nil {
		return nil, err
	}

	if _, err := fileExists(schemaFile); err != nil {
		return nil, err
	}

	jsonData, err := fileContentsNormalizer(configFile)
	if err != nil {
		return nil, err
	}

	documentLoader := gojsonschema.NewBytesLoader(jsonData)
	// XXX allow reference loader to use URLs as well as local files
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaFile)
	validated, err := gojsonschema.Validate(schemaLoader, documentLoader)

	result := ValidatorResult{
		IsValid: false,
		Exceptions: []ExceptionDetail{},
		Config: configFile,
		Schema: schemaFile,
	}

	if err != nil {
		result.appendExceptionWithPath(err, "a general exception occurred; probably an invalid schema; see https://github.com/xeipuuv/gojsonschema/issues/160")
	} else {
		if validated.Valid() {
			result.IsValid = true
		} else {
			for _, desc := range validated.Errors() {
				result.appendExceptionMessage(desc.String(), desc.Context().String())
			}
		}
	}

	return json.Marshal(result)
}

// jsonStrRespValidate calls JSONDataRespValidate() and marshalls the JSON response to
// a string returning the string to the caller.
func jsonStrRespValidate(schemaFile string, configFile string) (jsonOutput string, err error) {
	jsonResponse, err := validate(schemaFile, configFile)
	if err == nil {
		return string(jsonResponse), nil
	}

	result := ValidatorResult{
		IsValid: false,
		Exceptions: []ExceptionDetail{},
		Config: configFile,
		Schema: schemaFile,
	}

	result.appendException(err)

	errResult, err := json.Marshal(result)

	if err != nil {
		result.appendException(err)
	}

	return string(errResult), err
}

// validate is simply a method alias which points to JSONDataRespValidate making
// `validate` by default return a go data structure to the caller. If one wants
// to get a JSON string returned explicitly then one must call `jsonStrRespValidate()`
func validate(schemaFile string, configFile string) (jsonOutput []byte, err error) {
	return JSONDataRespValidate(schemaFile, configFile)
}

// DoValidate is the entry point into validating JSON documents
func doValidate(schemaFile string, configFile string) (err error) {
	registerCustomFormatters()

	if jsonstr, err := jsonStrRespValidate(schemaFile, configFile); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(jsonstr)
	}

	return err
}
