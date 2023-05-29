/*
must provides a set of functions to get rid of a following way-to-common pattern:

	if err != nil {
		log(...)
		return ResultType{}, err
	}

Currently can wrap up to four returned values (not counting last error). Uses generics.

Example

	file1 := must.Do(os.Open("file1")).R()
	file2 := must.Do(os.Open("file2")).Rf("Open(%s)", "file2")

	cmd := exec.Command(...)
	must.Mustf(cmd.Start(), "cmd start")

	val, hasVal := someMap[key]
	must.Hold(hasVal)

	v := must.Get(someMap, key)
	v2 := must.Getf(someMap, key2, "someMap[%s]", key2)

	var iFace interface{}
	iFace := "hello"
	helloString := must.Castf[string](iFace, "iFace.(string)")

	func LotsOfEarlyReturns() (retval int, reterr err) {
	    defer func() {
			if mustErr, ok := must.AsErrOrPanic(recover()); ok {
				retval = 0
				reterr = mustErr.Err
				log("LotsOfEarlyReturns failed: %v", mustErr)
				return
			}
	    }()
		r1 := must.Do(mightFail1()).R()
		r2 := must.Do(mightFail2(r1)).Rf("mightFail2(%d)", r1)
		r3 := must.Do(mightFail3(r2)).R()
		return r3, nil
	}
*/
package must

import (
	"errors"
	"fmt"
	"runtime"
)

// Do captures 1 value and an error.
// Resulting proxy object has methods R() or Rf() to retrieve a result.
func Do[T any](result T, err error) *Holder1[T] {
	return &Holder1[T]{result, err}
}

// Do2 captures 2 values and an error. For an example see Overview
func Do2[T1, T2 any](r1 T1, r2 T2, err error) *Holder2[T1, T2] {
	return &Holder2[T1, T2]{r1, r2, err}
}

// Do3 captures 3 values and an error. For an example see Overview
func Do3[T1, T2, T3 any](r1 T1, r2 T2, r3 T3, err error) *Holder3[T1, T2, T3] {
	return &Holder3[T1, T2, T3]{r1, r2, r3, err}
}

// Do4 captures 4 values and an error. For an example see Overview
func Do4[T1, T2, T3, T4 any](r1 T1, r2 T2, r3 T3, r4 T4, err error) *Holder4[T1, T2, T3, T4] {
	return &Holder4[T1, T2, T3, T4]{r1, r2, r3, r4, err}
}

// Get accesses a map m with the key k and panics with Err.Err == ErrGet if the key is not present.
func Get[Key comparable, Val any](m map[Key]Val, k Key) Val {
	v, ok := m[k]
	if !ok {
		panic(newErr(ErrGet, ""))
	}
	return v
}


// Getf behaves as Get but allows to specify an arbitrary context string.
// ctxFormat and ctxArgs are supplied as is to fmt.Sprintf. This context
// is available in Err.Ctx
func Getf[Key comparable, Val any](m map[Key]Val, k Key, ctxFormat string, ctxArgs... any) Val {
	v, ok := m[k]
	if !ok {
		panic(newErr(ErrGet, ctxFormat, ctxArgs...))
	}
	return v
}

// Cast tries to make a type casting (assertion) as val.(To) and panics with Err.Err == ErrCast if fails.
func Cast[To any](val any) To {
	casted, ok := val.(To)
	if !ok {
		panic(newErr(ErrCast, ""))
	}
	return casted
}

// Castf behaves as Cast but allows to specify an arbitrary context string.
// ctxFormat and ctxArgs are supplied as is to fmt.Sprintf. This context
// is available in Err.Ctx
func Castf[To any](val any, ctxFormat string, ctxArgs... any) To {
	casted, ok := val.(To)
	if !ok {
		panic(newErr(ErrCast, ctxFormat, ctxArgs...))
	}
	return casted
}

// R checks if a holder (proxy object returned by Do) is holding an error and if not - returns 1 result.
// If the error is present - panics with a non-nil *Err object with Err.Ctx == ""
func (h *Holder1[T]) R() T {
	if h.err != nil {
		panic(newErr(h.err, ""))
	}
	return h.result
}

// Rf behaves as R but allows to specify an arbitrary context string.
// ctxFormat and ctxArgs are supplied as is to fmt.Sprintf. This context
// is available in Err.Ctx
func (h *Holder1[T]) Rf(ctxFormat string, ctxArgs ...any) T {
	if h.err != nil {
		panic(newErr(h.err, ctxFormat, ctxArgs...))
	}
	return h.result
}

// R checks if a holder (proxy object returned by Do2) is holding an error and if not - returns 2 results.
// If the error is present - panics with a non-nil *Err object with Err.Ctx == ""
func (h *Holder2[T1, T2]) R() (T1, T2) {
	if h.err != nil {
		panic(newErr(h.err, ""))
	}
	return h.r1, h.r2
}

// Rf behaves as R but allows to specify an arbitrary context string.
// ctxFormat and ctxArgs are supplied as is to fmt.Sprintf. This context
// is available in Err.Ctx
func (h *Holder2[T1, T2]) Rf(ctxFormat string, ctxArgs ...any) (T1, T2) {
	if h.err != nil {
		panic(newErr(h.err, ctxFormat, ctxArgs...))
	}
	return h.r1, h.r2
}

