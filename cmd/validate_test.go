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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
)

var testYAML = "../tests.yaml"

type testCase struct {
		name       string
		schema     string
		config     string
		jsonString string
		success    bool
		have       string
		expect     string
}

type YAMLConf struct {
		ConfigsDir string  `yaml:"configs_dir"`
		SchemasDir string  `yaml:"schemas_dir"`
		Tests      []Tests `yaml:"tests"`
}

type Tests struct {
		Name   string `yaml:"name"`
		Schema string `yaml:"schema"`
		Config string `yaml:"config"`
		Expect string `yaml:"expect"`
}

func TestTablesUsingYAML(t *testing.T) {
	var config YAMLConf

	testYamlFile, _ := ioutil.ReadFile(testYAML)
	err := yaml.Unmarshal(testYamlFile, &config)

	if err != nil {
		t.Fatalf(err.Error())
	}

	SuccessMap := map[string]bool{"success": true, "fail": false}
	SuccessMapRev := map[bool]string{true: "success", false: "fail"}

	for _, thisTest := range config.Tests {
		var testCase testCase
		var validated ValidatorResult

		// And the name is
		testCase.name = thisTest.Name

		// set default expectation.
		testCase.expect = thisTest.Expect

		cwd, err := os.Getwd()
		if err != nil {
			t.Error("could not get working directory: ", err)
		}


		// Schemas require absolute path
		schema := filepath.Join(cwd, config.SchemasDir, thisTest.Schema)

		// Configs do not
		config := filepath.Join(cwd, config.ConfigsDir, thisTest.Config)

		// Verify schema and config file for this test run exist
		testCase.schema, testCase.config = schema, config

		// register custom formatters
		registerCustomFormatters()

		// Run validation between schema and config
		jsondata, err := validate(schema, config);

		if err != nil {
			t.Error("validating failed: ", err)
		}

		commonOutStr := "\n\tTest |    %-35s| %-30s\n\tConfig: `%s`\n\tSchema: `%s`.\n\tExpected: %-20v\n\tHad: %v\n"
		commonOutErr := "\tError(s): `%+v`\n\n"

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
				testCase.schema, testCase.expect, testCase.have, validated.Exceptions)
		}

	}
}
