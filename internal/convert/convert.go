// Package convert provides string converters shaped func(string) (T, error),
// suitable for passing to errors.Field.
//
// strconv.Atoi already has this shape and can be passed directly; the helpers
// here exist for the strconv functions that take extra arguments.
package convert

import "strconv"

// Float returns a converter that parses a float of the given bit size.
func Float(bitSize int) func(string) (float64, error) {
	return func(s string) (float64, error) {
		return strconv.ParseFloat(s, bitSize)
	}
}

// Int returns a converter that parses a signed integer in the given base.
func Int(base, bitSize int) func(string) (int64, error) {
	return func(s string) (int64, error) {
		return strconv.ParseInt(s, base, bitSize)
	}
}
