package must

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

var errorNotAndOddInt = errors.New("not and odd int")

func getThisFile() string {
	_, file, _, _ := runtime.Caller(0)
	return file
}

var thisFile = getThisFile()

func incOdd(arg int) (result int, err error) {
	if arg%2 != 0 {
		return arg + 1, nil
	}
	return 0, errorNotAndOddInt
}

func incOdd2(arg int) (r1 int, r2 int, err error) {
	if arg%2 != 0 {
		return arg + 1, arg + 2, nil
	}
	return 0, 0, errorNotAndOddInt
}

func incOdd3(arg int) (r1 int, r2 int, r3 int, err error) {
	if arg%2 != 0 {
		return arg + 1, arg + 2, arg + 3, nil
	}
	return 0, 0, 0, errorNotAndOddInt
}

func incOdd4(arg int) (r1 int, r2 int, r3 int, r4 int, err error) {
	if arg%2 != 0 {
		return arg + 1, arg + 2, arg + 3, arg + 4, nil
	}
	return 0, 0, 0, 0, errorNotAndOddInt
}

func runIncR(t *testing.T) (err *Err) {
	defer func() {
		maybeErr := recover()
		if errMust, ok := maybeErr.(*Err); ok {
			err = errMust
			return
		}
	}()

	res := Do(incOdd(1)).R()
	assert.Equal(t, 2, res)
	Do2(incOdd2(res)).R()
	return nil
}

func runIncRf(t *testing.T) (err *Err) {
	defer func() {
		maybeErr := recover()
		if errMust, ok := maybeErr.(*Err); ok {
			err = errMust
			return
		}
	}()

	r1, r2, r3 := Do3(incOdd3(1)).Rf("incOdd%d", 3)
	assert.Equal(t, 2, r1)
	assert.Equal(t, 3, r2)
	assert.Equal(t, 4, r3)
	Do4(incOdd4(r3)).Rf("incOdd%d", 4)
	return
}

func runMust(t *testing.T) (reterr *Err) {
	defer func() {
		maybeErr := recover()
		if errMust, ok := maybeErr.(*Err); ok {
			reterr = errMust
			return
		}
	}()
	_, err := incOdd(0)
	Must(err)
	return
}

func runMustf(t *testing.T) (reterr *Err) {
	defer func() {
		maybeErr := recover()
		if errMust, ok := maybeErr.(*Err); ok {
			reterr = errMust
			return
		}
	}()

	_, err := incOdd(0)
	Mustf(err, "incOdd(%d)", 0)
	return
}

func runHold(t *testing.T) (reterr *Err) {
	defer func() {
		maybeErr := recover()
		if errMust, ok := maybeErr.(*Err); ok {
			reterr = errMust
			return
		}
	}()

	_, err := incOdd(0)
	Hold(err == nil)
	return
}

func runHoldf(t *testing.T) (reterr *Err) {
	defer func() {
		if errMust, ok := AsErrOrPanic(recover()); ok {
			reterr = errMust
			return
		}
	}()
	_, err := incOdd(0)
	Holdf(err == nil, "incOdd(%d)", 0)
	return
}

func runGet(t *testing.T) (reterr *Err) {
	defer func() {
		if errMust, ok := AsErrOrPanic(recover()); ok {
			reterr = errMust
			return
		}
	}()
	m := map[string]int{"1": 1}
	v := Get(m, "1")
	assert.Equal(t, 1, v)
	_ = Get(m, "key")
	return
}

func runGetf(t *testing.T) (reterr *Err) {
	defer func() {
		if errMust, ok := AsErrOrPanic(recover()); ok {
			reterr = errMust
			return
		}
	}()
	m := map[string]int{"1": 1}
	v := Getf(m, "1", "m[\"1\"]")
	assert.Equal(t, 1, v)
	_ = Getf(m, "2", "m[\"2\"]")
	return
}

func runCast(t *testing.T) (reterr *Err) {
	defer func() {
		if errMust, ok := AsErrOrPanic(recover()); ok {
			reterr = errMust
			return
		}
	}()
	s := "hello"
	var iface any
	iface = s
	sCasted := Cast[string](iface)
	assert.Equal(t, "hello", sCasted)
	_ = Cast[int](iface)
	return
}

func runCastf(t *testing.T) (reterr *Err) {
	defer func() {
		if errMust, ok := AsErrOrPanic(recover()); ok {
			reterr = errMust
			return
		}
	}()
	s := "hello"
	var iface any
	iface = s
	sCasted := Cast[string](iface)
	assert.Equal(t, "hello", sCasted)
	_ = Castf[int](iface, "iface.(int)")
	return
}

