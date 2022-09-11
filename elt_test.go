package shuntingyard

import "fmt"
import "reflect"
import "testing"

type test struct {
	precedence int
	associativity int
	kind int
	symbol string
	input_types [][]int
	output_types [][]int
}

func (t *test)Precedence()(int) { return t.precedence }
func (t *test)Associativity()(int) { return t.associativity }
func (t *test)Kind()(int) { return t.kind }
func (t *test)String()(string) { return t.symbol }
func (t *test)Input_types()([][]int) { return t.input_types }
func (t *test)Output_types()([][]int) { return t.output_types }
func (t *test)Execute(vs []*Value)([]*Value, error) {
	switch t.symbol {
	case "true": return []*Value{Value_bool(true)}, nil
	case "false": return []*Value{Value_bool(false)}, nil
	case "or": return []*Value{Value_bool(vs[0].Value_bool || vs[1].Value_bool)}, nil
	case "and": return []*Value{Value_bool(vs[0].Value_bool && vs[1].Value_bool)}, nil
	case "not":  return []*Value{Value_bool(!vs[0].Value_bool)}, nil
	case "neg":  return []*Value{Value_float64(-vs[0].Value_float64)}, nil
	case "*":  return []*Value{Value_float64(vs[0].Value_float64 * vs[1].Value_float64)}, nil
	case "/":  return []*Value{Value_float64(vs[0].Value_float64 / vs[1].Value_float64)}, nil
	case "+":  return []*Value{Value_float64(vs[0].Value_float64 + vs[1].Value_float64)}, nil
	case "+e": return nil, fmt.Errorf("This is a %s error", t.symbol)
	case "-":  return []*Value{Value_float64(vs[0].Value_float64 - vs[1].Value_float64)}, nil
	case "2.3":  return []*Value{Value_float64(2.3)}, nil
	case "2.4":  return []*Value{Value_float64(2.4)}, nil
	case "2.5":  return []*Value{Value_float64(2.5)}, nil
	case "2.6":  return []*Value{Value_float64(2.6)}, nil
	case "2.6_or_nil":  return []*Value{Value_float64(2.6)}, nil
	case "coalesce_float": return func(in []*Value)([]*Value, error) {
		if in[0].Kind == Type_float64 {
			return []*Value{Value_float64(in[0].Value_float64)}, nil
		} else {
			return []*Value{Value_float64(in[1].Value_float64)}, nil
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

var op_add_error *test = &test{
	precedence: 1,
	associativity: Associativity_left,
	kind: Kind_operator,
	symbol: "+e",
	input_types: [][]int{[]int{Type_float64},[]int{Type_float64}},
	output_types: [][]int{[]int{Type_float64}},
}

var op_add *test = &test{
	precedence: 1,
	associativity: Associativity_left,
	kind: Kind_operator,
	symbol: "+",
	input_types: [][]int{[]int{Type_float64},[]int{Type_float64}},
	output_types: [][]int{[]int{Type_float64}},
}

var op_coalesce_float *test = &test{
	precedence: 2,
	associativity: Associativity_left,
	kind: Kind_operator,
	symbol: "coalesce_float",
	input_types: [][]int{[]int{Type_float64,Type_nil},[]int{Type_float64}},
	output_types: [][]int{[]int{Type_float64}},
}

var op_mul *test = &test{
	precedence: 2,
	associativity: Associativity_left,
	kind: Kind_operator,
	symbol: "*",
	input_types: [][]int{[]int{Type_float64},[]int{Type_float64}},
	output_types: [][]int{[]int{Type_float64}},
}

var op_or *test = &test{
	precedence: 1,
	associativity: Associativity_left,
	kind: Kind_operator,
	symbol: "or",
	input_types: [][]int{[]int{Type_bool},[]int{Type_bool}},
	output_types: [][]int{[]int{Type_bool}},
}

var op_and *test = &test{
	precedence: 2,
	associativity: Associativity_left,
	kind: Kind_operator,
	symbol: "and",
	input_types: [][]int{[]int{Type_bool},[]int{Type_bool}},
	output_types: [][]int{[]int{Type_bool}},
}

var op_neg *test = &test{
	precedence: 3,
	associativity: Associativity_right,
	kind: Kind_operator,
	symbol: "neg",
	input_types: [][]int{[]int{Type_float64}},
	output_types: [][]int{[]int{Type_float64}},
}

var op_true *test = &test{
	kind: Kind_value,
	symbol: "true",
	input_types: [][]int{},
	output_types: [][]int{[]int{Type_bool}},
}

var op_false *test = &test{
	kind: Kind_value,
	symbol: "false",
	input_types: [][]int{},
	output_types: [][]int{[]int{Type_bool}},
}

var op_23 *test = &test{
	kind: Kind_value,
	symbol: "2.3",
	input_types: [][]int{},
	output_types: [][]int{[]int{Type_float64}},
}

var op_24 *test = &test{
	kind: Kind_value,
	symbol: "2.4",
	input_types: [][]int{},
	output_types: [][]int{[]int{Type_float64}},
}

var op_25 *test = &test{
	kind: Kind_value,
	symbol: "2.5",
	input_types: [][]int{},
	output_types: [][]int{[]int{Type_float64}},
}

var op_26 *test = &test{
	kind: Kind_value,
	symbol: "2.6",
	input_types: [][]int{},
	output_types: [][]int{[]int{Type_float64}},
}

var op_26_or_nil *test = &test{
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

	v = Value_bool(true)
	if v.Kind != Type_bool {
		t.Errorf("Expect type bool, got %s", Type_string(v.Kind))
	}
	if !v.Value_bool {
		t.Errorf("Expect value true, got %t", v.Value_bool)
	}

	v = Value_float64(1.2)
	if v.Kind != Type_float64 {
		t.Errorf("Expect type float64, got %s", Type_string(v.Kind))
	}
	if v.Value_float64 != 1.2 {
		t.Errorf("Expect value true, got %t", v.Value_bool)
	}

	v = Value_nil()
	if v.Kind != Type_nil {
		t.Errorf("Expect type nil, got %s", Type_string(v.Kind))
	}
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

	e = New(nil)
	e.done = true
	err = e.Append(op_23)
	if err == nil {
		t.Errorf("Expect error, got no error")
	}
	err = e.Finalize()
	if err == nil {
		t.Errorf("Expect error, got no error")
	}

	e = New(nil)
	e.Append(op_23)
	e.Append(op_add)
	e.Append(op_24)
	err = e.Finalize()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	verif(t, e, "2.3|2.4|+|")

	e = New(nil)
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

	e = New(nil)
	err = e.Append(op_close)
	if err == nil {
		t.Errorf("Expect error, got no error")
	}

	e = New(nil)
	e.Append(op_23)
	e.Append(op_add)
	e.Append(op_24)
	e.Append(op_mul)
	e.Append(op_25)
	err = e.Append(op_close)
	if err == nil {
		t.Errorf("Expect error, got no error")
	}

	e = New(nil)
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

	e = New(nil)
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

	e = New(nil)
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

	e = New(nil)
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

	e = New(nil)
	err = e.Append(&test{kind: 8000})
	if err == nil {
		t.Errorf("Expect error, got no error")
	}

	e = New(nil)
	err = e.Finalize()
	if e.input_types != nil {
		t.Errorf("Expect empty input types, got %q", Type_list(e.input_types))
	}
	if e.output_types != nil {
		t.Errorf("Expect empty output types, got %q", Type_list(e.output_types))
	}

	e = New([][]int{[]int{Type_bool}, []int{Type_float64}})
	err = e.Finalize()
	if !reflect.DeepEqual(e.input_types, [][]int{[]int{Type_bool}, []int{Type_float64}}) {
		t.Errorf("Expect empty input types, got %q", Type_list(e.input_types))
	}
	if !reflect.DeepEqual(e.output_types, [][]int{[]int{Type_bool}, []int{Type_float64}}) {
		t.Errorf("Expect empty output types, got %q", Type_list(e.output_types))
	}

	e = New(nil)
	e.Append(op_23)
	e.Append(op_23)
	err = e.Finalize()
	if !reflect.DeepEqual(e.output_types, [][]int{[]int{Type_float64}, []int{Type_float64}}) {
		t.Errorf("Expect \"float64, float64\" output types, got %q", Type_list(e.output_types))
	}

	e = New(nil)
	e.Append(op_add)
	err = e.Finalize()
	if err == nil {
		t.Errorf("Expect error, got no error")
	}

	e = New(nil)
	e.Append(op_23)
	e.Append(op_add)
	e.Append(op_true)
	err = e.Finalize()
	if err == nil {
		t.Errorf("Expected error")
	}

	e = New(nil)
	e.Append(op_26_or_nil)
	e.Append(op_coalesce_float)
	e.Append(op_26)
	err = e.Finalize()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	e = New(nil)
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

	e = New(nil)
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
		if v.Kind != Type_float64 {
			t.Errorf("Expect float, got %s", Type_string(v.Kind))
		}
		if v.Value_float64 != 8.3 {
			t.Errorf("Expect result 8.3, got %f", v.Value_float64)
		}
	}

	e = New(nil)
	e.Append(op_add)
	e.Append(op_add)
	e.Finalize() // error intentionnaly not check
	e.done = true // force done
	_, err = e.Exec()
	if err == nil {
		t.Errorf("Expect error, got no error")
	}

	e = New(nil)
	e.Append(op_23)
	e.Append(op_24)
	e.Finalize() // error intentionnaly not check
	e.done = true // force done
	_, err = e.Exec()
	if err == nil {
		t.Errorf("Expect error, got no error")
	}
	
	e = New(nil)
	e.Append(op_23)
	e.Append(op_add_error)
	e.Append(op_24)
	err = e.Finalize()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	_, err = e.Exec()
	if err == nil {
		t.Errorf("Expect error")
	}

	e = New(nil)
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
		if v.Kind != Type_bool {
			t.Errorf("Expect bool, got %s", Type_string(v.Kind))
		}
		if !v.Value_bool {
			t.Errorf("Expect result true, got %t", v.Value_bool)
		}
	}

	test_nopanic(t, func() {
		var e *Expr

		e = New(nil)
		e.Append(op_true)
		e.Append(op_and)
		e.Append(op_false)
		e.Append(op_or)
		e.Append(op_true)
		e.Finalize()
		e.Dump() /* test dump function */
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

	res = Type_string(Type_bool)
	if res != "bool" {
		t.Errorf("expect \"bool\", got %q", res)
	}

	res = Type_list([][]int{[]int{Type_float64, Type_nil},[]int{Type_bool}})
	if res != "float64|nil, bool" {
		t.Errorf("expect \"float64|nil, bool\", got %q", res)
	}

	test_panic(t, func() {
		Type_desc([]int{9000})
	})
}

func Test_sub_expression(t *testing.T) {
	var se *Expr
	var err error

	se = New([][]int{[]int{Type_float64}})
	se.Push(op_25)
	se.Push(op_add)
	se.Finalize()
	err = se.Push(op_26)
	if err == nil {
		t.Errorf("Expect error, got no error")
	}

	se = New(nil)
	se.Push(op_25)
	se.Push(op_add)
	se.Push(op_26)
	err = se.Finalize()
	if err == nil {
		t.Errorf("Expect error, got no error")
	}

	se = New(nil)
	se.Push(op_25)
	se.Push(op_26)
	se.Push(op_add)
	err = se.Finalize()
	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
	}
	if se.input_types != nil {
		t.Errorf("Expect empty input type")
	}
	if !reflect.DeepEqual(se.output_types, [][]int{[]int{Type_float64}}) {
		t.Errorf("Expect \"float64\" output types, got %q", Type_list(se.output_types))
	}
}
