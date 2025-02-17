package eval

import (
	"fmt"
	"kaze/object"
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
}
