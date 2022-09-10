package shuntingyard

import "testing"

type test struct {
   precedence int
   associativity int
   needs_elements int
   return_elements int
   kind int
	symbol string
	input_types [][]int
	output_types [][]int
}

func (t *test)Precedence()(int) { return t.precedence }
func (t *test)Associativity()(int) { return t.associativity }
func (t *test)Needs_elements()(int) { return t.needs_elements }
func (t *test)Return_elements()(int) { return t.return_elements }
func (t *test)Kind()(int) { return t.kind }
func (t *test)String()(string) { return t.symbol }
func (t *test)Input_types()([][]int) { return t.input_types }
func (t *test)Output_types()([][]int) { return t.output_types }
func (t *test)Execute(vs []*Value)([]*Value) {
	switch t.symbol {
	case "true": return []*Value{Value_bool(true)}
	case "false": return []*Value{Value_bool(false)}
	case "or": return []*Value{Value_bool(vs[0].Bool() || vs[1].Bool())}
	case "and": return []*Value{Value_bool(vs[0].Bool() && vs[1].Bool())}
	case "not":  return []*Value{Value_bool(!vs[0].Bool())}
	case "neg":  return []*Value{Value_float64(-vs[0].Float64())}
	case "*":  return []*Value{Value_float64(vs[0].Float64() * vs[1].Float64())}
	case "/":  return []*Value{Value_float64(vs[0].Float64() / vs[1].Float64())}
	case "+":  return []*Value{Value_float64(vs[0].Float64() + vs[1].Float64())}
	case "-":  return []*Value{Value_float64(vs[0].Float64() - vs[1].Float64())}
	case "2.3":  return []*Value{Value_float64(2.3)}
	case "2.4":  return []*Value{Value_float64(2.4)}
	case "2.5":  return []*Value{Value_float64(2.5)}
	case "2.6":  return []*Value{Value_float64(2.6)}
	case "2.6_or_nil":  return []*Value{Value_float64(2.6)}
	case "coalesce_float": return func(in []*Value)([]*Value) {
		if in[0].Kind == Type_float64 {
			return []*Value{Value_float64(in[0].Value_float64)}
		} else {
			return []*Value{Value_float64(in[1].Value_float64)}
		}
	}(vs)
	default: panic("unknown operator")
	}
}

var op_open *test = &test{
   kind: Kind_group_open,
	symbol: "(",
}

var op_close *test = &test{
   kind: Kind_group_close,
	symbol: ")",
}

var op_add *test = &test{
   precedence: 1,
   associativity: Associativity_left,
   needs_elements: 2,
   return_elements: 1,
   kind: Kind_operator,
	symbol: "+",
	input_types: [][]int{[]int{Type_float64},[]int{Type_float64}},
	output_types: [][]int{[]int{Type_float64}},
}

var op_coalesce_float *test = &test{
   precedence: 2,
   associativity: Associativity_left,
   needs_elements: 2,
   return_elements: 1,
   kind: Kind_operator,
	symbol: "coalesce_float",
	input_types: [][]int{[]int{Type_float64,Type_nil},[]int{Type_float64}},
	output_types: [][]int{[]int{Type_float64}},
}

var op_mul *test = &test{
   precedence: 2,
   associativity: Associativity_left,
   needs_elements: 2,
   return_elements: 1,
   kind: Kind_operator,
	symbol: "*",
	input_types: [][]int{[]int{Type_float64},[]int{Type_float64}},
	output_types: [][]int{[]int{Type_float64}},
}

var op_or *test = &test{
   precedence: 1,
   associativity: Associativity_left,
   needs_elements: 2,
   return_elements: 1,
   kind: Kind_operator,
	symbol: "or",
	input_types: [][]int{[]int{Type_bool},[]int{Type_bool}},
	output_types: [][]int{[]int{Type_bool}},
}

var op_and *test = &test{
   precedence: 2,
   associativity: Associativity_left,
   needs_elements: 2,
   return_elements: 1,
   kind: Kind_operator,
	symbol: "and",
	input_types: [][]int{[]int{Type_bool},[]int{Type_bool}},
	output_types: [][]int{[]int{Type_bool}},
}

var op_neg *test = &test{
   precedence: 3,
   associativity: Associativity_right,
   needs_elements: 1,
   return_elements: 1,
   kind: Kind_operator,
	symbol: "neg",
	input_types: [][]int{[]int{Type_float64}},
	output_types: [][]int{[]int{Type_float64}},
}

var op_true *test = &test{
   needs_elements: 0,
   return_elements: 1,
   kind: Kind_value,
	symbol: "true",
	input_types: [][]int{},
	output_types: [][]int{[]int{Type_bool}},
}

var op_false *test = &test{
   needs_elements: 0,
   return_elements: 1,
   kind: Kind_value,
	symbol: "false",
	input_types: [][]int{},
	output_types: [][]int{[]int{Type_bool}},
}

var op_23 *test = &test{
   needs_elements: 0,
   return_elements: 1,
   kind: Kind_value,
	symbol: "2.3",
	input_types: [][]int{},
	output_types: [][]int{[]int{Type_float64}},
}

