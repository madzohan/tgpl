package tests

import (
	"image/color"
	"os"
	"testing"

	"github.com/madzohan/tgpl/ch1/lissajous_series/lis"
)

func TestLissajous(t *testing.T) {
	t.Run("lis1", func(t *testing.T) {
		_, outWriter, _ := os.Pipe()
		lis.Lissajous(outWriter, []color.Color{}, true)
		lis.LissajousToFile("testdir1/img.gif", []color.Color{}, false, FSUtil)
	})
	t.Run("lis1-stderr", func(t *testing.T) {
		outReader, outWriter, errReader, errWriter := SetUp([]string{})
		lis.LissajousToFile("", []color.Color{}, false, FSUtil)
		TearDown(t, outReader, outWriter, errReader, errWriter, "",
			"LissajousToFile: invalid file\n")
	})
}
