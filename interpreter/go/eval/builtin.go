package eval

import (
	"fmt"
	"kaze/object"
	"os"
	"strconv"
	"strings"
)

var builtins = map[string]*object.Builtin{
	"print": {
		Fn: func(args ...object.Object) object.Object {
			var _args []string
			for _, arg := range args {
				if arg, ok := arg.(object.Printable); ok {
					_args = append(_args, arg.String())
				} else {
					return &object.Error{Message: fmt.Sprintf("cannot print type: %s", arg.Type())}
				}
			}
			fmt.Print(strings.Join(_args, " "))
			return NULL
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			var _args []string

			for _, arg := range args {
				if arg, ok := arg.(object.Printable); ok {
					_args = append(_args, arg.String())
				} else {
					return &object.Error{Message: fmt.Sprintf("cannot print type: %s", arg.Type())}
				}
			}

			fmt.Println(strings.Join(_args, " "))
			return NULL
		},
	},
	"args": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=0", len(args))}
			}

			var _args []object.Object
			for _, arg := range os.Args {
				_args = append(_args, &object.String{Value: arg})
			}
			return &object.Array{Elements: _args}
		},
	},
	"string": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			if arg, ok := args[0].(object.Printable); ok {
				return &object.String{Value: arg.String()}
			}
			return &object.Error{Message: fmt.Sprintf("cannot convert type: %s to string", args[0].Type())}
		},
	},
	"int": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				return arg
			case *object.String:
				value, ok := strconv.ParseInt(arg.Value, 10, 64)
				if ok != nil {
					return NAN
				}
				return &object.Integer{Value: value}
			default:
				return NAN
			}
		},
	},
	"ord": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.String:
				if len(arg.Value) != 1 {
					return NAN
				}
				return &object.Integer{Value: int64(arg.Value[0])}
			default:
				return NAN
			}
		},
	},
	"chr": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				return &object.String{Value: string(rune(arg.Value))}
			default:
				return NAN
			}
		},
	},
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.Hash:
				return &object.Integer{Value: int64(len(arg.Pairs))}
			default:
				return &object.Error{Message: fmt.Sprintf("argument to `len` not supported, got %s", args[0].Type())}
			}
		},
	},
	"append": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args))}
			}
			if arg, ok := args[0].(*object.Array); ok {
				return &object.Array{Elements: append(arg.Elements, args[1])}
			}
			return &object.Error{Message: fmt.Sprintf("cannot append to type: %s", args[0].Type())}
		},
	},
	"readFile": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			if arg, ok := args[0].(*object.String); ok {
				data, err := os.ReadFile(arg.Value)
				if err != nil {
					return &object.Error{Message: err.Error()}
				}
				return &object.String{Value: string(data)}
			}
			return &object.Error{Message: fmt.Sprintf("cannot read file from type: %s", args[0].Type())}
		},
	},
}
