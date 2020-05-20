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
	expected := "same_line: 3"

	t.Run("echo1", func(t *testing.T) {
		r, w := std.SetUp([]string{})
		dup1.FindDuplicateLines(stdin)
		std.TearDown(t, r, w, expected)
	})
}
