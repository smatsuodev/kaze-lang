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
				_args = append(_args, arg.Inspect())
			}
			fmt.Print(strings.Join(_args, " "))
			return NULL
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			var _args []string
			for _, arg := range args {
				_args = append(_args, arg.Inspect())
			}
			fmt.Println(strings.Join(_args, " "))
			return NULL
		},
	},
	"args": {
		Fn: func(_ ...object.Object) object.Object {
			var args []object.Object
			for _, arg := range os.Args {
				args = append(args, &object.String{Value: arg})
			}
			return &object.Array{Elements: args}
		},
	},
}
