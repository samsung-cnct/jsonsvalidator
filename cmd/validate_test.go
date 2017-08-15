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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	y "gopkg.in/yaml.v2"
)

var testYAML string = "./tests.yaml"
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

	my_yaml, _ := ioutil.ReadFile(testYAML)
	err := y.Unmarshal(my_yaml, &config)

	if err != nil {
		t.Fatalf(err.Error())
	}

	schemaTestsDir = config.SchemasDir
	configTestsDir = config.ConfigsDir
	tests_to_run := config.Tests

	SuccessMap := map[string]bool{"success": true, "fail": false}
	SuccessMapRev := map[bool]string{true: "success", false: "fail"}

	for _, this_test := range tests_to_run {
		var test_case testCase
		var validated Result

		// And the name is
		test_case.name = this_test.Name

		// set default expectation.
		test_case.expect = this_test.Expect

		cwd, err := os.Getwd()

		// Schemas require absolute path
		schema := filepath.Join(cwd, schemaTestsDir, this_test.Schema)

		// Configs do not
		config := filepath.Join(configTestsDir, this_test.Config)

		// Verify schema and config file for this test run exist
		test_case.schema, test_case.config = confFiles(schema, config)

		// register custom formatters
		RegisterCustomFormatters()

		// Run validation between schema and config
		if jsondata, ok := Validate(schema, config); ok == nil {
			common_out := "\n\tTest |    %-35s| %-30s\n\tConfig: `%s`\n\tSchema: `%s`.\n\tExpected: %-20v\n\tHad: %v\n"
			with_error := "\tError(s): `%s`\n\n"

			if err = json.Unmarshal(jsondata, &validated); err != nil {
				t.Fatalf(err.Error())
			}

			test_case.have = SuccessMapRev[validated.IsValid]

			if validated.IsValid == SuccessMap[this_test.Expect] {
				test_case.success = true
				t.Logf(common_out+"\n", test_case.name, "SUCCEEDED", test_case.config,
					test_case.schema, test_case.expect, test_case.have)
			} else {
				test_case.success = false
				t.Errorf(common_out+with_error, test_case.name, "FAILED!!", test_case.config,
					test_case.schema, test_case.expect, test_case.have, validated.Exception)
			}
		}
	}
}
