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

func scan(scanner bufio.Scanner, filename string, stats map[string][]interface{}) {
	//map[line 1:[map[file 1:{} file 2:{} file 5:{}] map[0:453]]]
	for scanner.Scan() {
		line := scanner.Text()
		if item, exist := stats[line]; exist {
			lineCounter := item[1].(map[int]int)
			lineCounter[0]++
			if filename != "" {
				fileset := item[0].(map[string]struct{})
				if _, ok := fileset[filename]; !ok {
					fileset[filename] = struct{}{}
				}
			}
		} else {
			fileset := make(map[string]struct{})
			fileset[filename] = struct{}{}
			lineCounter := make(map[int]int)
			lineCounter[0]++
			stats[line] = []interface{}{fileset, lineCounter}
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

func printStats(stats map[string][]interface{}) {
	for line, item := range stats {
		lineCounter := item[1].(map[int]int)
		if lineCounter[0] > 1 {
			fileset := item[0].(map[string]struct{})
			filesetSorted := make([]string, len(fileset))
			i := 0
			for key := range fileset {
				filesetSorted[i] = key
				i++
			}
			sort.Strings(filesetSorted)
			filenames := strings.Join(filesetSorted, ", ")
			if filenames != "" {
				// exercise 1.4
				fmt.Printf("line: %s, count: %x, filenames: %s", line, lineCounter[0], filenames)
			} else {
				fmt.Printf("line: %s, count: %x", line, lineCounter[0])
			}

		}
	}
}

func FindDuplicateLines(FS afero.Fs, reader io.Reader) {
	stats := make(map[string][]interface{})
	files := os.Args[1:]
	if len(files) == 0 {
		scanner := getStdScanner(reader)
		scan(*scanner, "", stats)
	} else {
		for _, filename := range files {
			file, err := FS.Open(filename)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "dup2: %v", err)
				continue
			}
			scanner, filename := getFileScanner(file)
			scan(*scanner, filename, stats)
			_ = file.Close()
		}
	}
	printStats(stats)
}
