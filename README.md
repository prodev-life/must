must
====

Overview
--------

`must` provides a set of functions to get rid of a following way-to-common go pattern:
```go
if err != nil {
	log(...)
	return ResultType{}, err
}
```

Currently it can wrap up to four returned values (not counting one last error). Uses generics.

### Example
```go
file1 := must.Do(os.Open("file1")).R()
file2 := must.Do(os.Open("file2")).Rf("Open(%s)", "file2")

cmd := exec.Command(...)
must.Mustf(cmd.Start(), "cmd start")

val, hasVal := someMap[key]
must.Hold(hasVal)

func CouldBeLotsOfEarlyReturns() (retval int, reterr err) {
    defer func() {
        if mustErr, ok := must.AsErrOrPanic(recover()); ok {
			retval = 0
			reterr = mustErr.Err
			log("CouldBeLotsOfEarlyReturns failed: %v", mustErr)
		}
    }()
	r1 := must.Do(mightFail1()).R()
	r2 := must.Do(mightFail2(r1)).Rf("mightFail2(%d)", r1)
	r3 := must.Do(mightFail3(r2)).R()
	return r3, nil
}
```

Documentation
-------------
TBD with a link to pkg.go.dev