// R checks if a holder (proxy object returned by Do3) is holding an error and if not - returns 3 results.
// If the error is present - panics with a non-nil *Err object with Err.Ctx == ""
func (h *Holder3[T1, T2, T3]) R() (T1, T2, T3) {
	if h.err != nil {
		panic(newErr(h.err, ""))
	}
	return h.r1, h.r2, h.r3
}

// Rf behaves as R but allows to specify an arbitrary context string.
// ctxFormat and ctxArgs are supplied as is to fmt.Sprintf. This context
// is available in Err.Ctx
func (h *Holder3[T1, T2, T3]) Rf(ctxFormat string, ctxArgs ...any) (T1, T2, T3) {
	if h.err != nil {
		panic(newErr(h.err, ctxFormat, ctxArgs...))
	}
	return h.r1, h.r2, h.r3
}

// R checks if a holder (proxy object returned by Do4) is holding an error and if not - returns 4 results.
// If the error is present - panics with a non-nil *Err object with Err.Ctx == ""
func (h *Holder4[T1, T2, T3, T4]) R() (T1, T2, T3, T4) {
	if h.err != nil {
		panic(newErr(h.err, ""))
	}
	return h.r1, h.r2, h.r3, h.r4
}

// Rf behaves as R but allows to specify an arbitrary context string.
// ctxFormat and ctxArgs are supplied as is to fmt.Sprintf. This context
// is available in Err.Ctx
func (h *Holder4[T1, T2, T3, T4]) Rf(ctxFormat string, ctxArgs ...any) (T1, T2, T3, T4) {
	if h.err != nil {
		panic(newErr(h.err, ctxFormat, ctxArgs...))
	}
	return h.r1, h.r2, h.r3, h.r4
}

// ErrHold is an error for Err.Err. Used for panic in Hold/Holdf
var ErrHold = errors.New("condition did not hold true")


// ErrGet is an error for Err.Err. Used for panic in Get/Getf
var ErrGet = errors.New("value is not in a map")

var ErrCast = errors.New("type assertion failed")

// Err is an error that is used for panic (by pointer) if must functions catch an error.
type Err struct {
	// Err is a caller error if Do*/Must functions were used or an ErrHold if Hold/Holdf were used
	Err error
	// Ctx is a fmt.Sprintf(ctxFormat, ctxArgs...) if Rf/Mustf/Holdf were used. Otherwise "".
	Ctx string
	// File is a file where package function was used and caught an error as provided by a runtime package
	File string
	// Line is a source code line number where package function was used
	// and caught an error as provided by a runtime package
	Line int
}

// Error implements an error interface and returns the following string
//
//	"must({Err.Ctx}) |{Err.File}:{Err.Line}| failed with: {Err.Err}"
func (err *Err) Error() string {
	return fmt.Sprintf("must(%s) |%s:%d| failed with: %v", err.Ctx, err.File, err.Line, err.Err)
}

// Hold can be used to panic on a failed condition. panic() with *Err will be called, Err.Err == ErrHold.
//
// Example
//
//	val, hasVal := someMap[key]
//	must.Hold(hasVal)
func Hold(cond bool) {
	if !cond {
		panic(newErr(ErrHold, ""))
	}
}

// Holdf behaves as Hold but allows to specify an arbitrary context string.
// ctxFormat and ctxArgs are supplied as is to fmt.Sprintf. This context
// is available in Err.Ctx
func Holdf(cond bool, ctxFormat string, ctxArgs ...any) {
	if !cond {
		panic(newErr(ErrHold, ctxFormat, ctxArgs...))
	}
}

// Must can be used to panic on non-nil errors.
// Useful when caller's function does not return values aside from an error.
// panic() with *Err will be called, Err.Err == err
//
// Example
//
//	cmd := exec.Command(...)
//	must.Must(cmd.Start())
func Must(err error) {
	if err != nil {
		panic(newErr(err, ""))
	}
}

// Mustf behaves as Must but allows to specify an arbitrary context string.
// ctxFormat and ctxArgs are supplied as is to fmt.Sprintf. This context
// is available in Err.Ctx
func Mustf(err error, ctxFormat string, ctxArgs ...any) {
	if err != nil {
		panic(newErr(err, ctxFormat, ctxArgs...))
	}
}

// AsErrOrPanic is intended for use in defer closures with recover()
//
// Any recovered object that is not a *must.Err will be repaniced.
// nil object will be ignored with (nil, false) return values.
func AsErrOrPanic(maybeErr any) (*Err, bool) {
	if maybeErr == nil {
		return nil, false
	}
	if mustErr, ok := maybeErr.(*Err); ok {
		return mustErr, true
	}
	panic(maybeErr)
}

func newErr(err error, ctxFormat string, ctxArgs ...any) *Err {
	_, file, line, _ := runtime.Caller(2)
	return &Err{
		Err:  err,
		Ctx:  fmt.Sprintf(ctxFormat, ctxArgs...),
		File: file,
		Line: line,
	}
}

type Holder1[T any] struct {
	result T
	err    error
}

type Holder2[T1, T2 any] struct {
	r1  T1
	r2  T2
	err error
}

type Holder3[T1, T2, T3 any] struct {
	r1  T1
	r2  T2
	r3  T3
	err error
}

type Holder4[T1, T2, T3, T4 any] struct {
	r1  T1
	r2  T2
	r3  T3
	r4  T4
	err error
}
