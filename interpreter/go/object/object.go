package object

import (
	"fmt"
	"hash/fnv"
	"kaze/ast"
	"strconv"
	"strings"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Printable interface {
	Object
	String() string
}

const (
	ERROR_OBJ    = "ERROR"
	NULL_OBJ     = "NULL"
	INTEGER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	STRING_OBJ   = "STRING"
	RETURN_OBJ   = "RETURN"
	FUNCTION_OBJ = "FUNCTION"
	BREAK_OBJ    = "BREAK"
	CONTINUE_OBJ = "CONTINUE"
	BUILTIN_OBJ  = "BUILTIN"
	HASH_OBJ     = "HASH"
	ARRAY_OBJ    = "ARRAY"
	LVALUE_OBJ   = "LVALUE"
)

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) String() string {
	return e.Inspect()
}

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }
func (n *Null) String() string {
	return n.Inspect()
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return strconv.FormatInt(i.Value, 10) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}
func (i *Integer) String() string {
	return i.Inspect()
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string {
	if b.Value {
		return "true"
	}
	return "false"
}
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}
func (b *Boolean) String() string {
	return b.Inspect()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return fmt.Sprintf(`"%s"`, s.Value) }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
func (s *String) String() string {
	return s.Value
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Function struct {
	Parameters []*ast.Identifier
	Body       ast.Expression
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	return fmt.Sprintf("fn(%s) {\n%s\n}", strings.Join(params, ", "), f.Body.String())
}
func (f *Function) String() string {
	return f.Inspect()
}

type Break struct{}

func (b *Break) Type() ObjectType { return BREAK_OBJ }
func (b *Break) Inspect() string  { return "break" }

type Continue struct{}

func (c *Continue) Type() ObjectType { return CONTINUE_OBJ }
func (c *Continue) Inspect() string  { return "continue" }

type Builtin struct {
	Fn func(args ...Object) Object
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) String() string {
	return b.Inspect()
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	pairs := make([]string, 0)
	for _, pair := range h.Pairs {
		pairs = append(pairs, pair.Key.Inspect()+": "+pair.Value.Inspect())
	}

	return "#{ " + strings.Join(pairs, ", ") + " }"
}
func (h *Hash) String() string {
	return h.Inspect()[1:]
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	elements := make([]string, 0)
	for _, el := range a.Elements {
		elements = append(elements, el.Inspect())
	}

	return "[ " + strings.Join(elements, ", ") + " ]"
}
func (a *Array) String() string {
	return a.Inspect()
}

type LValue interface {
	Object
	Get() (Object, bool)
	Update(Object) (Object, bool)
}

type Variable struct {
	Name string
	Env  *Environment
}

func (v *Variable) Type() ObjectType {
	return LVALUE_OBJ
}
func (v *Variable) Inspect() string {
	return LVALUE_OBJ
}
func (v *Variable) Get() (Object, bool) {
	return v.Env.Get(v.Name)
}
func (v *Variable) Update(val Object) (Object, bool) {
	return v.Env.Update(v.Name, val)
}

type IndexRef struct {
	Left  LValue
	Index Object
}

func (ir *IndexRef) Type() ObjectType {
	return LVALUE_OBJ
}
func (ir *IndexRef) Inspect() string {
	return ir.Left.Inspect()
}
func (ir *IndexRef) Get() (Object, bool) {
	left, ok := ir.Left.Get()
	if !ok {
		return nil, false
	}

	switch obj := left.(type) {
	case *Array:
		index, ok := ir.Index.(*Integer)
		if !ok || index.Value < 0 || int(index.Value) >= len(obj.Elements) {
			return nil, false
		}

		return obj.Elements[index.Value], true
	case *Hash:
		hashKey, ok := ir.Index.(Hashable)
		if !ok {
			return nil, false
		}

		key := hashKey.HashKey()
		pair, ok := obj.Pairs[key]
		if !ok {
			return nil, false
		}

		return pair.Value, true
	case *String:
		index, ok := ir.Index.(*Integer)
		if !ok || index.Value < 0 || int(index.Value) >= len(obj.Value) {
			return nil, false
		}

		return &String{Value: string(obj.Value[index.Value])}, true
	}
	return nil, false
}
func (ir *IndexRef) Update(val Object) (Object, bool) {
	left, ok := ir.Left.Get()
	if !ok {
		return nil, false
	}

	switch obj := left.(type) {
	case *Array:
		index, ok := ir.Index.(*Integer)
		if !ok || index.Value < 0 || int(index.Value) >= len(obj.Elements) {
			return nil, false
		}

		obj.Elements[index.Value] = val
		return val, true
	case *Hash:
		hashKey, ok := ir.Index.(Hashable)
		if !ok {
			return nil, false
		}

		key := hashKey.HashKey()
		obj.Pairs[key] = HashPair{Key: ir.Index, Value: val}
		return val, true
	case *String:
		index, ok := ir.Index.(*Integer)
		if !ok || index.Value < 0 || int(index.Value) >= len(obj.Value) {
			return nil, false
		}

		val, ok := val.(*String)
		if !ok {
			return nil, false
		}

		if len(val.Value) != 1 {
			return nil, false
		}

		obj.Value = obj.Value[:index.Value] + val.Value + obj.Value[index.Value+1:]
		return val, true
	}

	return nil, false
}
