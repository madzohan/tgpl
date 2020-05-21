package std

import (
	"bytes"
	"io"
	"os"
	"testing"
)

var osArgs, osStdout = os.Args, os.Stdout

func SetUp(inArgs []string) (outReader *os.File, outWriter *os.File, errReader *os.File, errWriter *os.File) {
	outReader, outWriter, _ = os.Pipe()
	errReader, errWriter, _ = os.Pipe()
	os.Stdout = outWriter
	os.Stderr = errWriter
	os.Args = inArgs
	return outReader, outWriter, errReader, errWriter
}

func TearDown(c testing.TB, outReader *os.File, outWriter *os.File,
	errReader *os.File, errWriter *os.File, expectedOut string, expectedErr string) {
	outChan := make(chan string)
	errChan := make(chan string)
	chanMap := map[chan string][]interface{}{
		outChan: {[]*os.File{outWriter, outReader}, expectedOut, "StdOut"},
		errChan: {[]*os.File{errWriter, errReader}, expectedErr, "StdErr"},
	}
	for stdChan, i := range chanMap {
		writer, reader := i[0].([]*os.File)[0], i[0].([]*os.File)[1]
		expected := i[1].(string)
		expectedPrefix := i[2].(string)
		// https://stackoverflow.com/a/10476304/3033586
		// copy the output in a separate goroutine so printing can't block indefinitely
		go func(reader *os.File, stdChan chan string) {
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, reader)
			stdChan <- buf.String()
		}(reader, stdChan)
		_ = writer.Close()
		got := <-stdChan
		if expected != got {
			c.Errorf("%s expected %s; got %s", expectedPrefix, expected, got)
		}
	}
	os.Args, os.Stdout = osArgs, osStdout
}
