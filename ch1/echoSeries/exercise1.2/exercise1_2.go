package exercise1_2

import (
	"fmt"
	"os"
)

func Print() {
	for i, arg := range os.Args {
		fmt.Printf("index: %x, value: %s\n", i, arg)
	}
}
