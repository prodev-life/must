package must

import (
	"errors"
	"fmt"
	"runtime"
)

/*
Usage: see must_test.go
*/

func Do[T any](result T, err error) *oneHolder[T] {
	return &oneHolder[T]{result, err}
}
func Do2[T1, T2 any](r1 T1, r2 T2, err error) *twoHolder[T1, T2] {
	return &twoHolder[T1, T2]{r1, r2, err}
}
func Do3[T1, T2, T3 any](r1 T1, r2 T2, r3 T3, err error) *threeHolder[T1, T2, T3] {
	return &threeHolder[T1, T2, T3]{r1, r2, r3, err}
}
func Do4[T1, T2, T3, T4 any](r1 T1, r2 T2, r3 T3, r4 T4, err error) *fourHolder[T1, T2, T3, T4] {
	return &fourHolder[T1, T2, T3, T4]{r1, r2, r3, r4, err}
}

func newErr(err error, ctxFormat string, ctxArgs... any) *Err {
	_, file, line, _ := runtime.Caller(2)
	return &Err{
		Err: err,
		Ctx: fmt.Sprintf(ctxFormat, ctxArgs...),
		File: file,
		Line: line,
	}
}

func (h *oneHolder[T]) R() T {
	if h.err != nil {
		panic(newErr(h.err, ""))
	}
	return h.result
}

func (h *oneHolder[T]) Rf(ctxFormat string, ctxArgs... any) T {
	if h.err != nil {
		panic(newErr(h.err, ctxFormat, ctxArgs...))
	}
	return h.result
}

func (h *twoHolder[T1, T2]) R() (T1, T2) {
	if h.err != nil {
		panic(newErr(h.err, ""))
	}
	return h.r1, h.r2
}

func (h *twoHolder[T1, T2]) Rf(ctxFormat string, ctxArgs... any) (T1, T2) {
	if h.err != nil {
		panic(newErr(h.err, ctxFormat, ctxArgs...))
	}
	return h.r1, h.r2
}


func (h *threeHolder[T1, T2, T3]) R() (T1, T2, T3) {
	if h.err != nil {
		panic(newErr(h.err, ""))
	}
	return h.r1, h.r2, h.r3
}

func (h *threeHolder[T1, T2, T3]) Rf(ctxFormat string, ctxArgs... any) (T1, T2, T3) {
	if h.err != nil {
		panic(newErr(h.err, ctxFormat, ctxArgs...))
	}
	return h.r1, h.r2, h.r3
}

func (h *fourHolder[T1, T2, T3, T4]) R() (T1, T2, T3, T4) {
	if h.err != nil {
		panic(newErr(h.err, ""))
	}
	return h.r1, h.r2, h.r3, h.r4
}

func (h *fourHolder[T1, T2, T3, T4]) Rf(ctxFormat string, ctxArgs... any) (T1, T2, T3, T4) {
	if h.err != nil {
		panic(newErr(h.err, ctxFormat, ctxArgs...))
	}
	return h.r1, h.r2, h.r3, h.r4
}


var ErrHold = errors.New("condition did not hold true")

type Err struct {
	Err  error
	Ctx string
	File string
	Line int
}

func (err *Err) Error() string {
	return fmt.Sprintf("must(%s) |%s:%d| failed with: %v", err.Ctx, err.File, err.Line, err.Err)
}

func Hold(cond bool) {
	if !cond {
		panic(newErr(ErrHold, ""))
	}
}

func Holdf(cond bool, ctxFormat string, ctxArgs ...any) {
	if !cond {
		panic(newErr(ErrHold, ctxFormat, ctxArgs...))
	}
}

func Must(err error) {
	if err != nil {
		panic(newErr(err, ""))
	}
}

func Mustf(err error, ctxFormat string, ctxArgs ...any) {
	if err != nil {
		panic(newErr(err, ctxFormat, ctxArgs...))
	}
}


type oneHolder[T any] struct {
	result T
	err    error
}

type twoHolder[T1, T2 any] struct {
	r1  T1
	r2  T2
	err error
}

type threeHolder[T1, T2, T3 any] struct {
	r1  T1
	r2  T2
	r3  T3
	err error
}

type fourHolder[T1, T2, T3, T4 any] struct {
	r1  T1
	r2  T2
	r3  T3
	r4  T4
	err error
}