func TestMust(t *testing.T) {
	err := runIncR(t)
	assert.NotNil(t, err)
	assert.Equal(t, thisFile, err.File)
	assert.Equal(t, 59, err.Line)
	assert.Equal(t, "", err.Ctx)
	assert.Equal(t, errorNotAndOddInt, err.Err)
	assert.Equal(t, fmt.Sprintf("must() |%s:%d| failed with: not and odd int", err.File, err.Line), err.Error())

	err = runIncRf(t)
	assert.NotNil(t, err)
	assert.Equal(t, thisFile, err.File)
	assert.Equal(t, 76, err.Line)
	assert.Equal(t, "incOdd4", err.Ctx)
	assert.Equal(t, errorNotAndOddInt, err.Err)
	assert.Equal(t, fmt.Sprintf("must(incOdd4) |%s:%d| failed with: not and odd int", err.File, err.Line), err.Error())

	err = runMust(t)
	assert.NotNil(t, err)
	assert.Equal(t, thisFile, err.File)
	assert.Equal(t, 89, err.Line)
	assert.Equal(t, "", err.Ctx)
	assert.Equal(t, errorNotAndOddInt, err.Err)
	assert.Equal(t, fmt.Sprintf("must() |%s:%d| failed with: not and odd int", err.File, err.Line), err.Error())

	err = runMustf(t)
	assert.NotNil(t, err)
	assert.Equal(t, thisFile, err.File)
	assert.Equal(t, 103, err.Line)
	assert.Equal(t, "incOdd(0)", err.Ctx)
	assert.Equal(t, errorNotAndOddInt, err.Err)
	assert.Equal(t, fmt.Sprintf("must(incOdd(0)) |%s:%d| failed with: not and odd int", err.File, err.Line), err.Error())

	err = runHold(t)
	assert.NotNil(t, err)
	assert.Equal(t, thisFile, err.File)
	assert.Equal(t, 117, err.Line)
	assert.Equal(t, "", err.Ctx)
	assert.Equal(t, ErrHold, err.Err)
	assert.Equal(t, fmt.Sprintf("must() |%s:%d| failed with: condition did not hold true", err.File, err.Line), err.Error())

	err = runHoldf(t)
	assert.NotNil(t, err)
	assert.Equal(t, thisFile, err.File)
	assert.Equal(t, 129, err.Line)
	assert.Equal(t, "incOdd(0)", err.Ctx)
	assert.Equal(t, ErrHold, err.Err)
	assert.Equal(t, fmt.Sprintf("must(incOdd(0)) |%s:%d| failed with: condition did not hold true", err.File, err.Line), err.Error())

	err = runGet(t)
	assert.NotNil(t, err)
	assert.Equal(t, thisFile, err.File)
	assert.Equal(t, 143, err.Line)
	assert.Equal(t, "", err.Ctx)
	assert.Equal(t, ErrGet, err.Err)
	assert.Equal(t, fmt.Sprintf("must() |%s:%d| failed with: value is not in a map", err.File, err.Line), err.Error())

	err = runGetf(t)
	assert.NotNil(t, err)
	assert.Equal(t, thisFile, err.File)
	assert.Equal(t, 157, err.Line)
	assert.Equal(t, "m[\"2\"]", err.Ctx)
	assert.Equal(t, ErrGet, err.Err)
	assert.Equal(t, fmt.Sprintf("must(m[\"2\"]) |%s:%d| failed with: value is not in a map", err.File, err.Line), err.Error())

	err = runCast(t)
	assert.NotNil(t, err)
	assert.Equal(t, thisFile, err.File)
	assert.Equal(t, 173, err.Line)
	assert.Equal(t, "", err.Ctx)
	assert.Equal(t, ErrCast, err.Err)
	assert.Equal(t, fmt.Sprintf("must() |%s:%d| failed with: type assertion failed", err.File, err.Line), err.Error())

	err = runCastf(t)
	assert.NotNil(t, err)
	assert.Equal(t, thisFile, err.File)
	assert.Equal(t, 189, err.Line)
	assert.Equal(t, "iface.(int)", err.Ctx)
	assert.Equal(t, ErrCast, err.Err)
	assert.Equal(t, fmt.Sprintf("must(iface.(int)) |%s:%d| failed with: type assertion failed", err.File, err.Line), err.Error())
}
