package echo_exercise1

import (
	"fmt"
	"os"
	"strings"
)

func Print() {
	fmt.Printf("program: %s, args: %s\n", os.Args[0], strings.Join(os.Args[1:], " "))
}
