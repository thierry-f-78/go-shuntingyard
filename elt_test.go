package shuntingyard

import "context"
import "fmt"
import "reflect"
import "testing"

type test struct {
	precedence int
	associativity int
	kind int
	symbol string
	input_types [][]Type
	output_types [][]Type
}

type type_t int

func (t type_t) Name()(string) {
	switch t {
	case type_nil: return "nil"
	case type_bool: return "bool"
	case type_float64: return "float64"
	}
	panic(fmt.Errorf("Unhandled type %#v", t))
}

var type_float64 type_t = 0
var type_bool type_t = 1
var type_nil type_t = 2
var type_other type_t = 3

type value_t struct {
	// type of value. use Type_*
	kind type_t
	// Float64
	value_float64 float64
	// Bool
	value_bool bool
}

func (v *value_t)Descr()(string) {
	switch v.kind {
	case type_nil:     return "nil"
	case type_bool:    return fmt.Sprintf("%t", v.value_bool)
	case type_float64: return fmt.Sprintf("%f", v.value_float64)
	}
	panic(fmt.Errorf("Unhandled type %#v", v.kind))
}

func (v *value_t)Type()(Type) {
	return v.kind
}

func value_nil()(Value) {
	return &value_t{
		kind: type_t(type_nil),
	}
}

func value_bool(v bool)(Value) {
	return &value_t{
		kind: type_t(type_bool),
		value_bool: v,
	}
}

func value_float64(v float64)(Value) {
	return &value_t{
		kind: type_t(type_float64),
		value_float64: v,
	}
}

