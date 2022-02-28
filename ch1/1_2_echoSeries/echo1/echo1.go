package echo1

import (
	"fmt"
	"os"
)

// Print prints os.Args to Stdout
// using for loop
func Print() {
	echoStr, sep := "", ""
	for i := 1; i < len(os.Args); i++ {
		echoStr += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(echoStr)
}
