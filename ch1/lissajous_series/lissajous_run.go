package main

import (
	"image/color"

	"github.com/madzohan/tgpl/ch1/lissajous_series/lis"
	"github.com/spf13/afero"
)

func main() {
	fs := afero.NewOsFs()
	lis.LissajousToFile("./img.gif", []color.Color{}, true, fs)
}