var op_24 *test = &test{
   needs_elements: 0,
   return_elements: 1,
   kind: Kind_value,
	symbol: "2.4",
	input_types: [][]int{},
	output_types: [][]int{[]int{Type_float64}},
}

var op_25 *test = &test{
   needs_elements: 0,
   return_elements: 1,
   kind: Kind_value,
	symbol: "2.5",
	input_types: [][]int{},
	output_types: [][]int{[]int{Type_float64}},
}

var op_26 *test = &test{
   needs_elements: 0,
   return_elements: 1,
   kind: Kind_value,
	symbol: "2.6",
	input_types: [][]int{},
	output_types: [][]int{[]int{Type_float64}},
}

var op_26_or_nil *test = &test{
   needs_elements: 0,
   return_elements: 1,
   kind: Kind_value,
	symbol: "2.6_or_nil",
	input_types: [][]int{},
	output_types: [][]int{[]int{Type_float64,Type_nil}},
}

func verif(t *testing.T, e *Expr, expect string) {
	var sol string
	var ec *elt_cache

	for _, ec = range e.rpn {
		sol += ec.elt.String() + "|"
	}

	if sol != expect {
		t.Errorf("Expect %s, got %s", expect, sol)
	}
}

func test_panic(t *testing.T, f func()) {
	defer func() {
		var r interface{}

		r = recover()
		if r == nil {
			t.Errorf("Expect panic")
		}
	}()
	f()
}

func test_nopanic(t *testing.T, f func()) {
	defer func() {
		var r interface{}

		r = recover()
		if r != nil {
			t.Errorf("Expect success, got panic: %#v", r)
		}
	}()
	f()
}

func Test_value(t *testing.T) {
	var v *Value

	v = Value_bool(false)
	if v.Bool() {
		t.Errorf("Expect result false, got %t", v.Bool())
	}

	v = Value_bool(true)
	if !v.Bool() {
		t.Errorf("Expect result true, got %t", v.Bool())
	}

	v = Value_float64(0.0)
	if v.Bool() {
		t.Errorf("Expect result false, got %t", v.Bool())
	}

	v = Value_float64(1.0)
	if !v.Bool() {
		t.Errorf("Expect result true, got %t", v.Bool())
	}

	v = Value_nil()
	if v.Bool() {
		t.Errorf("Expect result false, got %t", v.Bool())
	}

	v = Value_float64(1.0)
	if v.Float64() != 1.0 {
		t.Errorf("Expect result 1.0, got %f", v.Float64())
	}

	v = Value_bool(false)
	if v.Float64() != 0.0 {
		t.Errorf("Expect result 0.0, got %f", v.Float64())
	}

	v = Value_bool(true)
	if v.Float64() != 1.0 {
		t.Errorf("Expect result 1.0, got %f", v.Float64())
	}

	v = Value_nil()
	if v.Float64() != 0.0 {
		t.Errorf("Expect result 0.0, got %f", v.Float64())
	}

	v = Value_float64(1.44)
	if v.String() != "1.440000" {
		t.Errorf("Expect result \"1.440000\", got %q", v.String())
	}

	v = Value_bool(true)
	if v.String() != "true" {
		t.Errorf("Expect result \"true\", got %q", v.String())
	}

	v = Value_nil()
	if v.String() != "nil" {
		t.Errorf("Expect result \"nil\", got %q", v.String())
	}

	test_panic(t, func(){
		v = Value_float64(0.0)
		v.Kind = 5000
		v.Bool()
	})

	test_panic(t, func(){
		v = Value_float64(0.0)
		v.Kind = 5000
		v.Float64()
	})

	test_panic(t, func(){
		v = Value_float64(0.0)
		v.Kind = 5000
		v.String()
	})
}

func Test_compat(t *testing.T) {

	if (has_compat([]int{1}, []int{2, 3, 4})) {
		t.Errorf("expect false, got true")
	}

	if (has_compat([]int{10, 20}, []int{10})) {
		t.Errorf("expect false, got true")
	}

	if (!has_compat([]int{10}, []int{10, 20})) {
		t.Errorf("expect true, got false")
	}

	if (!has_compat([]int{1,2,3}, []int{3,2,1})) {
		t.Errorf("expect true, got false")
	}

	if (has_compat([]int{10, 20}, []int{})) {
		t.Errorf("expect true, got false")
	}

	if (has_compat([]int{}, []int{10, 20})) {
		t.Errorf("expect true, got false")
	}

}

