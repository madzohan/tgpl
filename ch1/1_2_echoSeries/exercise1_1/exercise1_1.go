package exercise1_1

import (
	"fmt"
	"os"
	"strings"
)

// Print prints os.Args to Stdout
// using strings.Join
func Print() {
	fmt.Printf("program: %s, args: %s\n", os.Args[0], strings.Join(os.Args[1:], " "))
}
