package tests

import (
	"ch1/dupSeries/dup1"
	"ch1/dupSeries/dup2"
	"github.com/spf13/afero"
	"os"
	"strings"
	"testing"
)

var (
	FS     afero.Fs
	FSUtil *afero.Afero
)

func init() {
	FS = afero.NewMemMapFs()
	FSUtil = &afero.Afero{Fs: FS}

	_ = FSUtil.WriteFile("test1", []byte("same_line\nyet_another_line\nsame_line\nanother_line\n"), 0644)
	_ = FSUtil.WriteFile("test2", []byte("same_line\n"), 0644)
}

func TestFindDuplicateLines(t *testing.T) {
	stdin := "same_line\nyet_another_line\nsame_line\nanother_line\nsame_line\n"
	expO := "line: same_line, count: 3"

	stdinBuf := strings.NewReader(stdin)
	t.Run("dup1", func(t *testing.T) {
		or, ow, er, ew := SetUp([]string{})
		dup1.FindDuplicateLines(stdinBuf)
		TearDown(t, or, ow, er, ew, expO, "")
	})

	stdinBuf = strings.NewReader(stdin)
	t.Run("dup2-stdin", func(t *testing.T) {
		or, ow, er, ew := SetUp([]string{"dup2-stdin"})
		dup2.FindDuplicateLines(FS, stdinBuf)
		TearDown(t, or, ow, er, ew, expO, "")
	})

	t.Run("dup2-files", func(t *testing.T) {
		or, ow, er, ew := SetUp([]string{"dup2-files", "test1", "test2"})
		dup2.FindDuplicateLines(FS, os.Stdin)
		TearDown(t, or, ow, er, ew, expO, "")
	})

	t.Run("dup2-files-stderr", func(t *testing.T) {
		or, ow, er, ew := SetUp([]string{"dup2-files", "file_not_exist"})
		dup2.FindDuplicateLines(FS, os.Stdin)
		TearDown(t, or, ow, er, ew, "",
			"dup2: open file_not_exist: file does not exist")
	})
}
