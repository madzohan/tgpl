// first example in ch 1.5 adapted to be reusable and testable

package fetch1

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Response struct {
	Body io.Reader
	URL  string
}
type PageGetter func(url string) (resp *http.Response, err error)
type ResponsePrinter func(response *Response, outWrite io.Writer) (err error)
type URLsFetcher struct {
	urls          []string
	httpGet       PageGetter
	outWrite      io.Writer
	errWrite      io.Writer
	PrintResponse ResponsePrinter
}

func NewURLsFetcher(urls []string, httpGet PageGetter, outWrite io.Writer, errWrite io.Writer,
	PrintResponse ResponsePrinter) *URLsFetcher {
	f := URLsFetcher{urls, httpGet, outWrite, errWrite, PrintResponse}
	if PrintResponse == nil {
		f.PrintResponse = f._PrintResponse
	}
	return &f
}

func (f *URLsFetcher) _PrintResponse(response *Response, outWrite io.Writer) (err error) {
	b, err := io.ReadAll(ioutil.NopCloser(response.Body))
	fmt.Fprintf(outWrite, "%s\n", b)
	return
}

func (f *URLsFetcher) Fetch() {
	for _, url := range f.urls {
		response, err := f.httpGet(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fetch: while getting url=\"%s\" error occurs: \"%v\"\n", url, err)
			return
		}
		err = f.PrintResponse(&Response{Body: response.Body, URL: url}, f.outWrite)
		if err != nil {
			fmt.Fprintf(f.errWrite, "Fetch: while reading response from url=\"%s\" error occurs: \"%v\"\n", url, err)
			return
		}
	}
}
