package object

import (
	"fmt"
)

var kernelModule = newModule("Kernel", kernelMethodSet)
var kernelFunctions = NewEnclosedEnvironment(classes)

func init() {
	classes.Set("Kernel", kernelModule)
	kernelFunctions.Set("puts", &Builtin{
		Fn: func(args ...RubyObject) RubyObject {
			out := ""
			for _, arg := range args {
				out += arg.Inspect()
			}
			fmt.Println(out)
			return NIL
		},
	},
	)
}

var kernelMethodSet = map[string]RubyMethod{
	"nil?":    withArity(0, publicMethod(kernelIsNil)),
	"methods": withArity(0, publicMethod(kernelMethods)),
	"class":   withArity(0, publicMethod(kernelClass)),
	"puts":    privateMethod(kernelPuts),
	"require": withArity(1, privateMethod(kernelRequire)),
}

func kernelPuts(context CallContext, args ...RubyObject) (RubyObject, error) {
	out := ""
	for _, arg := range args {
		out += arg.Inspect()
	}
	fmt.Println(out)
	return NIL, nil
}

func kernelMethods(context CallContext, args ...RubyObject) (RubyObject, error) {
	var methodSymbols []RubyObject
	class := context.Receiver().Class()
	for class != nil {
		methods := class.Methods()
		for meth, fn := range methods {
			if fn.Visibility() == PUBLIC_METHOD {
				methodSymbols = append(methodSymbols, &Symbol{meth})
			}
		}
		class = class.SuperClass()
	}

	return &Array{Elements: methodSymbols}, nil
}

func kernelIsNil(context CallContext, args ...RubyObject) (RubyObject, error) {
	return FALSE, nil
}

func kernelClass(context CallContext, args ...RubyObject) (RubyObject, error) {
	class := context.Receiver().Class()
	if eigenClass, ok := class.(*eigenclass); ok {
		class = eigenClass.Class()
	}
	classObj := class.(RubyClassObject)
	return classObj, nil
}

func kernelRequire(context CallContext, args ...RubyObject) (RubyObject, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError(1, len(args))
	}
	name, ok := args[0].(*String)
	if !ok {
		return nil, NewImplicitConversionTypeError(name, args[0])
	}
	return &RequireStatement{name}, nil
}
