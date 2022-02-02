package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/madzohan/tgpl/ch1/1_3_dupSeries/dup1"
	"github.com/madzohan/tgpl/ch1/1_3_dupSeries/dup2"
	"github.com/madzohan/tgpl/ch1/1_3_dupSeries/dup3"
	"github.com/madzohan/tgpl/ch1/1_3_dupSeries/exercise1_4"
)

func TestFindDuplicateLines(t *testing.T) {
	stdin := "same_line\nyet_another_line\nsame_line\nanother_line\nsame_line\n"

	stdinBuf := strings.NewReader(stdin)
	t.Run("dup1", func(t *testing.T) {
		outReader, outWriter, errReader, errWriter := SetUp([]string{"dup1"})
		dup1.FindDuplicateLines(stdinBuf)
		TearDown(t, outReader, outWriter, errReader, errWriter, "line: same_line, count: 3", "")
	})

	stdinBuf = strings.NewReader(stdin)
	t.Run("dup2-stdin", func(t *testing.T) {
		outReader, outWriter, errReader, errWriter := SetUp([]string{"dup2-stdin"})
		dup2.FindDuplicateLines(FS, stdinBuf)
		TearDown(t, outReader, outWriter, errReader, errWriter, "\nline: same_line, count: 3", "")
	})

	t.Run("dup2-files", func(t *testing.T) {
		outReader, outWriter, errReader, errWriter := SetUp([]string{"dup2-files", "test1", "test2"})
		dup2.FindDuplicateLines(FS, os.Stdin)
		TearDown(t, outReader, outWriter, errReader, errWriter,
			"\nline: same_line, count: 3, filenames: test1, test2", "")
	})

	t.Run("dup2-files-stderr", func(t *testing.T) {
		outReader, outWriter, errReader, errWriter := SetUp([]string{"dup2-files-stderr", "file_not_exist"})
		dup2.FindDuplicateLines(FS, os.Stdin)
		TearDown(t, outReader, outWriter, errReader, errWriter, "",
			"dup: open file_not_exist: file does not exist")
	})

	stdinBuf = strings.NewReader(stdin)
	t.Run("dup3-stdin", func(t *testing.T) {
		outReader, outWriter, errReader, errWriter := SetUp([]string{"dup3-stdin"})
		dup3.FindDuplicateLines(FS, stdinBuf)
		TearDown(t, outReader, outWriter, errReader, errWriter, "line: same_line, count: 3", "")
	})

	t.Run("dup3-files", func(t *testing.T) {
		outReader, outWriter, errReader, errWriter := SetUp([]string{"dup3-files", "test1", "test2"})
		dup3.FindDuplicateLines(FS, os.Stdin)
		TearDown(t, outReader, outWriter, errReader, errWriter, "line: same_line, count: 3", "")
	})

	t.Run("dup3-files-stderr", func(t *testing.T) {
		outReader, outWriter, errReader, errWriter := SetUp([]string{"dup3-files-stderr", "file_not_exist"})
		dup3.FindDuplicateLines(FS, os.Stdin)
		TearDown(t, outReader, outWriter, errReader, errWriter, "",
			"dup3: open file_not_exist: file does not exist\n")
	})

	t.Run("exercise1_4-files", func(t *testing.T) {
		outReader, outWriter, errReader, errWriter := SetUp([]string{"exercise1_4-files", "testdir1", "testdir2"})
		exercise1_4.FindDuplicateLines(FS, os.Stdin)
		TearDown(t, outReader, outWriter, errReader, errWriter,
			"\nline: same_line, count: 6, filenames: testdir1/f1, testdir1/td2/f1, testdir2/td2/f1"+
				"\nline: another_line, count: 3, filenames: testdir1/f1, testdir1/td2/f1, testdir2/td2/f1", "")
	})

	t.Run("exercise1_4-files-stderr", func(t *testing.T) {
		outReader, outWriter, errReader, errWriter := SetUp([]string{"exercise1_4-files", "dir_not_exist"})
		exercise1_4.FindDuplicateLines(FS, os.Stdin)
		TearDown(t, outReader, outWriter, errReader, errWriter, "",
			"exercise1_4: open dir_not_exist: file does not exist\n")
	})
}

func BenchmarkFindDuplicateLines(b *testing.B) {
	stdin := "same_line\nyet_another_line\nsame_line\nanother_line\nsame_line\n"
	expectedOut := "line: same_line, count: 3"

	b.Run("dup1", func(b *testing.B) {
		stdinBuf := strings.NewReader(stdin)
		outReader, outWriter, errReader, errWriter := SetUp([]string{})
		dup1.FindDuplicateLines(stdinBuf)
		TearDown(b, outReader, outWriter, errReader, errWriter, expectedOut, "")
	})

	b.Run("dup2-stdin", func(b *testing.B) {
		stdinBuf := strings.NewReader(stdin)
		outReader, outWriter, errReader, errWriter := SetUp([]string{"dup2-stdin"})
		dup2.FindDuplicateLines(FS, stdinBuf)
		TearDown(b, outReader, outWriter, errReader, errWriter, expectedOut, "")
	})

	b.Run("dup3-stdin", func(b *testing.B) {
		stdinBuf := strings.NewReader(stdin)
		outReader, outWriter, errReader, errWriter := SetUp([]string{"dup3-stdin"})
		dup3.FindDuplicateLines(FS, stdinBuf)
		TearDown(b, outReader, outWriter, errReader, errWriter, expectedOut, "")
	})

	b.Run("dup2-files", func(b *testing.B) {
		outReader, outWriter, errReader, errWriter := SetUp([]string{"dup2-files", "test1", "test2"})
		dup2.FindDuplicateLines(FS, os.Stdin)
		TearDown(b, outReader, outWriter, errReader, errWriter,
			"line: same_line, count: 3, filenames: test1, test2", "")
	})

	b.Run("dup3-files", func(b *testing.B) {
		outReader, outWriter, errReader, errWriter := SetUp([]string{"dup3-files", "test1", "test2"})
		dup3.FindDuplicateLines(FS, os.Stdin)
		TearDown(b, outReader, outWriter, errReader, errWriter, expectedOut, "")
	})
}
