// first example in ch 1.5 adapted to be used as an importable function and to be testable :)
// function Fetch
//	accepts 4 arguments:
//		- "urls" 			slice of strings (urls to be fetched)
//		- "httpGet" 		http client's GET function (for mockable in tests)
//		- "responseRead" 	http client response body read function (for mockable in tests)
//		- "outWriter"		io.Writer in practice default value is os.Stdout (for checking whole actual output in tests)
// 	returns nothing, prints text of entire response body from all urls to outWriter

package fetch1

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type PageGetter func(url string) (resp *http.Response, err error)
type ResponseReader func(r io.Reader) ([]byte, error)

func Fetch(urls []string, httpGet PageGetter, responseRead ResponseReader, outWriter io.Writer) {
	for _, url := range urls {
		response, err := httpGet(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fetch: while getting url=\"%s\" error occurs: \"%v\"\n", url, err)
			return
		}
		b, err := responseRead(ioutil.NopCloser(response.Body))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fetch: while reading response from url=\"%s\" error occurs: \"%v\"\n", url, err)
			return
		}
		fmt.Printf("%s\n", b)
	}
}
