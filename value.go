package shuntingyard

// This interface describe a value. The value is a value used
// for calculus. The value is associated with a type.
type Value interface {
	Descr()(string)
	Type()(Type)
}
