package evaluator

import (
	"fmt"

	"github.com/GenericEntity/interpreter-go/monkey/object"
)

func checkArgsLen(wantedLength int, args ...object.Object) *object.Error {
	if len(args) != wantedLength {
		return newError("wrong number of arguments. got=%d, want=%d", len(args), wantedLength)
	}
	return nil
}

func formatPosition(position int) string {
	var suffix string
	if position > 10 && position < 20 {
		// 11th, 12th, 13th, ..., 19th
		suffix = "th"
	} else if position%10 == 1 {
		// 1st, 21st, 31st, etc.
		suffix = "st"
	} else if position%10 == 2 {
		// 2nd, 22nd, etc.
		suffix = "nd"
	} else if position%10 == 3 {
		// 3rd, 23rd, etc.
		suffix = "rd"
	} else {
		// 0th, 4th, 120th, etc.
		suffix = "th"
	}

	return fmt.Sprintf("%d%s", position, suffix)
}

func newTypeNotSupportedError(funcName string, position int, arg object.Object) *object.Error {
	return newError("type of %s argument to `%s` not supported, got %s", formatPosition(position), funcName, arg.Type())
}

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if err := checkArgsLen(1, args...); err != nil {
				return err
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}

			default:
				return newTypeNotSupportedError("len", 1, arg)
			}
		},
	},

	"first": {
		Fn: func(args ...object.Object) object.Object {
			if err := checkArgsLen(1, args...); err != nil {
				return err
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return newError("`first` should not be called on empty array")
				}
				return arg.Elements[0]

			default:
				return newTypeNotSupportedError("first", 1, arg)
			}
		},
	},

	"last": {
		Fn: func(args ...object.Object) object.Object {
			if err := checkArgsLen(1, args...); err != nil {
				return err
			}

			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				if length == 0 {
					return newError("`last` should not be called on empty array")
				}
				return arg.Elements[length-1]

			default:
				return newTypeNotSupportedError("last", 1, arg)
			}
		},
	},

	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if err := checkArgsLen(1, args...); err != nil {
				return err
			}

			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				if length == 0 {
					return newError("`rest` should not be called on empty array")
				}

				elems := make([]object.Object, length-1)
				copy(elems, arg.Elements[1:length])
				return &object.Array{Elements: elems}

			default:
				return newTypeNotSupportedError("rest", 1, arg)
			}
		},
	},

	"push": {
		Fn: func(args ...object.Object) object.Object {
			if err := checkArgsLen(2, args...); err != nil {
				return err
			}

			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				elems := make([]object.Object, length+1)
				copy(elems, arg.Elements)
				elems[length] = args[1]
				return &object.Array{Elements: elems}

			default:
				return newTypeNotSupportedError("push", 1, arg)
			}
		},
	},

	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},

	"put": {
		Fn: func(args ...object.Object) object.Object {
			const funcName = "put"
			if err := checkArgsLen(3, args...); err != nil {
				return err
			}

			switch arg := args[0].(type) {
			case *object.Hash:
				// second arg must be Hashable
				h, ok := args[1].(object.Hashable)
				if !ok {
					return newTypeNotSupportedError(funcName, 2, args[1])
				}

				// copy map
				newPairs := make(map[object.HashKey]object.HashPair)
				for k, v := range arg.Pairs {
					newPairs[k] = v
				}

				// add (or replace) key-value pair
				newPairs[h.HashKey()] = object.HashPair{Key: args[1], Value: args[2]}

				return &object.Hash{Pairs: newPairs}

			default:
				return newTypeNotSupportedError(funcName, 1, arg)
			}
		},
	},
}
