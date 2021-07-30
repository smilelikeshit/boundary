package dbtemplate

import "errors"

var (
	// StartDbFromTemplate creates a new test database from a postgres template database.
	StartDbFromTemplate func(dialect string, opt ...Option) (func() error, string, string, error) = startDbFromTemplateUnsupported

	// ErrDbTemplateUnsupported is returned when the platform does not support generating a database from template.
	ErrDbTemplateUnsupported = errors.New("dbtemplates is not currently supported on this platform")
)

func startDbFromTemplateUnsupported(dialect string, opt ...Option) (cleanup func() error, retURL, container string, err error) {
	return nil, "", "", ErrDbTemplateUnsupported
}
