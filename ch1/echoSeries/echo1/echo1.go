package echo1

import (
	"fmt"
	"os"
)

func Print() {
	s, sep := "", ""
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}
