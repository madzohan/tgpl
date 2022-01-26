// 	Dup2 prints the count and text of lines that appear more than once
// in the input, and filenames where they occurs (if input is from named files).
// It reads from stdin or from a list of named files.
package dup2

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/spf13/afero"
)

type LineStats struct {
	fileset interface{}
	count   int
}

func scan(scanner bufio.Scanner, filename string, stats map[string]*LineStats) {
	for scanner.Scan() {
		line := scanner.Text()
		if lineStats, exist := stats[line]; exist {
			lineStats.count++
			if filename != "" {
				fileset := lineStats.fileset.(map[string]struct{})
				if _, ok := fileset[filename]; !ok {
					fileset[filename] = struct{}{}
				}
			}
		} else {
			fileset := make(map[string]struct{})
			fileset[filename] = struct{}{}
			stats[line] = &LineStats{fileset: fileset, count: 1}
		}
	}
}

func getStdScanner(reader io.Reader) (scanner *bufio.Scanner) {
	scanner = bufio.NewScanner(reader)
	return scanner
}

func getFileScanner(reader afero.File) (scanner *bufio.Scanner, filename string) {
	filename = reader.Name()
	scanner = bufio.NewScanner(reader)
	return scanner, filename
}

func PrintStats(stats map[string]*LineStats) {
	for line, linestats := range stats {
		if linestats.count > 1 {
			fileset := linestats.fileset.(map[string]struct{})
			filenamesSorted := []string{}
			for key := range fileset {
				filenamesSorted = append(filenamesSorted, key)
			}
			sort.Strings(filenamesSorted)
			filenames := strings.Join(filenamesSorted, ", ")
			if filenames != "" {
				fmt.Printf("\nline: %s, count: %x, filenames: %s", line, linestats.count, filenames)
			} else {
				fmt.Printf("\nline: %s, count: %x", line, linestats.count)
			}

		}
	}
}

func PopulateLineStats(filenames []string, FS afero.Fs, reader io.Reader, lineStats map[string]*LineStats) {
	if len(filenames) == 0 {
		scanner := getStdScanner(reader)
		scan(*scanner, "", lineStats)
	} else {
		for _, filename := range filenames {
			file, err := FS.Open(filename)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "dup: %v", err)
				continue
			}
			scanner, filename := getFileScanner(file)
			scan(*scanner, filename, lineStats)
			_ = file.Close()
		}
	}
}

func FindDuplicateLines(FS afero.Fs, reader io.Reader) {
	filenames := os.Args[1:]
	lineStats := make(map[string]*LineStats)
	PopulateLineStats(filenames, FS, reader, lineStats)
	PrintStats(lineStats)
}
