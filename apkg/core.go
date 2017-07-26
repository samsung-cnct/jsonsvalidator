/*
Copyright 2017 The <project> Authors. All rights reserved.

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
	"fmt"
)

// Verbose is a debugging flag
var Verbose bool

// WebPort is a hack.  If present then make a webpage available at that port
var WebPort string

// Status is a componnet status enum.
type Status int

// Values for the status enum.
const (
	UNKNOWN Status = 1 + iota
	STARTING
	ALIVE
	TERMINATING
	DEAD
)

var statuses = [...]string{
	"UNKNOWN",
	"STARTING",
	"ALIVE",
	"TERMINATING",
	"DEAD",
}

func (s Status) String() string {
	return statuses[s-1]
}

//
// Component represents something in the system.
type Component struct {
	ID     string
	Type   string
	Status Status
}

// NewComponent creates with an UNKNOWN status.
func NewComponent(id string, t string) *Component {
	return &Component{id, t, UNKNOWN}
}

func (c Component) String() string {
	return fmt.Sprintf("ID:%s Type:%s Status:%s", string(c.ID), string(c.Type), c.Status.String())
}

// Criteria represents some generic selection criteria.
type Criteria struct {
}

// Result is a list of selected components.
type Result struct {
	Comps []Component
}

// A Targeter Target returns a list of selected components.
//
// USAGE: (c Criteria) Target(a SomeConnection) (res Result, err error)
type Targeter interface {
	Target(a API) (res Result, err error)
}

// An API represents a targetable interface.
//
// Example Usage:  func (a API) Target(c Criteria) (res Result){ actual code }
type API interface {
}
