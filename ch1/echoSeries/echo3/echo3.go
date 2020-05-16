// Echo1 prints its command-line arguments. One line solution
package echo3

import (
	"fmt"
	"os"
	"strings"
)

func Print() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}
