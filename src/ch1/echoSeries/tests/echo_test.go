package tests

import (
	"ch1/echoSeries/echo1"
	"ch1/echoSeries/echo2"
	"ch1/echoSeries/echo3"
	"ch1/echoSeries/exercise1.1"
	"ch1/echoSeries/exercise1.2"
	"common/tests/std"
	"testing"
)

func TestPrint(t *testing.T) {
	inArgs := []string{"", "Cool,", "it", "works!"}
	expO, expE := "Cool, it works!\n", ""
	t.Run("echo1", func(t *testing.T) {
		or, ow, er, ew := std.SetUp(inArgs)
		echo1.Print()
		std.TearDown(t, or, ow, er, ew, expO, expE)
	})
	t.Run("echo2", func(t *testing.T) {
		or, ow, er, ew := std.SetUp(inArgs)
		echo2.Print()
		std.TearDown(t, or, ow, er, ew, expO, expE)
	})
	t.Run("echo3", func(t *testing.T) {
		or, ow, er, ew := std.SetUp(inArgs)
		echo3.Print()
		std.TearDown(t, or, ow, er, ew, expO, expE)
	})

	// exercise 1
	inArgs = []string{"exercise1", "Cool,", "it", "works!"}
	expO, expE = "program: exercise1, args: Cool, it works!\n", ""
	t.Run("exercise1", func(t *testing.T) {
		or, ow, er, ew := std.SetUp(inArgs)
		exercise1_1.Print()
		std.TearDown(t, or, ow, er, ew, expO, expE)
	})

	// exercise 2
	inArgs = []string{"exercise2", "Cool,", "it", "works!"}
	expO, expE = "index: 0, value: exercise2\n"+
		"index: 1, value: Cool,\n"+
		"index: 2, value: it\n"+
		"index: 3, value: works!\n", ""
	t.Run("exercise2", func(t *testing.T) {
		or, ow, er, ew := std.SetUp(inArgs)
		exercise1_2.Print()
		std.TearDown(t, or, ow, er, ew, expO, expE)
	})
}

func BenchmarkEcho(b *testing.B) {
	inArgs := []string{"", "Cool,", "it", "works!"}
	expO, expE := "Cool, it works!\n", ""
	b.Run("echo1", func(b *testing.B) {
		or, ow, er, ew := std.SetUp(inArgs)
		echo1.Print()
		std.TearDown(b, or, ow, er, ew, expO, expE)
	})
	b.Run("echo2", func(b *testing.B) {
		or, ow, er, ew := std.SetUp(inArgs)
		echo2.Print()
		std.TearDown(b, or, ow, er, ew, expO, expE)
	})
	b.Run("echo3", func(b *testing.B) {
		or, ow, er, ew := std.SetUp(inArgs)
		echo3.Print()
		std.TearDown(b, or, ow, er, ew, expO, expE)
	})
}
