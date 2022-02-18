// Exercise 1.9
package fetch3

import (
	"fmt"
	"io"

	"github.com/madzohan/tgpl/ch1/1_56_fetchSeries/fetch1"
)

func PrintResponse(response *fetch1.Response, outWrite io.Writer) (err error) {
	_, err = io.Copy(outWrite, io.NopCloser(response.Body))
	fmt.Fprint(outWrite, "\n")
	fmt.Fprintf(outWrite, "HTTP status code=%v\n", response.StatusCode)
	return
}
