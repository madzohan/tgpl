// Echo1 prints its command-line arguments. One line solution
package echo3

import (
	"fmt"
	"os"
	"strings"
)

// Print prints os.Args to Stdout
// using strings.Join
func Print() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}
