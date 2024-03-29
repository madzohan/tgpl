package tests

import (
	"bytes"
	"io"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

var (
	FS                         afero.Fs
	FSUtil                     *afero.Afero
	osArgs, osStdout, osStderr = os.Args, os.Stdout, os.Stderr
)

func init() {
	FS = afero.NewMemMapFs()
	FSUtil = &afero.Afero{Fs: FS}

	_ = FSUtil.WriteFile("test1", []byte("same_line\nyet_another_line\nsame_line\nanother_line\n"), 0644)
	_ = FSUtil.WriteFile("test2", []byte("same_line\n"), 0644)

	_ = FSUtil.MkdirAll("testdir1/td2", 0644)
	_ = FSUtil.MkdirAll("testdir2/td2", 0644)
	_ = FSUtil.WriteFile("testdir1/f1", []byte("same_line\nyet_another_line\nsame_line\nanother_line\n"), 0644)
	_ = FSUtil.WriteFile("testdir1/td2/f1", []byte("same_line\nsame_line\nanother_line\n"), 0644)
	_ = FSUtil.WriteFile("testdir2/td2/f1", []byte("same_line\nsame_line\nanother_line\n"), 0644)
}

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

func TearDown(tb testing.TB, outReader *os.File, outWriter *os.File,
	errReader *os.File, errWriter *os.File, expectedOut string, expectedErr string) {
	tb.Helper()
	outChan := make(chan string)
	errChan := make(chan string)
	chanMap := map[chan string]chanSlice{
		outChan: {fileBuf{outWriter, outReader}, expectedOut, "StdOut"},
		errChan: {fileBuf{errWriter, errReader}, expectedErr, "StdErr"},
	}
	for stdChan, i := range chanMap {
		// https://stackoverflow.com/a/10476304/3033586
		// copy the output in a separate goroutine so printing can't block indefinitely
		go func(_ *os.File, stdChan chan string) {
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, i.filebuf.reader)
			stdChan <- buf.String()
		}(i.filebuf.reader, stdChan)
		i.filebuf.writer.Close()
		got := <-stdChan
		if len(got) != len(i.expected) {
			tb.Errorf("In \"%s\" expected \"%s\"; got \"%s\"", i.expectedPrefix, i.expected, got)

			break
		}
		// compare sorted values
		gotS := strings.Split(got, "\n")
		expectedS := strings.Split(i.expected, "\n")
		sort.Strings(gotS)
		sort.Strings(expectedS)
		for j, v := range expectedS {
			if v != gotS[j] {
				tb.Errorf("In \"%s\" expected \"%s\"; got \"%s\"", i.expectedPrefix, i.expected, got)

				break
			}
		}
	}
	os.Args, os.Stdout, os.Stderr = osArgs, osStdout, osStderr
}
