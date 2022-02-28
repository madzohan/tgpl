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

type OrderedMap struct {
	Lines *[]string
	M     map[string]*LineStats
}

func (orderedMap *OrderedMap) Set(line string, lineStats *LineStats) {
	_, present := orderedMap.M[line]
	orderedMap.M[line] = lineStats
	if !present {
		*orderedMap.Lines = append(*orderedMap.Lines, line)
	}
}

func (orderedMap *OrderedMap) Get(line string) (lineStats *LineStats, exist bool) {
	lineStats, exist = orderedMap.M[line]

	return
}

func scan(scanner bufio.Scanner, filename string, stats OrderedMap) {
	for scanner.Scan() {
		line := scanner.Text()
		if lineStats, exist := stats.Get(line); exist {
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
			stats.Set(line, &LineStats{fileset: fileset, count: 1})
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

func PrintStats(stats OrderedMap) {
	for _, line := range *stats.Lines {
		linestats, _ := stats.Get(line)
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

func PopulateLineStats(filenames []string, FS afero.Fs, reader io.Reader, stats OrderedMap) {
	if len(filenames) == 0 {
		scanner := getStdScanner(reader)
		scan(*scanner, "", stats)
	} else {
		for _, filename := range filenames {
			file, err := FS.Open(filename)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "dup: %v", err)

				continue
			}
			scanner, filename := getFileScanner(file)
			scan(*scanner, filename, stats)
			_ = file.Close()
		}
	}
}

// FindDuplicateLines prints the count and text of lines that appear more than once
// in the input, and filenames where they occurs (if input is from named files).
// Reads from stdin or from a list of named files.
func FindDuplicateLines(FS afero.Fs, reader io.Reader) {
	filenames := os.Args[1:]
	stats := OrderedMap{Lines: &[]string{}, M: make(map[string]*LineStats)}
	PopulateLineStats(filenames, FS, reader, stats)
	PrintStats(stats)
}
