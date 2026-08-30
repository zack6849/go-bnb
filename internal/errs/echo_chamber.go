package errs

import (
	"errors"
	"fmt"
)

// EchoChamber accumulates errors so a caller can run a batch of fallible
// operations and report everything that went wrong at once, rather than
// bailing on the first failure.
//
// The zero value is ready to use.
type EchoChamber struct {
	errs []error
}

func CreateEchoChamber() *EchoChamber {
	return &EchoChamber{}
}

// PushError records err, err can be nil
func (ec *EchoChamber) PushError(err error) {
	if err != nil {
		ec.errs = append(ec.errs, err)
	}
}

// PushErrorf records a formatted error
func (ec *EchoChamber) PushErrorf(format string, args ...any) {
	ec.errs = append(ec.errs, fmt.Errorf(format, args...))
}

func (ec *EchoChamber) Empty() bool { return len(ec.errs) == 0 }

func (ec *EchoChamber) HasErrors() bool { return !ec.Empty() }

// Errors are the individual errors collected
func (ec *EchoChamber) Errors() []error {
	return ec.errs
}

// Summary returns every collected error as one error, or nil if none were collected.
func (ec *EchoChamber) Summary() error {
	return errors.Join(ec.errs...)
}

// Collect runs fn, recording any error and returning the zero value on failure.
// This is the general form: anything shaped (T, error) fits.
// It's a function because Go methods cannot be generic

// Collect runs fn recording any error and returning the zero value on fail
// general purpose for anything shaped (T, error)
func (ec *EchoChamber) Collect[T any](
	fn func() (T, error),
) T {
	v, err := fn()
	if err != nil {
		ec.PushError(err)
		var zero T
		return zero
	}
	return v
}

// ConvertValue looks up the field in the record and converts it, recording an error on failure.
// The field name is reused for both the lookup and the error message.
// A missing key is reported distinctly from an unparseable value
func (ec *EchoChamber) ConvertValue[T any](
	record map[string]string,
	field string,
	convert func(string) (T, error),
) T {
	var zero T
	raw, present := record[field]
	if !present {
		ec.PushErrorf("field %q is not present in the record", field)
		return zero
	}
	v, err := convert(raw)
	if err != nil {
		ec.PushErrorf("field %q: cannot convert %q: %w", field, raw, err)
		return zero
	}
	return v
}

// ConvertFieldOrDefaultValue behaves like ConvertValue, but a missing or empty value yields fallback instead of an error.
// Only a value that is present and malformed is reported.
func (ec *EchoChamber) ConvertFieldOrDefaultValue[T any](
	record map[string]string,
	field string,
	fallback T,
	convert func(string) (T, error),
) T {
	raw, present := record[field]
	if !present || raw == "" {
		return fallback
	}
	v, err := convert(raw)
	if err != nil {
		//goland:noinspection GoPrintFunctions (the %w is passed to errorf internally upstream.)
		ec.PushErrorf("field %q: cannot convert %q: %w", field, raw, err)
		return fallback
	}
	return v
}

// ConvertFieldOrZeroValue is ConvertFieldOrDefaultValue with the zero value as the fallback.
func (ec *EchoChamber) ConvertFieldOrZeroValue[T any](
	record map[string]string,
	field string,
	convert func(string) (T, error),
) T {
	var zero T
	return ec.ConvertFieldOrDefaultValue(record, field, zero, convert)
}
