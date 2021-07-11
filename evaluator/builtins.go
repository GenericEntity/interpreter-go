package evaluator

import "github.com/GenericEntity/interpreter-go/monkey/object"

func checkArgsLen(wantedLength int, args ...object.Object) *object.Error {
	if len(args) != wantedLength {
		return newError("wrong number of arguments. got=%d, want=%d", len(args), wantedLength)
	}
	return nil
}

func newTypeNotSupportedError(funcName string, arg object.Object) *object.Error {
	return newError("argument to `%s` not supported, got %s", funcName, arg.Type())
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
				return newTypeNotSupportedError("len", arg)
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
				return newTypeNotSupportedError("first", arg)
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
				return newTypeNotSupportedError("last", arg)
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
				return newTypeNotSupportedError("rest", arg)
			}
		},
	},
}
