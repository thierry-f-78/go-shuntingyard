package shuntingyard

import "fmt"
import "strings"

const (
	Kind_value = iota
	Kind_group_open
	Kind_group_close
	Kind_operator

	Associativity_left
	Associativity_right

	Type_nil
	Type_bool
	Type_float64
)

func Type_desc(types []int)(string) {
	var t int
	var names []string

	if len(types) == 0 {
		return "void"
	}

	for _, t = range types {
		switch t {
		case Type_nil:     names = append(names, "nil")
		case Type_bool:    names = append(names, "bool")
		case Type_float64: names = append(names, "float64")
		default:
			panic("unknown type")
		}
	}

	return strings.Join(names, "|")
}

func Type_list(t [][]int)(string) {
	var out []string
	var e []int

	for _, e = range t {
		out = append(out, Type_desc(e))
	}

	return strings.Join(out, ", ")
}

func Type_string(t int)(string) {
	return Type_desc([]int{t})
}

type Elt interface {
	// the precedence of the operator. high level is high priority on other symbol.
	// typically or=1, and=2, not=3
	Precedence()(int)
	// the associativity of the operator. expect symbol Associativity_*. typically
	// and/or have left associativity, not has right associativity
	Associativity()(int)
	// accepted input type. array of Type_*
	Input_types()([][]int)
	// returned output types. array of Type_*
	Output_types()([][]int)
	// execute this function
	Execute([]*Value)([]*Value, error)
	// nature of element. Use Kind_*
	Kind()(int)
	// display element
	String()(string)
}

/* used as cache of Elt, prevent execution of function which return constants */
type elt_cache struct {
	precedence int
	associativity int
	input_types [][]int
	output_types [][]int
	kind int
	elt Elt
}

type Value struct {
	// type of value. use Type_*
	Kind int
	// Float64
	Value_float64 float64
	// Bool
	Value_bool bool
}

func Value_nil()(*Value) {
	return &Value{
		Kind: Type_nil,
	}
}

func Value_bool(v bool)(*Value) {
	return &Value{
		Kind: Type_bool,
		Value_bool: v,
	}
}

func Value_float64(v float64)(*Value) {
	return &Value{
		Kind: Type_float64,
		Value_float64: v,
	}
}

type Expr struct {
	rpn []*elt_cache
	// precedence stack (only used during parsing of stack)
	precedence_stack []*elt_cache
	// indicate stack ready
	done bool
	// indicates kind of consumed value
	input_types [][]int
	// indicate kind of returned value
	output_types [][]int
	// expression representation
	name_elements []string
	name string
}

func (e *Expr) String()(string) {
	return e.name
}

func New(input_values [][]int)(*Expr) {
	return &Expr{
		input_types: input_values,
	}
}

func (e *Expr)Set_name(n string) {
	e.name = n
}

func (e *Expr)Dump()() {
	var ec *elt_cache

	for _, ec = range e.rpn {
		println(ec.elt.String())
	}
}

// Append element to the expression using shuntingyard algorithm
func (e *Expr)Append(elt Elt)(error) {
	var ec *elt_cache
	var ec_browse *elt_cache

	if e.done {
		return fmt.Errorf("Expression already finalized")
	}

	/* convert to elt cache */
	ec = &elt_cache{
		precedence: elt.Precedence(),
		associativity: elt.Associativity(),
		input_types: elt.Input_types(),
		output_types: elt.Output_types(),
		kind: elt.Kind(),
		elt: elt,
	}

	/* build name */
	e.name_elements = append(e.name_elements, elt.String())

	/* pass value */
	if ec.kind == Kind_value {
		e.rpn = append(e.rpn, ec)
		return nil
	}

	/* we have open group */
	if ec.kind == Kind_group_open {
		e.precedence_stack = append(e.precedence_stack, ec)
		return nil
	}

	/* we have close parenthesis. pop precedence stack to stack until open group */
	if ec.kind == Kind_group_close {
		for {
			/* check precedence stack is not empty */
			if len(e.precedence_stack) == 0 {
				return fmt.Errorf("Expression error, encounter %q, but this symbol is not associated", ec.elt.String())
			}

			/* pop element for precedence stack */
			ec_browse = e.precedence_stack[len(e.precedence_stack) - 1]
			e.precedence_stack = e.precedence_stack[:len(e.precedence_stack) - 1]

			/* if open group, stop poping */
			if ec_browse.kind == Kind_group_open {
				break
			}

			/* push element at the top of stack */
			e.rpn = append(e.rpn, ec_browse)
		}
		return nil
	}

	/* we have operator */
	if ec.kind == Kind_operator {

		/* process operator migration from precedence stack to stack */
		for {

			/* stop if the precedence stack is empty */
			if len(e.precedence_stack) == 0 {
				break
			}

			/* get top of precedence stack element */
			ec_browse = e.precedence_stack[len(e.precedence_stack) - 1]

			/* stop if the operator at the top of the operator stack is an group open */
			if ec_browse.kind == Kind_group_open {
				break
			}

			/* pop if there is an operator at the top of the precedence stack with greater precedence */
			if ec_browse.precedence > ec.precedence {
				e.rpn = append(e.rpn, ec_browse)
				e.precedence_stack = e.precedence_stack[:len(e.precedence_stack) - 1]
				continue
			}

			/* pop if the operator at the top of the operator stack has equal precedence and is left associative */
			if ec_browse.precedence == ec.precedence && ec_browse.associativity == Associativity_left {
				e.rpn = append(e.rpn, ec_browse)
				e.precedence_stack = e.precedence_stack[:len(e.precedence_stack) - 1]
				continue
			}

			/* no condition satisfy continuation */
			break
		}

		/* push operator in the stack */
		e.precedence_stack = append(e.precedence_stack, ec)
		return nil
	}

	return fmt.Errorf("Unexpected kind value %d", ec.kind)
}

