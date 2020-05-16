package tests

import (
	"book/tgpl/src/ch1/echoSeries/echo1"
	"book/tgpl/src/ch1/echoSeries/echo2"
	"book/tgpl/src/ch1/echoSeries/echo3"
	"io/ioutil"
	"os"
	"testing"
)

var osArgs, osStdout = os.Args, os.Stdout

func setUp() (reader *os.File, writer *os.File) {
	reader, writer, _ = os.Pipe()
	os.Stdout = writer
	os.Args = []string{"", "Cool,", "it", "works!"}
	return reader, writer
}

func tearDown(t *testing.T, reader *os.File, writer *os.File) {
	os.Args, os.Stdout = osArgs, osStdout
	expected := "Cool, it works!\n"
	_ = writer.Close()
	out, _ := ioutil.ReadAll(reader)
	got := string(out)
	if expected != got {
		t.Errorf("expected in output = %s; got %s", expected, got)
	}
}

func TestPrint(t *testing.T) {
	t.Run("echo1", func(t *testing.T) { r, w := setUp(); echo1.Print(); tearDown(t, r, w) })
	t.Run("echo2", func(t *testing.T) { r, w := setUp(); echo2.Print(); tearDown(t, r, w) })
	t.Run("echo3", func(t *testing.T) { r, w := setUp(); echo3.Print(); tearDown(t, r, w) })
}
