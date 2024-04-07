package shuntingyard

import "context"
import "fmt"
import "os"
import "strings"

/* used as cache of Elt, prevent execution of function which return constants */
type elt_cache struct {
	precedence int
	associativity int
	input_types [][]Type
	output_types [][]Type
	kind int
	elt Elt
}

type Expr struct {
	rpn []*elt_cache
	// precedence stack (only used during parsing of stack)
	precedence_stack []*elt_cache
	// indicate stack ready
	done bool
	// indicates kind of consumed value
	input_types [][]Type
	// indicate kind of returned value
	output_types [][]Type
	// expression representation
	name_elements []string
	name string
}

/* Implement Elt interface for Expr expression, except Execute which is located below */
func (e *Expr) Precedence()(int) {
	return 0
}
func (e *Expr) Associativity()(int) {
	return 0
}
func (e *Expr) Input_types()([][]Type) {
	return e.input_types
}
func (e *Expr) Output_types()([][]Type) {
	return e.output_types
}
func (e *Expr) Kind()(int) {
	return Kind_value
}
func (e *Expr) String()(string) {
	return e.name
}

func New(input_values [][]Type)(*Expr) {
	return &Expr{
		input_types: input_values,
	}
}

func (e *Expr)Set_name(n string) {
	e.name = n
}

func (e *Expr)dump(level int)() {
	var ec *elt_cache
	var ex *Expr
	var ok bool

	fmt.Fprintf(os.Stderr, "%s[%s]:\n", strings.Repeat("|   ", level - 1), e.String())
	for _, ec = range e.rpn {
		ex, ok = ec.elt.(*Expr)
		if ok {
			ex.dump(level + 1)
		} else {
			fmt.Fprintf(os.Stderr, "%s%s\n", strings.Repeat("|   ", level), ec.elt.String())
		}
	}
}

func (e *Expr)Dump()() {
	e.dump(1)
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

	return fmt.Errorf("Unexpected kind value %s for %s", kind_str(ec.kind), ec.elt.String())
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
func Has_compat(provide []Type, require []Type)(bool) {
	var provided_type Type
	var required_type Type

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
	var stack_types [][]Type
	var i int
	var stack_index int
	var value_type []Type

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
			if !Has_compat(stack_types[stack_index + i], ec_browse.input_types[i]) {
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

/* part of implementation of Elt interface for Expr expression */
func (e *Expr)Execute(ctx context.Context, in []Value)([]Value, error) {
	var stack []Value
	var ec *elt_cache
	var val []Value
	var err error

	/* push input value in the stack */
	stack = append(stack, in...)

	for _, ec = range e.rpn {
		if len(stack) < len(ec.input_types) {
			return nil, fmt.Errorf("%q needs %d elements, only %d available",
			                       ec.elt.String(), len(ec.input_types), len(stack))
		}
		val, err = ec.elt.Execute(ctx, stack[len(stack) - len(ec.input_types):])
		if err != nil {
			return nil, err
		}
		stack = stack[:len(stack) - len(ec.input_types)]
		stack = append(stack, val...)
	}

	return stack, nil
}
