// Dup2 prints the count and text of lines that appear more than once
//in the input.  It reads from stdin or from a list of named files.
package dup2

import (
	"bufio"
	"fmt"
	"github.com/spf13/afero"
	"io"
	"os"
)

func scanBuf(file io.Reader, stats map[string]int) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stats[scanner.Text()]++
	}
}

func FindDuplicateLines(FS afero.Fs, reader io.Reader) {
	stats := make(map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		scanBuf(reader, stats)
	} else {
		for _, arg := range files {
			file, err := FS.Open(arg)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "dup2: %v", err)
				continue
			}
			scanBuf(file, stats)
			_ = file.Close()
		}
	}
	for line, count := range stats {
		if count > 1 {
			fmt.Printf("line: %s, count: %x", line, count)
		}
	}
}
