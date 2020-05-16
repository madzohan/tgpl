// Echo1 prints its command-line arguments.
package echo2

import (
	"fmt"
	"os"
)

const sep = " "

func Print() {
	var s string
	for i, arg := range os.Args {
		switch i {
		case 0:
			continue
		case 1:
			s = arg
		default:
			s += sep + arg
		}
	}
	fmt.Println(s)
}
