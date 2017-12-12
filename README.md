# testleak

[![Build Status](https://travis-ci.org/ceshihao/testleak.svg?branch=master)](https://travis-ci.org/ceshihao/testleak)
[![Go Report Card](https://goreportcard.com/badge/github.com/ceshihao/testleak)](https://goreportcard.com/report/github.com/ceshihao/testleak)
[![GoDoc](https://godoc.org/github.com/ceshihao/testleak?status.svg)](https://godoc.org/github.com/ceshihao/testleak)
[![Coverage Status](https://coveralls.io/repos/github/ceshihao/testleak/badge.svg?branch=master)](https://coveralls.io/github/ceshihao/testleak?branch=master)


A small goroutine leak detection package using in testing.

This package utilizes [`runtime.Stack()`](https://golang.org/pkg/runtime/#Stack) to detect all goroutines, and compare the goroutine list before and after testing.

This package also defines a whitelist to ignore some expected gotoutines in testing.
Default whitelist built in package is 
```go
// defaultTestLeakWhiteList with default values
var defaultTestLeakWhiteList = []string{"testing.Main(", "runtime.goexit", "testing.(*T).Run"}
```

User can also customize `whitelist` to ignore any expected goroutines.

## Features
* SetTestLeakWhiteList
* AppendTestLeakWhiteList
* RestoreDefaultTestLeakWhiteList
* TestLeak
## Install
Run `go get github.com/ceshihao/testleak`

## Usage
```go
package main

import (
    "testing"

    "github.com/ceshihao/testleak"
)

func TestUseTestLeak(t *testing.T) {
	defer testleak.TestLeak(t)()
	// Testing your code here
}
```
More details, please refer to [demo](./testleak_test.go).
## Documentation
https://godoc.org/github.com/ceshihao/testleak
