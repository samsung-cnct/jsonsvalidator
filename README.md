# __JSON Schema Validator__ 

## Description

*jsonsvalidator* performs JSON schema validation using a specified schema and an arbitrary config (YAML or JSON) file.

## Synopsis
`./jsonsvalidator validate --schema /path/to/schema.json --config /path/to/config.yaml`

## Build process
Currently still manual IE:

`$ go build -o jsonsvalidator main.go`


## Dependencies

[gojsonschema library (spec 4)](https://github.com/xeipuuv/gojsonschema)

[YAML validation library](https://github.com/ghodss/yaml)

[Cobra](https://github.com/spf13/cobra)

[Cobra POSIX Flags library](https://github.com/spf13/pflag)

## More to come