func (t *test)Precedence()(int) { return t.precedence }
func (t *test)Associativity()(int) { return t.associativity }
func (t *test)Kind()(int) { return t.kind }
func (t *test)String()(string) { return t.symbol }
func (t *test)Input_types()([][]Type) { return t.input_types }
func (t *test)Output_types()([][]Type) { return t.output_types }
func (t *test)Execute(ctx context.Context, vs []Value)([]Value, error) {
	switch t.symbol {
	case "true": return []Value{value_bool(true)}, nil
	case "false": return []Value{value_bool(false)}, nil
	case "or": return []Value{value_bool(vs[0].(*value_t).value_bool || vs[1].(*value_t).value_bool)}, nil
	case "and": return []Value{value_bool(vs[0].(*value_t).value_bool && vs[1].(*value_t).value_bool)}, nil
	case "not":  return []Value{value_bool(!vs[0].(*value_t).value_bool)}, nil
	case "neg":  return []Value{value_float64(-vs[0].(*value_t).value_float64)}, nil
	case "*":  return []Value{value_float64(vs[0].(*value_t).value_float64 * vs[1].(*value_t).value_float64)}, nil
	case "/":  return []Value{value_float64(vs[0].(*value_t).value_float64 / vs[1].(*value_t).value_float64)}, nil
	case "+":  return []Value{value_float64(vs[0].(*value_t).value_float64 + vs[1].(*value_t).value_float64)}, nil
	case "+e": return nil, fmt.Errorf("This is a %s error", t.symbol)
	case "-":  return []Value{value_float64(vs[0].(*value_t).value_float64 - vs[1].(*value_t).value_float64)}, nil
	case "2.3":  return []Value{value_float64(2.3)}, nil
	case "2.4":  return []Value{value_float64(2.4)}, nil
	case "2.5":  return []Value{value_float64(2.5)}, nil
	case "2.6":  return []Value{value_float64(2.6)}, nil
	case "2.6_or_nil":  return []Value{value_float64(2.6)}, nil
	case "coalesce_float": return func(in []Value)([]Value, error) {
		if in[0].(*value_t).kind == type_float64 {
			return []Value{value_float64(in[0].(*value_t).value_float64)}, nil
		} else {
			return []Value{value_float64(in[1].(*value_t).value_float64)}, nil
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
	input_types: [][]Type{[]Type{type_float64},[]Type{type_float64}},
	output_types: [][]Type{[]Type{type_float64}},
}

var op_add *test = &test{
	precedence: 1,
	associativity: Associativity_left,
	kind: Kind_operator,
	symbol: "+",
	input_types: [][]Type{[]Type{type_float64},[]Type{type_float64}},
	output_types: [][]Type{[]Type{type_float64}},
}

var op_coalesce_float *test = &test{
	precedence: 2,
	associativity: Associativity_left,
	kind: Kind_operator,
	symbol: "coalesce_float",
	input_types: [][]Type{[]Type{type_float64,type_nil},[]Type{type_float64}},
	output_types: [][]Type{[]Type{type_float64}},
}

var op_mul *test = &test{
	precedence: 2,
	associativity: Associativity_left,
	kind: Kind_operator,
	symbol: "*",
	input_types: [][]Type{[]Type{type_float64},[]Type{type_float64}},
	output_types: [][]Type{[]Type{type_float64}},
}

var op_or *test = &test{
	precedence: 1,
	associativity: Associativity_left,
	kind: Kind_operator,
	symbol: "or",
	input_types: [][]Type{[]Type{type_bool},[]Type{type_bool}},
	output_types: [][]Type{[]Type{type_bool}},
}

var op_and *test = &test{
	precedence: 2,
	associativity: Associativity_left,
	kind: Kind_operator,
	symbol: "and",
	input_types: [][]Type{[]Type{type_bool},[]Type{type_bool}},
	output_types: [][]Type{[]Type{type_bool}},
}

var op_neg *test = &test{
	precedence: 3,
	associativity: Associativity_right,
	kind: Kind_operator,
	symbol: "neg",
	input_types: [][]Type{[]Type{type_float64}},
	output_types: [][]Type{[]Type{type_float64}},
}

var op_true *test = &test{
	kind: Kind_value,
	symbol: "true",
	input_types: [][]Type{},
	output_types: [][]Type{[]Type{type_bool}},
}

var op_false *test = &test{
	kind: Kind_value,
	symbol: "false",
	input_types: [][]Type{},
	output_types: [][]Type{[]Type{type_bool}},
}

var op_23 *test = &test{
	kind: Kind_value,
	symbol: "2.3",
	input_types: [][]Type{},
	output_types: [][]Type{[]Type{type_float64}},
}

var op_24 *test = &test{
	kind: Kind_value,
	symbol: "2.4",
	input_types: [][]Type{},
	output_types: [][]Type{[]Type{type_float64}},
}

var op_25 *test = &test{
	kind: Kind_value,
	symbol: "2.5",
	input_types: [][]Type{},
	output_types: [][]Type{[]Type{type_float64}},
}

var op_26 *test = &test{
	kind: Kind_value,
	symbol: "2.6",
	input_types: [][]Type{},
	output_types: [][]Type{[]Type{type_float64}},
}

var op_26_or_nil *test = &test{
	kind: Kind_value,
	symbol: "2.6_or_nil",
	input_types: [][]Type{},
	output_types: [][]Type{[]Type{type_float64,type_nil}},
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

func Test_compat(t *testing.T) {

	if (Has_compat([]Type{type_float64}, []Type{type_bool, type_nil, type_other})) {
		t.Errorf("expect false, got true")
	}

	if (Has_compat([]Type{type_float64, type_bool}, []Type{type_float64})) {
		t.Errorf("expect false, got true")
	}

	if (!Has_compat([]Type{type_float64}, []Type{type_float64, type_bool})) {
		t.Errorf("expect true, got false")
	}

	if (!Has_compat([]Type{type_float64,type_bool,type_nil}, []Type{type_nil,type_bool,type_float64})) {
		t.Errorf("expect true, got false")
	}

	if (Has_compat([]Type{type_float64, type_bool}, []Type{})) {
		t.Errorf("expect true, got false")
	}

	if (Has_compat([]Type{}, []Type{type_float64, type_bool})) {
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

	e = New([][]Type{[]Type{type_bool}, []Type{type_float64}})
	err = e.Finalize()
	if !reflect.DeepEqual(e.input_types, [][]Type{[]Type{type_bool}, []Type{type_float64}}) {
		t.Errorf("Expect empty input types, got %q", Type_list(e.input_types))
	}
	if !reflect.DeepEqual(e.output_types, [][]Type{[]Type{type_bool}, []Type{type_float64}}) {
		t.Errorf("Expect empty output types, got %q", Type_list(e.output_types))
	}

	e = New(nil)
	e.Append(op_23)
	e.Append(op_23)
	err = e.Finalize()
	if !reflect.DeepEqual(e.output_types, [][]Type{[]Type{type_float64}, []Type{type_float64}}) {
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
	var v []Value
	var ctx context.Context

	ctx = context.Background()

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
	v, err = e.Execute(ctx, nil)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	} else {
		if len(v) != 1 {
			t.Errorf("Expect only one value as return, got %d", len(v))
		} else {
			if v[0].(*value_t).kind != type_float64 {
				t.Errorf("Expect float, got %s", Type_string(v[0].(*value_t).kind))
			}
			if v[0].(*value_t).value_float64 != 8.3 {
				t.Errorf("Expect result 8.3, got %f", v[0].(*value_t).value_float64)
			}
		}
	}

	e = New(nil)
	e.Append(op_add)
	e.Append(op_add)
	e.Finalize() // error intentionnaly not check
	e.done = true // force done
	_, err = e.Execute(ctx, nil)
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
	_, err = e.Execute(ctx, nil)
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
	v, err = e.Execute(ctx, nil)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	} else {
		if len(v) != 1 {
			t.Errorf("Expect only one value as return, got %d", len(v))
		} else {
			if v[0].(*value_t).kind != type_bool {
				t.Errorf("Expect bool, got %s", Type_string(v[0].(*value_t).kind))
			}
			if !v[0].(*value_t).value_bool {
				t.Errorf("Expect result true, got %t", v[0].(*value_t).value_bool)
			}
		}
	}

	test_nopanic(t, func() {
		var e *Expr
		var e2 *Expr

		e2 = New(nil)
		e2.Append(op_true)
		e2.Append(op_and)
		e2.Append(op_false)
		e2.Append(op_or)
		e2.Append(op_true)
		e2.Finalize()

		e = New(nil)
		e.Append(op_true)
		e.Append(op_and)
		e.Append(e2)
		e.Append(op_or)
		e.Append(op_true)
		e.Finalize()
		e.Dump() /* test dump function */
	})
}

func Test_type_desc(t *testing.T) {
	var types []Type
	var res string

	types = []Type{}
	res = Type_desc(types)
	if res != "void" {
		t.Errorf("expect \"void\", got %q", res)
	}

	types = []Type{type_bool}
	res = Type_desc(types)
	if res != "bool" {
		t.Errorf("expect \"bool\", got %q", res)
	}

	types = []Type{type_bool, type_float64, type_nil}
	res = Type_desc(types)
	if res != "bool|float64|nil" {
		t.Errorf("expect \"bool|float64|nil\", got %q", res)
	}

	res = Type_string(type_bool)
	if res != "bool" {
		t.Errorf("expect \"bool\", got %q", res)
	}

	res = Type_list([][]Type{[]Type{type_float64, type_nil},[]Type{type_bool}})
	if res != "float64|nil, bool" {
		t.Errorf("expect \"float64|nil, bool\", got %q", res)
	}

	test_panic(t, func() {
		Type_desc([]Type{type_other})
	})
}

func Test_name(t *testing.T) {
	var se *Expr
	var err error

	se = New(nil)
	se.Push(op_25)
	se.Push(op_26)
	se.Push(op_add)
	err = se.Finalize()
	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
	}
	if se.String() != "2.5 2.6 +" {
		t.Errorf("Expect \"2.5 2.6 +\", got %q", se.String())
	}

	se = New(nil)
	se.Append(op_25)
	se.Append(op_add)
	se.Append(op_26)
	err = se.Finalize()
	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
	}
	if se.String() != "2.5 + 2.6" {
		t.Errorf("Expect \"2.5 + 2.6\", got %q", se.String())
	}

	se = New(nil)
	se.Set_name("name")
	err = se.Finalize()
	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
	}
	if se.String() != "name" {
		t.Errorf("Expect \"name\", got %q", se.String())
	}
}

func Test_sub_expression(t *testing.T) {
	var e *Expr
	var se *Expr
	var err error
	var v []Value
	var ctx context.Context

	ctx = context.Background()

	se = New([][]Type{[]Type{type_float64}})
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
	if !reflect.DeepEqual(se.output_types, [][]Type{[]Type{type_float64}}) {
		t.Errorf("Expect \"float64\" output types, got %q", Type_list(se.output_types))
	}

	e = New(nil)
	e.Append(se)
	e.Append(op_mul)
	e.Append(se)
	err = e.Finalize()
	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
	}
	v, err = e.Execute(ctx, nil)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	} else {
		if len(v) != 1 {
			t.Errorf("Expect only one value as return, got %d", len(v))
		} else {
			if v[0].(*value_t).kind != type_float64 {
				t.Errorf("Expect float, got %s", Type_string(v[0].(*value_t).kind))
			}
			/* 26.01 equality not match with floats */
			if v[0].(*value_t).value_float64 < 26.0099999 && v[0].(*value_t).value_float64 > 26.01 {
				t.Errorf("Expect result 26.01, got %f", v[0].(*value_t).value_float64)
			}
		}
	}
}
