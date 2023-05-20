package must

import (
	"errors"
	"fmt"
)

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

func P[T any](result T, err error) *oneHolder[T] {
	return &oneHolder[T]{result, err}
}
func P2[T1, T2 any](r1 T1, r2 T2, err error) *twoHolder[T1, T2] {
	return &twoHolder[T1, T2]{r1, r2, err}
}
func P3[T1, T2, T3 any](r1 T1, r2 T2, r3 T3, err error) *threeHolder[T1, T2, T3] {
	return &threeHolder[T1, T2, T3]{r1, r2, r3, err}
}
func P4[T1, T2, T3, T4 any](r1 T1, r2 T2, r3 T3, r4 T4, err error) *fourHolder[T1, T2, T3, T4] {
	return &fourHolder[T1, T2, T3, T4]{r1, r2, r3, r4, err}
}

var ErrHold = errors.New("condition did not hold true")

type ErrMust struct {
	Err  error
	Name string
}

func (err *ErrMust) Error() string {
	return fmt.Sprintf("must.*(%s) failed with: %v", err.Name, err.Err)
}

func Do[T any](holder *oneHolder[T], format string, args ...any) T {
	if holder.err != nil {
		panic(&ErrMust{Err: holder.err, Name: fmt.Sprintf(format, args...)})
	}
	return holder.result
}

func Do2[T1, T2 any](holder *twoHolder[T1, T2], format string, args ...any) (T1, T2) {
	if holder.err != nil {
		panic(&ErrMust{Err: holder.err, Name: fmt.Sprintf(format, args...)})
	}
	return holder.r1, holder.r2
}

func Do3[T1, T2, T3 any](holder *threeHolder[T1, T2, T3], format string, args ...any) (T1, T2, T3) {
	if holder.err != nil {
		panic(&ErrMust{Err: holder.err, Name: fmt.Sprintf(format, args...)})
	}
	return holder.r1, holder.r2, holder.r3
}

func Do4[T1, T2, T3, T4 any](holder *fourHolder[T1, T2, T3, T4], format string, args ...any) (T1, T2, T3, T4) {
	if holder.err != nil {
		panic(&ErrMust{Err: holder.err, Name: fmt.Sprintf(format, args...)})
	}
	return holder.r1, holder.r2, holder.r3, holder.r4
}

func Hold(cond bool, format string, args ...any) {
	if !cond {
		panic(&ErrMust{Err: ErrHold, Name: fmt.Sprintf(format, args...)})
	}
}

func Must(err error, format string, args ...any) {
	if err != nil {
		panic(&ErrMust{Err: err, Name: fmt.Sprintf(format, args...)})
	}
}

func Work(int) (int, error) { return 0, nil }

func Test() (int, error) {
	v1, err := Work(0)
	if err != nil {
		return 0, err
	}
	v2, err := Work(v1)
	if err != nil {
		return 0, err
	}
	v3, err := Work(v2)
	if err != nil {
		return 0, err
	}
	return v3, nil
}

func Test2() (v3 int, err error) {
	defer func() {
		maybeErr := recover()
		if errMust, ok := maybeErr.(*ErrMust); ok {
			v3 = 0
			err = errMust.Err
			return
		}
		panic(maybeErr)
	}()
	v1 := Do(P(Work(0)), "Work(0)")
	v2 := Do(P(Work(v1)), "Work(v1)")
	v3 = Do(P(Work(v2)), "Work(v2)")
	return
}
