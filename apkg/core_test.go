/*
Copyright 2017 The <package> Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apkg

import (
	"reflect"
	"testing"
)

func TestNewComponent(t *testing.T) {
	const cid string = "test"
	const ctype string = "one"
	const cstat Status = UNKNOWN
	bc := &Component{cid, ctype, cstat}
	nc := NewComponent(cid, ctype)
	if *nc != *bc {
		t.Error(
			"Expected ", bc,
			" got ", nc,
		)
	}
}

func TestToString(t *testing.T) {
	nc := NewComponent("test2", "another")
	ns := nc.String()
	const rs string = "ID:test2 Type:another Status:UNKNOWN"
	if rs != ns {
		t.Error(
			"Expected ", rs,
			" got ", ns,
		)
	}
}

func (c *Criteria) Target(a API) (res *Result, err error) {
	res = a.(*Result)
	return res, nil
}

func TestExampleTarget(t *testing.T) {

	c := Criteria{}
	r := &Result{}

	res, err := c.Target(r)
	if err != nil {
		t.Error("Unexpected Target err ", err)
	}
	if !reflect.DeepEqual(r, res) {
		t.Error("Expected ", r, " got ", res)
	}
}
