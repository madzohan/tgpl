package exercise1_2

import (
	"fmt"
	"os"
)

// Print prints os.Args to Stdout
// using for range
func Print() {
	for i, arg := range os.Args {
		fmt.Printf("index: %x, value: %s\n", i, arg)
	}
}
