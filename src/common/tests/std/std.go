package std

import (
	"bytes"
	"io"
	"os"
	"testing"
)

var osArgs, osStdout = os.Args, os.Stdout

func SetUp(inArgs []string) (reader *os.File, writer *os.File) {
	reader, writer, _ = os.Pipe()
	os.Stdout = writer
	os.Args = inArgs
	return reader, writer
}

func TearDown(c testing.TB, reader *os.File, writer *os.File, expected string) {
	os.Args, os.Stdout = osArgs, osStdout
	// https://stackoverflow.com/a/10476304/3033586
	// copy the output in a separate goroutine so printing can't block indefinitely
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, reader)
		outC <- buf.String()
	}()
	_ = writer.Close()
	got := <-outC
	if expected != got {
		c.Errorf("expected in output = %s; got %s", expected, got)
	}
}
