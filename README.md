# __JSON Schema Validator__ 

## Description

*jsonsvalidator* performs JSON schema validation using a specified schema and an arbitrary config (YAML or JSON) file.

## Synopsis
`./jsonsvalidator validate --schema /path/to/schema.json --config /path/to/config.yaml`

## Gotchas
1. The `/path/to/schema` must be a fully qualified path.
2. Currently, the validator does not handle remote schemas, yet.
3. According to the author's observations, the library which does the actual
   validation ignores `$ref`s during validation if a schema file is used for
   validation. Currently, this application only recognizes validation by 
   local file.

## Build process
`$ make all`

## TODO
1. Update so that remote schemas IE URLs can be used for validation.

## Non-standard Dependencies

[gojsonschema library (spec 4)](https://github.com/xeipuuv/gojsonschema)

[YAML validation library](https://github.com/ghodss/yaml)

[Cobra](https://github.com/spf13/cobra)

[Cobra POSIX Flags library](https://github.com/spf13/pflag)

[Semantic Version Library (spec ver. 2.0.0 compliant)](https://github.com/blang/semver)

## More to come