// Just push element. This is not compatible with the Append function.
// the Push function is used to build your own expression stack using
// RPN order. Mixing Push and Append return undefined result.
func (e *Expr)Push(elt Elt)(error) {
	var ec *elt_cache

	if e.done {
		return fmt.Errorf("Expression already finalized")
	}

	/* convert to elt cache */
	ec = &elt_cache{
		precedence: elt.Precedence(),
		associativity: elt.Associativity(),
		input_types: elt.Input_types(),
		output_types: elt.Output_types(),
		kind: elt.Kind(),
		elt: elt,
	}

	/* build name */
	e.name_elements = append(e.name_elements, elt.String())

	/* push value */
	e.rpn = append(e.rpn, ec)

	return nil
}

/* all provided type must be found in required type */
func has_compat(provide []int, require []int)(bool) {
	var provided_type int
	var required_type int

	if len(provide) == 0 || len(require) == 0 {
		return false
	}

	for _, provided_type = range provide {
		for _, required_type = range require {
			if provided_type == required_type {
				break
			}
		}
		if provided_type != required_type {
			return false
		}
	}
	return true
}

func (e *Expr)Finalize()(error) {
	var ec_browse *elt_cache
	var stack_types [][]int
	var i int
	var stack_index int
	var value_type []int

	if e.done {
		return fmt.Errorf("Expression already finalized")
	}	

	/* flush the stack */
	for {

		/* sprecedence stack flushed */
		if len(e.precedence_stack) == 0 {
			break
		}

		/* get top of precedence stack element */
		ec_browse = e.precedence_stack[len(e.precedence_stack) - 1]

		/* error if we encounter open group */
		if ec_browse.kind == Kind_group_open {
			return fmt.Errorf("Expression error, encounter %q, but this symbol is not associated", ec_browse.elt.String())
		}

		/* push element in the stack and pop it from precedence stack */
		e.rpn = append(e.rpn, ec_browse)
		e.precedence_stack = e.precedence_stack[:len(e.precedence_stack) - 1]
	}

	/* push inputs in the type_stack */
	for _, value_type = range e.input_types {
		stack_types = append(stack_types, value_type)
	}

	/* check the returned result */
	for _, ec_browse = range e.rpn {

		/* check number of inputs */
		if len(stack_types) < len(ec_browse.input_types) {
			return fmt.Errorf("Inconsistent expression, need %d entries, only %d available at symbol %q",
			                  len(ec_browse.input_types), len(stack_types), ec_browse.elt.String())
		}

		/* check types of inputs */
		stack_index = len(stack_types) - len(ec_browse.input_types)
		for i = 0; i < len(ec_browse.input_types); i++ {
			if !has_compat(stack_types[stack_index + i], ec_browse.input_types[i]) {
				return fmt.Errorf("Inconsistent expression, %q needs %s, got %s",
				                  ec_browse.elt.String(), Type_desc(ec_browse.input_types[i]),
				                  Type_desc(stack_types[stack_index + i]))
			}
		}

		/* pop entries from stack */
		stack_types = stack_types[:stack_index]

		/* push output in stack */
		stack_types = append(stack_types, ec_browse.output_types...)
	}

	/* store kind of returned value */
	for _, value_type = range stack_types {
		e.output_types = append(e.output_types, value_type)
	}

	/* set name */
	if e.name == "" {
		e.name =  strings.Join(e.name_elements, " ")
	}
	e.name_elements = nil

	e.done = true
	return nil
}

func (e *Expr)Exec()(*Value, error) {
	var stack []*Value
	var ec *elt_cache
	var val []*Value
	var err error

	for _, ec = range e.rpn {
		if len(stack) < len(ec.input_types) {
			return nil, fmt.Errorf("%q needs %d elements, only %d available",
			                       ec.elt.String(), len(ec.input_types), len(stack))
		}
		val, err = ec.elt.Execute(stack[len(stack) - len(ec.input_types):])
		if err != nil {
			return nil, err
		}
		stack = stack[:len(stack) - len(ec.input_types)]
		stack = append(stack, val...)
	}

	if len(stack) != 1 {
		return nil, fmt.Errorf("Expression return too many values")
	}

	return stack[0], nil
}


