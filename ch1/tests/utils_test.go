package tests

import (
	"bytes"
	"io"
	"os"
	"testing"
)

var osArgs, osStdout = os.Args, os.Stdout

type fileBuf struct {
	writer *os.File
	reader *os.File
}

type chanSlice struct {
	filebuf        fileBuf
	expected       string
	expectedPrefix string
}

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
	chanMap := map[chan string]chanSlice{
		outChan: {fileBuf{outWriter, outReader}, expectedOut, "StdOut"},
		errChan: {fileBuf{errWriter, errReader}, expectedErr, "StdErr"},
	}
	for stdChan, i := range chanMap {
		// https://stackoverflow.com/a/10476304/3033586
		// copy the output in a separate goroutine so printing can't block indefinitely
		go func(reader *os.File, stdChan chan string) {
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, i.filebuf.reader)
			stdChan <- buf.String()
		}(i.filebuf.reader, stdChan)
		_ = i.filebuf.writer.Close()
		got := <-stdChan
		if i.expected != got {
			c.Errorf("%s expected %s; got %s", i.expectedPrefix, i.expected, got)
		}
	}
	os.Args, os.Stdout = osArgs, osStdout
}
