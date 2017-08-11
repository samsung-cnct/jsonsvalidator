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
	re "regexp"
	"testing"
	//v "jsonsvalidator/cmd"
	p "github.com/kr/pretty"
)

var testSchemaDir string = "./test_schemas"
var testConfigDir string = "./test_configs"

func pickFiles(files []os.FileInfo, regex string) (matching_files []os.FileInfo, err error) {
	match := re.MustCompile(regex)

	for _, file := range files {
		if match.MatchString(file.Name()) {
			matching_files = append(matching_files, file)
		}
	}

	return matching_files, err
}

type cidrTest struct {
	schema     string
	config     string
	jsonstring string
	success    bool
	have       bool
	expect     bool
}

var tests cidrTest

func buildTable() (err error) {
	schema := fmt.Sprintf("%s/config.json", testSchemaDir)
	configs, err := ioutil.ReadDir(testConfigDir)

	if err == nil {
		configs, err = pickFiles(configs, "cidr.*")

		if len(configs) == 0 {
			fmt.Println("No config files for CIDR checks. Stop.")
			return
		}

		for _, cf := range configs {
			jsondata, err1 := Validate(schema, cf.Name())
			jsonstr, err2 := json.Marshal(string(jsondata))

			if err1 == nil {
				if err2 != nil {
					err = err2
					return err
				}
			}

			if tests.have == tests.expect {
				tests.success = true
			}

			tests.schema = schema
			tests.config = cf.Name()
			tests.jsonstring = string(jsonstr)
			tests.expect = true
		}
	}

	p.Println(tests)
	return nil
}

func TestValidCIDR(t *testing.T) {
}

func main() {
	table := buildTable()
	//TestValidCIDR(validCIDRtable)
}
