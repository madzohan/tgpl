package tests

import (
	"testing"

	helloWorld "github.com/madzohan/tgpl/ch1/1_1_helloWorld"
)

func TestHelloWorldFunc(t *testing.T) {
	expO, expE := "Hello World!\n", ""
	t.Run("HelloWorld", func(t *testing.T) {
		or, ow, er, ew := SetUp([]string{})
		helloWorld.PrintHelloWorld()
		TearDown(t, or, ow, er, ew, expO, expE)
	})
}
