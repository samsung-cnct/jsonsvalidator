/*
Copyright 2017 The Samsung CNCT
All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.


Command validates your JSON config file using the JSON schema specification
draft 4.0.

	$ go get samsung-cnct/golang-tools/jsonsvalidator

To build:

	$ ./build.sh
	$ make all install

Which should create a binary called jsonsvalidator

To use:

	$ jsonsvalidator validate --config /path/to/instance.json --schema /path/to/your/schema.json

Get help:

	$ jsonsvalidator help

FAQ:

	Q. Does jsonsvalidator support YAML instance (config) files?
	A. Yes. Simply specifying the YAML config instead of a JSON config will succeed.
*/

package jsv
