package lis

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
	"regexp"

	"github.com/spf13/afero"
)

var palette_defaults = []color.Color{color.Black, color.RGBA{0, 0xff, 0, 1}} // green on black
var colorIndex uint8 = 1

func Lissajous(out io.Writer, palette []color.Color, dynamicPalette bool) {
	if len(palette) != 2 && !dynamicPalette {
		palette = palette_defaults
	}
	const (
		cycles  = 5.0
		res     = 0.001
		size    = 500
		nframes = 64
		delay   = 8
	)

	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		if dynamicPalette {
			r, g, b := uint8(i), uint8(2*i), uint8(3*i)
			palette = []color.Color{color.Black, color.RGBA{r, g, b, 1}}
		}
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2.0*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			if dynamicPalette {
				r, g, b := uint8(3*t), uint8(t), uint8(2*t)
				colorIndex = uint8(color.Palette.Index(palette, color.RGBA{r, g, b, 1}))
			}
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), colorIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}

func LissajousToFile(filename string, palette []color.Color, dynamicPalette bool, fs afero.Fs) {
	f, err := fs.Create(filename)
	isAlphaNumerical := regexp.MustCompile(`^[\/\w.]+\.gif$`).MatchString
	if err != nil || !isAlphaNumerical(f.Name()) {
		fmt.Fprintf(os.Stderr, "LissajousToFile: invalid file\n")
		return
	}
	defer f.Close()
	Lissajous(f, palette, dynamicPalette)
	f.Sync()
}
