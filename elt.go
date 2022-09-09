package shuntingyard

import "fmt"

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

type Elt interface {
	// the precedence of the operator. high level is high priority on other symbol.
	// typically or=1, and=2, not=3
	Precedence()(int)
	// the associativity of the operator. expect symbol Associativity_*. typically
	// and/or have left associativity, not has right associativity
	Associativity()(int)
	// number of stack elements required
	Needs_elements()(int)
	// number of stack elements returned
	Return_elements()(int)
	// execute this function
	Execute([]*Value)([]*Value)
	// nature of element. Use Kind_*
	Kind()(int)
	// display element
	String()(string)
}

/* used as cache of Elt, prevent execution of function which return constants */
type elt_cache struct {
	precedence int
	associativity int
	needs_elements int
	return_elements int
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

func (v *Value)Float64()(float64) {
	switch v.Kind {
	case Type_float64: return v.Value_float64
	case Type_bool: if v.Value_bool { return 1.0 } else { return 0.0 }	
	case Type_nil: return 0.0
	default: panic("unexpected type")
	}
}

func (v *Value)Bool()(bool) {
	switch v.Kind {
	case Type_bool: return v.Value_bool
	case Type_float64: if v.Value_float64 == 0.0 { return false } else { return true }
	case Type_nil: return false
	default: panic("unexpected type")
	}
}

func (v *Value)String()(string) {
	switch v.Kind {
	case Type_float64: return fmt.Sprintf("%f", v.Value_float64)
	case Type_bool:    return fmt.Sprintf("%t", v.Value_bool)
	case Type_nil:     return "nil"
	default: panic("unexpected type")
	}
}

type Expr struct {
	rpn []*elt_cache
	// precedence stack (only used during parsing of stack)
	precedence_stack []*elt_cache
	// indicate stack ready
	done bool
}

func New()(*Expr) {
	return &Expr{}
}

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
		needs_elements: elt.Needs_elements(),
		return_elements: elt.Return_elements(),
		kind: elt.Kind(),
		elt: elt,
	}

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

func (e *Expr)Finalize()(error) {
	var ec_browse *elt_cache
	var stack_size int

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

	/* check the compute return one result */
	stack_size = 0
	for _, ec_browse = range e.rpn {
		stack_size -= ec_browse.needs_elements
		if stack_size < 0 {
			return fmt.Errorf("Inconsistent expression, need more entries than exists")
		}
		stack_size += ec_browse.return_elements
	}
	if stack_size == 0 {
		return fmt.Errorf("Expression doesn't return value")
	}
	if stack_size != 1 {
		return fmt.Errorf("Expression return too many values")
	}

	e.done = true
	return nil
}

func (e *Expr)Exec()(*Value, error) {
	var stack []*Value
	var ec *elt_cache
	var val []*Value

	for _, ec = range e.rpn {
		if len(stack) < ec.needs_elements {
			return nil, fmt.Errorf("%q needs %d elements, only %d available",
			                       ec.elt.String(), ec.needs_elements, len(stack))
		}
		val = ec.elt.Execute(stack[len(stack) - ec.needs_elements:])
		stack = stack[:len(stack) - ec.needs_elements]
		stack = append(stack, val...)
	}

	if len(stack) != 1 {
		return nil, fmt.Errorf("Expression return too many values")
	}

	return stack[0], nil
}


