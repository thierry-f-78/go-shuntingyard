package shuntingyard

import "context"
import "fmt"

const (
	Kind_value = iota
	Kind_group_open
	Kind_group_close
	Kind_operator

	Associativity_left
	Associativity_right
)

func kind_str(k int)(string) {
	switch k {
	case Kind_value: return "value"
	case Kind_group_open: return "group-open"
	case Kind_group_close: return "group-close"
	case Kind_operator: return "operator"
	}
	return fmt.Sprintf("unknown #%d", k)
}

type Elt interface {
	// the precedence of the operator. high level is high priority on other symbol.
	// typically or=1, and=2, not=3
	Precedence()(int)
	// the associativity of the operator. expect symbol Associativity_*. typically
	// and/or have left associativity, not has right associativity
	Associativity()(int)
	// accepted input type. array of Type_*
	Input_types()([][]Type)
	// returned output types. array of Type_*
	Output_types()([][]Type)
	// execute this function
	Execute(context.Context, []Value)([]Value, error)
	// nature of element. Use Kind_*
	Kind()(int)
	// display element
	String()(string)
}
