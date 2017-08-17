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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	y "gopkg.in/yaml.v2"
)

var testYAML = "./tests.yaml"
var schemaTestsDir string
var configTestsDir string

type (
	testCase struct {
		name       string
		schema     string
		config     string
		jsonstring string
		success    bool
		have       string
		expect     string
	}

	YAMLConf struct {
		ConfigsDir string  `yaml:"configs_dir"`
		SchemasDir string  `yaml:"schemas_dir"`
		Tests      []Tests `yaml:"tests"`
	}

	Tests struct {
		Name   string `yaml:"name"`
		Schema string `yaml:"schema"`
		Config string `yaml:"config"`
		Expect string `yaml:"expect"`
	}
)

func confFiles(schema, config string) (string, string) {
	if _, err := FileExists(schema); err != nil {
		fmt.Println(err.Error())
		return "", ""
	}

	if _, err := FileExists(config); err != nil {
		fmt.Println(err.Error())
		return "", ""
	}

	return schema, config
}

func TestTablesUsingYAML(t *testing.T) {
	var config YAMLConf

	myYAML, _ := ioutil.ReadFile(testYAML)
	err := y.Unmarshal(myYAML, &config)

	if err != nil {
		t.Fatalf(err.Error())
	}

	schemaTestsDir = config.SchemasDir
	configTestsDir = config.ConfigsDir
	testsToRun := config.Tests

	SuccessMap := map[string]bool{"success": true, "fail": false}
	SuccessMapRev := map[bool]string{true: "success", false: "fail"}

	for _, thisTest := range testsToRun {
		var testCase testCase
		var validated Result

		// And the name is
		testCase.name = thisTest.Name

		// set default expectation.
		testCase.expect = thisTest.Expect

		cwd, err := os.Getwd()

		// Schemas require absolute path
		schema := filepath.Join(cwd, schemaTestsDir, thisTest.Schema)

		// Configs do not
		config := filepath.Join(configTestsDir, thisTest.Config)

		// Verify schema and config file for this test run exist
		testCase.schema, testCase.config = confFiles(schema, config)

		// register custom formatters
		RegisterCustomFormatters()

		// Run validation between schema and config
		if jsondata, ok := Validate(schema, config); ok == nil {
			commonOutStr := "\n\tTest |    %-35s| %-30s\n\tConfig: `%s`\n\tSchema: `%s`.\n\tExpected: %-20v\n\tHad: %v\n"
			commonOutErr := "\tError(s): `%s`\n\n"

			if err = json.Unmarshal(jsondata, &validated); err != nil {
				t.Fatalf(err.Error())
			}

			testCase.have = SuccessMapRev[validated.IsValid]

			if validated.IsValid == SuccessMap[thisTest.Expect] {
				testCase.success = true
				t.Logf(commonOutStr+"\n", testCase.name, "SUCCEEDED", testCase.config,
					testCase.schema, testCase.expect, testCase.have)
			} else {
				testCase.success = false
				t.Errorf(commonOutStr+commonOutErr, testCase.name, "FAILED!!", testCase.config,
					testCase.schema, testCase.expect, testCase.have, validated.Exception)
			}
		}
	}
}
