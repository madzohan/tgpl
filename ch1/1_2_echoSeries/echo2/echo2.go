// Echo1 prints its command-line arguments.
package echo2

import (
	"fmt"
	"os"
)

const sep = " "

// Print prints os.Args to Stdout
// using for range switch
func Print() {
	var echoStr string

	for i, arg := range os.Args {
		switch i {
		case 0:
			continue
		case 1:
			echoStr = arg
		default:
			echoStr += sep + arg
		}
	}
	fmt.Println(echoStr)
}
