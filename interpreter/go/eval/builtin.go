package eval

import (
	"fmt"
	"kaze/object"
	"os"
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
}
