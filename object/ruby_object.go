package object

import (
	"bytes"
	"strings"

	"github.com/goruby/goruby/ast"
)

type ObjectType string

const (
	FUNCTION_OBJ           = "FUNCTION"
	RETURN_VALUE_OBJ       = "RETURN_VALUE"
	BASIC_OBJECT_OBJ       = "BASIC_OBJECT"
	BASIC_OBJECT_CLASS_OBJ = "BASIC_OBJECT_CLASS"
	INTEGER_OBJ            = "INTEGER"
	INTEGER_CLASS_OBJ      = "INTEGER_CLASS"
	STRING_OBJ             = "STRING"
	SYMBOL_OBJ             = "SYMBOL"
	BOOLEAN_OBJ            = "BOOLEAN"
	NIL_OBJ                = "NIL"
	ERROR_OBJ              = "ERROR"
	EXCEPTION_OBJ          = "EXCEPTION"
	BUILTIN_OBJ            = "BUILTIN"
)

var (
	NIL   = &Nil{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

type RubyObject interface {
	Type() ObjectType
	Send(name string, args ...RubyObject) RubyObject
	Inspect() string
}

type BuiltinFunction func(args ...RubyObject) RubyObject

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType                                { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string                                 { return "builtin function" }
func (b *Builtin) Send(name string, args ...RubyObject) RubyObject { return NIL }

type Nil struct{}

func (n *Nil) Inspect() string                                 { return "nil" }
func (n *Nil) Type() ObjectType                                { return NIL_OBJ }
func (n *Nil) Send(name string, args ...RubyObject) RubyObject { return n }

type ReturnValue struct {
	Value RubyObject
}

func (rv *ReturnValue) Type() ObjectType                                { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string                                 { return rv.Value.Inspect() }
func (rv *ReturnValue) Send(name string, args ...RubyObject) RubyObject { return NIL }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType                                { return ERROR_OBJ }
func (e *Error) Inspect() string                                 { return "ERROR: " + e.Message }
func (e *Error) Send(name string, args ...RubyObject) RubyObject { return NIL }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}
func (f *Function) Send(name string, args ...RubyObject) RubyObject { return NIL }