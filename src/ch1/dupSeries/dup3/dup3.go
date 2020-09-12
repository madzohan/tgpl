//	Prints the count and text of lines that appear more than once
//	Yet another approach of reading file - read entire input in memory
// using afero.ReadFile
package dup3

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/afero"
)

func scanBuf(data []byte, err error, stats map[string]int) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		if line == "" {
			continue
		}
		stats[line]++
	}
}

func FindDuplicateLines(FS afero.Fs, reader io.Reader) {
	stats := make(map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		data, err := afero.ReadAll(reader)
		scanBuf(data, err, stats)
	} else {
		for _, filename := range files {
			data, err := afero.ReadFile(FS, filename)
			scanBuf(data, err, stats)
		}
	}
	for line, count := range stats {
		if count > 1 {
			fmt.Printf("line: %s, count: %x", line, count)
		}
	}
}
