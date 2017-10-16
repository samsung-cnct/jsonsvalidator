package cmd

import (
	"net"

	"github.com/blang/semver"
	"github.com/xeipuuv/gojsonschema"
)

// ValidatorResult is the structure containing the validation results from JSON schema validation
type ValidatorResult struct {
	IsValid    bool              `json:"is_valid"`
	Exceptions []ExceptionDetail `json:"exception"`
	Config     string            `json:"config"`
	Schema     string            `json:"schema"`
}

// ExceptionDetail contains error messages and path. It is part of the ValidatorResult struct.
type ExceptionDetail struct {
	ErrorString string `json:"error_string"`
	Path        string `json:"path"`
}

func (r *ValidatorResult) appendException(err error) {
	r.appendExceptionWithPath(err, "Not Reported")
}

func (r *ValidatorResult) appendExceptionWithPath(err error, path string) {
	r.appendExceptionMessage(err.Error(), path)
}

func (r *ValidatorResult) appendExceptionMessage(errMsg string, path string) {
	exception := ExceptionDetail {
		ErrorString: errMsg,
		Path: path,
	}

	r.Exceptions = append(r.Exceptions, exception)
}

// CIDRFormatChecker struct to extend gojsonschema FormatCheckers
type CIDRFormatChecker struct{}

// IsFormat for CIDRFormatChecker - custom format checker for CIDRs
// extending gojsonschema.FormatChecker
// https://github.com/xeipuuv/gojsonschema
func (f CIDRFormatChecker) IsFormat(input interface{}) bool {
	if input == nil {
		return false
	}

	_, _, err := net.ParseCIDR(input.(string))

	return err == nil
}


// SemVerFormatChecker struct to extend gojsonschema FormatCheckers
type SemVerFormatChecker struct{}

// IsFormat for SemVerFormatChecker - custom format checker for semantics version format
// extending gojsonschema.FormatChecker
// https://github.com/xeipuuv/gojsonschema
func (f SemVerFormatChecker) IsFormat(input interface{}) bool {
	if input == nil {
		return false
	}

	_, err := semver.Make(input.(string))

	return err == nil
}

// RegisterCustomFormatters initializes any custom formatters for jsongoschema to
// validate our special JSON types.
func registerCustomFormatters() {
	// extend the checker to handle CIDRs
	gojsonschema.FormatCheckers.Add("cidr", CIDRFormatChecker{})

	// extend the checker to handle symver
	gojsonschema.FormatCheckers.Add("semver", SemVerFormatChecker{})
}
