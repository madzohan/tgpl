//	Prints the count and text of lines that appear more than once
// also it should display all files in which each duplicated line occures
// file search recursively through the directory

package exercise1_4

import (
	"fmt"
	"io"
	"os"

	"github.com/madzohan/tgpl/ch1/1_3_dupSeries/dup2"
	"github.com/spf13/afero"
)

var filePaths []string

func walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Fprintf(os.Stderr, "exercise1_4: %v\n", err)
		return err
	}
	if !info.IsDir() {
		filePaths = append(filePaths, path)
	}
	return nil
}

func FindDuplicateLines(FS afero.Fs, reader io.Reader) {
	stats := dup2.OrderedMap{Lines: &[]string{}, M: make(map[string]*dup2.LineStats)}
	for _, dirname := range os.Args[1:] {
		err := afero.Walk(FS, dirname, walk)
		if err != nil {
			return
		}
	}
	dup2.PopulateLineStats(filePaths, FS, reader, stats)
	dup2.PrintStats(stats)
}
