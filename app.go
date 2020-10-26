package trier

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/shupkg/trier/log"

	"go.uber.org/dig"
)

type (
	Out = dig.Out
	In  = dig.In
)

var (
	app = dig.New()
)

var Exit = func() {
	os.Exit(1)
}

var HandleMustError = func(err error) {
	if err != nil {
		if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			log.Errorf("%v", err)
			Exit()
		}
	}
}

func Provide(constructors ...interface{}) error {
	for _, constructor := range constructors {
		if err := app.Provide(constructor); err != nil {
			return dig.RootCause(err)
		}
	}
	return nil
}

func ProvideMust(constructors ...interface{}) {
	HandleMustError(Provide(constructors...))
}

func Context(ctx context.Context) {
	_ = Provide(func() context.Context { return ctx })
}

func Supply(values ...interface{}) error {
	for _, v := range values {
		provide, err := newSupplyConstructor(v)
		if err != nil {
			return err
		}
		if err := app.Provide(provide); err != nil {
			return dig.RootCause(err)
		}
	}
	return nil
}

func SupplyMust(values ...interface{}) {
	HandleMustError(Supply(values...))
}

func Invoke(functions ...interface{}) error {
	for _, function := range functions {
		if err := app.Invoke(function); err != nil {
			return dig.RootCause(err)
		}
	}
	return nil
}

func InvokeMust(functions ...interface{}) {
	HandleMustError(Invoke(functions...))
}

func Populate(targets ...interface{}) error {
	function, err := newPopulateFunc(targets...)
	if err != nil {
		return err
	}

	if err := app.Invoke(function); err != nil {
		return dig.RootCause(err)
	}

	return nil
}

func PopulateMust(targets ...interface{}) {
	HandleMustError(Populate(targets...))
}

func newSupplyConstructor(value interface{}) (interface{}, error) {
	switch value.(type) {
	case nil:
		return nil, errors.New("untyped nil passed to fx.Supply")
	case error:
		return nil, errors.New("error value passed to fx.Supply")
	}

	returnTypes := []reflect.Type{reflect.TypeOf(value)}
	returnValues := []reflect.Value{reflect.ValueOf(value)}

	ft := reflect.FuncOf([]reflect.Type{}, returnTypes, false)
	fv := reflect.MakeFunc(ft, func([]reflect.Value) []reflect.Value {
		return returnValues
	})

	return fv.Interface(), nil
}

func newPopulateFunc(targets ...interface{}) (interface{}, error) {
	// Validate all targets are non-nil pointers.
	targetTypes := make([]reflect.Type, len(targets))
	for i, t := range targets {
		if t == nil {
			return nil, fmt.Errorf("failed to Populate: target %v is nil", i+1)
		}
		rt := reflect.TypeOf(t)
		if rt.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("failed to Populate: target %v is not a pointer type, got %T", i+1, t)
		}
		targetTypes[i] = reflect.TypeOf(t).Elem()
	}

	// Build a function that looks like:
	//
	// func(t1 T1, t2 T2, ...) {
	//   *targets[0] = t1
	//   *targets[1] = t2
	//   [...]
	// }
	//

	fnType := reflect.FuncOf(targetTypes, nil, false /* variadic */)
	fn := reflect.MakeFunc(fnType, func(args []reflect.Value) []reflect.Value {
		for i, arg := range args {
			reflect.ValueOf(targets[i]).Elem().Set(arg)
		}
		return nil
	})

	return fn.Interface(), nil
}

//func invokeErr(err error) interface{} {
//	return func() error {
//		return err
//	}
//}
