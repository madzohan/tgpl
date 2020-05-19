package tests

import (
	"ch1/echoSeries/echo1"
	"ch1/echoSeries/echo2"
	"ch1/echoSeries/echo3"
	"ch1/echoSeries/exercise1.1"
	"ch1/echoSeries/exercise1.2"
	"io/ioutil"
	"os"
	"testing"
)

var osArgs, osStdout = os.Args, os.Stdout

func setUp(inArgs []string) (reader *os.File, writer *os.File) {
	reader, writer, _ = os.Pipe()
	os.Stdout = writer
	os.Args = inArgs
	return reader, writer
}

func tearDown(c testing.TB, reader *os.File, writer *os.File, expected string) {
	os.Args, os.Stdout = osArgs, osStdout
	_ = writer.Close()
	out, _ := ioutil.ReadAll(reader)
	got := string(out)
	if expected != got {
		c.Errorf("expected in output = %s; got %s", expected, got)
	}
}

func TestPrint(t *testing.T) {
	inArgs := []string{"", "Cool,", "it", "works!"}
	expected := "Cool, it works!\n"
	t.Run("echo1", func(t *testing.T) { r, w := setUp(inArgs); echo1.Print(); tearDown(t, r, w, expected) })
	t.Run("echo2", func(t *testing.T) { r, w := setUp(inArgs); echo2.Print(); tearDown(t, r, w, expected) })
	t.Run("echo3", func(t *testing.T) { r, w := setUp(inArgs); echo3.Print(); tearDown(t, r, w, expected) })

	// exercise 1
	inArgs = []string{"exercise1", "Cool,", "it", "works!"}
	expected = "program: exercise1, args: Cool, it works!\n"
	t.Run("exercise1", func(t *testing.T) {
		r, w := setUp(inArgs)
		exercise1_1.Print()
		tearDown(t, r, w, expected)
	})

	// exercise 2
	inArgs = []string{"exercise2", "Cool,", "it", "works!"}
	expected = "index: 0, value: exercise2\n" +
		"index: 1, value: Cool,\n" +
		"index: 2, value: it\n" +
		"index: 3, value: works!\n"
	t.Run("exercise2", func(t *testing.T) {
		r, w := setUp(inArgs)
		exercise1_2.Print()
		tearDown(t, r, w, expected)
	})
}

func BenchmarkEcho(b *testing.B) {
	inArgs := []string{"", "Cool,", "it", "works!"}
	expected := "Cool, it works!\n"
	b.Run("echo1", func(b *testing.B) { r, w := setUp(inArgs); echo1.Print(); tearDown(b, r, w, expected) })
	b.Run("echo2", func(b *testing.B) { r, w := setUp(inArgs); echo2.Print(); tearDown(b, r, w, expected) })
	b.Run("echo3", func(b *testing.B) { r, w := setUp(inArgs); echo3.Print(); tearDown(b, r, w, expected) })
}
