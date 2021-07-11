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
}