func Test_expr(t *testing.T) {
	var e *Expr
	var err error

	e = New()
	e.done = true
	err = e.Append(op_23)
	if err == nil {
		t.Errorf("Expect error, got no error")
	}
	err = e.Finalize()
	if err == nil {
		t.Errorf("Expect error, got no error")
	}

	e = New()
	e.Append(op_23)
	e.Append(op_add)
	e.Append(op_24)
	err = e.Finalize()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	verif(t, e, "2.3|2.4|+|")

	e = New()
	e.Append(op_23)
	e.Append(op_add)
	e.Append(op_24)
	e.Append(op_mul)
	e.Append(op_25)
	err = e.Finalize()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	verif(t, e, "2.3|2.4|2.5|*|+|")

	e = New()
	err = e.Append(op_close)
	if err == nil {
		t.Errorf("Expect error, got no error")
	}

	e = New()
	e.Append(op_23)
	e.Append(op_add)
	e.Append(op_24)
	e.Append(op_mul)
	e.Append(op_25)
	err = e.Append(op_close)
	if err == nil {
		t.Errorf("Expect error, got no error")
	}

	e = New()
	e.Append(op_23)
	e.Append(op_open)
	e.Append(op_add)
	e.Append(op_24)
	e.Append(op_mul)
	e.Append(op_25)
	err = e.Finalize()
	if err == nil {
		t.Errorf("Expect error, got no error")
	}

	e = New()
	e.Append(op_open)
	e.Append(op_23)
	e.Append(op_add)
	e.Append(op_24)
	e.Append(op_close)
	e.Append(op_mul)
	e.Append(op_25)
	err = e.Finalize()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	verif(t, e, "2.3|2.4|+|2.5|*|")

	e = New()
	e.Append(op_open)
	e.Append(op_23)
	e.Append(op_add)
	e.Append(op_24)
	e.Append(op_close)
	e.Append(op_mul)
	e.Append(op_25)
	e.Append(op_add)
	e.Append(op_26)
	err = e.Finalize()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	verif(t, e, "2.3|2.4|+|2.5|*|2.6|+|")

	e = New()
	e.Append(op_23)
	e.Append(op_add)
	e.Append(op_24)
	e.Append(op_add)
	e.Append(op_neg)
	e.Append(op_25)
	err = e.Finalize()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	verif(t, e, "2.3|2.4|+|2.5|neg|+|")

	e = New()
	err = e.Append(&test{kind: 8000})
	if err == nil {
		t.Errorf("Expect error, got no error")
	}

	e = New()
	err = e.Finalize()
	if err == nil {
		t.Errorf("Expect error, got no error")
	}

	e = New()
	e.Append(op_23)
	e.Append(op_23)
	err = e.Finalize()
	if err == nil {
		t.Errorf("Expect error, got no error")
	}

	e = New()
	e.Append(op_add)
	err = e.Finalize()
	if err == nil {
		t.Errorf("Expect error, got no error")
	}

	e = New()
	e.Append(op_23)
	e.Append(op_add)
	e.Append(op_true)
	err = e.Finalize()
	if err == nil {
		t.Errorf("Expected error")
	}

	e = New()
	e.Append(op_26_or_nil)
	e.Append(op_coalesce_float)
	e.Append(op_26)
	err = e.Finalize()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	e = New()
	e.Append(op_26)
	e.Append(op_coalesce_float)
	e.Append(op_26_or_nil)
	err = e.Finalize()
	if err == nil {
		t.Errorf("Expect error")
	}
}

func Test_exec(t *testing.T) {
	var e *Expr
	var err error
	var v *Value

	e = New()
	e.Append(op_23)
	e.Append(op_add)
	e.Append(op_24)
	e.Append(op_mul)
	e.Append(op_25)
	err = e.Finalize()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	v, err = e.Exec()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	} else {
		if v.Float64() != 8.3 {
			t.Errorf("Expect result 8.3, got %f", v.Float64())
		}
	}

	e = New()
	e.Append(op_add)
	e.Append(op_add)
	e.Finalize() // error intentionnaly not check
	e.done = true // force done
	_, err = e.Exec()
	if err == nil {
		t.Errorf("Expect error, got no error")
	}

	e = New()
	e.Append(op_23)
	e.Append(op_24)
	e.Finalize() // error intentionnaly not check
	e.done = true // force done
	_, err = e.Exec()
	if err == nil {
		t.Errorf("Expect error, got no error")
	}
	
	e = New()
	e.Append(op_true)
	e.Append(op_and)
	e.Append(op_false)
	e.Append(op_or)
	e.Append(op_true)
	err = e.Finalize()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	v, err = e.Exec()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	} else {
		if !v.Bool() {
			t.Errorf("Expect result true, got %t", v.Bool())
		}
	}

	test_nopanic(t, func() {
		var e *Expr

		e = New()
		e.Append(op_true)
		e.Append(op_and)
		e.Append(op_false)
		e.Append(op_or)
		e.Append(op_true)
		e.Finalize()
		e.Dump()
	})
}

func Test_type_desc(t *testing.T) {
	var types []int
	var res string

	types = []int{}
	res = Type_desc(types)
	if res != "void" {
		t.Errorf("expect \"void\", got %q", res)
	}

	types = []int{Type_bool}
	res = Type_desc(types)
	if res != "bool" {
		t.Errorf("expect \"bool\", got %q", res)
	}

	types = []int{Type_bool, Type_float64, Type_nil}
	res = Type_desc(types)
	if res != "bool|float64|nil" {
		t.Errorf("expect \"bool|float64|nil\", got %q", res)
	}

	test_panic(t, func() {
		Type_desc([]int{9000})
	})
}
