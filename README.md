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
   local schema file.

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

## Maintainer Versioning
This may change IE it's experimental. In order to bump the version on this app one must
tag a commit to the version one wants:

1. git commit -am 'some message here'
2. git tag x.x.x-foo OR some other valid semantic in conformance with semver 2.0.0 (semver.org)
3. The build will automatically add the build date, hash, and other details into the binary so
   `app version` reflects the proper information
