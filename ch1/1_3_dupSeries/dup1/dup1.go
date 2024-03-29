package dup1

import (
	"bufio"
	"fmt"
	"io"
)

// FindDuplicateLines prints the text of each line that appears more than
// once in the standard input, preceded by its count. using `map` and `bufio.NewScanner`
func FindDuplicateLines(reader io.Reader) {
	stats := make(map[string]int)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		stats[scanner.Text()]++
	}

	for line, count := range stats {
		if count > 1 {
			fmt.Printf("\nline: %s, count: %x", line, count)
		}
	}
}
