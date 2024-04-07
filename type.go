package shuntingyard

import "strings"

// This interface describe a type. The type is a type of
// value and it is used to check grammatical compatibility
type Type interface {
	Name()(string)
}

// Describe alternative of types
func Type_desc(types []Type)(string) {
	var t Type
	var names []string

	if len(types) == 0 {
		return "void"
	}

	for _, t = range types {
		names = append(names, t.Name())
	}

	return strings.Join(names, "|")
}

// Describe list of altrenatives of types.
func Type_list(t [][]Type)(string) {
	var out []string
	var e []Type

	for _, e = range t {
		out = append(out, Type_desc(e))
	}

	return strings.Join(out, ", ")
}

// Describe simple type
func Type_string(t Type)(string) {
	return Type_desc([]Type{t})
}
