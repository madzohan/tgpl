package tests

import (
	"os"
	"strings"
	"testing"
	"tgpl/ch1/dupSeries/dup1"
	"tgpl/ch1/dupSeries/dup2"
	"tgpl/ch1/dupSeries/dup3"

	"github.com/spf13/afero"
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
		TearDown(t, or, ow, er, ew,
			"line: same_line, count: 3, filenames: test1, test2", "")
	})

	t.Run("dup2-files-stderr", func(t *testing.T) {
		or, ow, er, ew := SetUp([]string{"dup2-files", "file_not_exist"})
		dup2.FindDuplicateLines(FS, os.Stdin)
		TearDown(t, or, ow, er, ew, "",
			"dup2: open file_not_exist: file does not exist")
	})

	stdinBuf = strings.NewReader(stdin)
	t.Run("dup3-stdin", func(t *testing.T) {
		or, ow, er, ew := SetUp([]string{"dup3-stdin"})
		dup3.FindDuplicateLines(FS, stdinBuf)
		TearDown(t, or, ow, er, ew, expO, "")
	})

	t.Run("dup3-files", func(t *testing.T) {
		or, ow, er, ew := SetUp([]string{"dup3-files", "test1", "test2"})
		dup3.FindDuplicateLines(FS, os.Stdin)
		TearDown(t, or, ow, er, ew, expO, "")
	})

	t.Run("dup3-files-stderr", func(t *testing.T) {
		or, ow, er, ew := SetUp([]string{"dup3-files", "file_not_exist"})
		dup3.FindDuplicateLines(FS, os.Stdin)
		TearDown(t, or, ow, er, ew, "",
			"dup3: open file_not_exist: file does not exist\n")
	})
}

func BenchmarkFindDuplicateLines(b *testing.B) {
	stdin := "same_line\nyet_another_line\nsame_line\nanother_line\nsame_line\n"
	expO := "line: same_line, count: 3"

	b.Run("dup1", func(b *testing.B) {
		stdinBuf := strings.NewReader(stdin)
		or, ow, er, ew := SetUp([]string{})
		dup1.FindDuplicateLines(stdinBuf)
		TearDown(b, or, ow, er, ew, expO, "")
	})

	b.Run("dup2-stdin", func(b *testing.B) {
		stdinBuf := strings.NewReader(stdin)
		or, ow, er, ew := SetUp([]string{"dup2-stdin"})
		dup2.FindDuplicateLines(FS, stdinBuf)
		TearDown(b, or, ow, er, ew, expO, "")
	})

	b.Run("dup3-stdin", func(b *testing.B) {
		stdinBuf := strings.NewReader(stdin)
		or, ow, er, ew := SetUp([]string{"dup3-stdin"})
		dup3.FindDuplicateLines(FS, stdinBuf)
		TearDown(b, or, ow, er, ew, expO, "")
	})

	b.Run("dup2-files", func(b *testing.B) {
		or, ow, er, ew := SetUp([]string{"dup2-files", "test1", "test2"})
		dup2.FindDuplicateLines(FS, os.Stdin)
		TearDown(b, or, ow, er, ew,
			"line: same_line, count: 3, filenames: test1, test2", "")
	})

	b.Run("dup3-files", func(b *testing.B) {
		or, ow, er, ew := SetUp([]string{"dup3-files", "test1", "test2"})
		dup3.FindDuplicateLines(FS, os.Stdin)
		TearDown(b, or, ow, er, ew, expO, "")
	})
}
