package tests

import (
	"ch1/helloworld"
	"testing"
)

func TestHelloWorldFunc(t *testing.T) {
	expO, expE := "Hello World!\n", ""
	t.Run("HelloWorld", func(t *testing.T) {
		or, ow, er, ew := SetUp([]string{})
		helloworld.PrintHelloWorld()
		TearDown(t, or, ow, er, ew, expO, expE)
	})
}
