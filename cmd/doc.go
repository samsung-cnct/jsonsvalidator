/*

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

package main // import "samsung-cnct/golang-tools/jsonsvalidator"