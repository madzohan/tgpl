package tests

import (
	"bytes"
	"ch1/dupSeries/dup1"
	"common/tests/std"
	"testing"
)

func TestFindDuplicateLines(t *testing.T) {
	stdin := bytes.NewBufferString("")
	stdin.WriteString("same_line\n")
	stdin.WriteString("yet_another_line\n")
	stdin.WriteString("same_line\n")
	stdin.WriteString("another_line\n")
	stdin.WriteString("same_line\n")
	expO := "same_line: 3"

	t.Run("dup1", func(t *testing.T) {
		or, ow, er, ew := std.SetUp([]string{})
		dup1.FindDuplicateLines(stdin)
		std.TearDown(t, or, ow, er, ew, expO, "")
	})
}
