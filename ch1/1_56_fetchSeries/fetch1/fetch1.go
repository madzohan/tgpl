// first example in ch 1.5 adapted to be reusable and testable

package fetch1

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	Body       io.Reader
	URL        string
	StatusCode int
}
type PageGetter func(url string) (resp *http.Response, err error)
type ResponsePrinter func(response *Response, outWrite io.Writer) (err error)
type URLMaker func(url string) (validURL string)
type URLsFetcher struct {
	urls          []string
	httpGet       PageGetter
	outWrite      io.Writer
	errWrite      io.Writer
	PrintResponse ResponsePrinter
	MakeURL       URLMaker
}

func NewURLsFetcher(urls []string, httpGet PageGetter, outWrite io.Writer, errWrite io.Writer,
	PrintResponse ResponsePrinter, MakeURL URLMaker) *URLsFetcher {
	fetcher := URLsFetcher{urls, httpGet, outWrite, errWrite, PrintResponse, MakeURL}
	if PrintResponse == nil {
		fetcher.PrintResponse = fetcher._PrintResponse
	}
	if MakeURL == nil {
		fetcher.MakeURL = fetcher._MakeURL
	}

	return &fetcher
}

func (f *URLsFetcher) _PrintResponse(response *Response, outWrite io.Writer) error {
	b, err := io.ReadAll(ioutil.NopCloser(response.Body))
	fmt.Fprintf(outWrite, "%s\n", b)

	return err
}

func (f *URLsFetcher) _MakeURL(url string) (validURL string) {
	// Exercise 1.8
	prefix := "http://"
	validURL = url
	if !strings.HasPrefix(url, prefix) {
		validURL = prefix + url
	}

	return validURL
}

func (f *URLsFetcher) Fetch() {
	for _, url := range f.urls {
		response, err := f.httpGet(f.MakeURL(url))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fetch: while getting url=\"%s\" error occurs: \"%v\"\n", url, err)

			return
		}
		err = f.PrintResponse(&Response{Body: response.Body, URL: url, StatusCode: response.StatusCode}, f.outWrite)
		if err != nil {
			fmt.Fprintf(f.errWrite, "Fetch: while reading response from url=\"%s\" error occurs: \"%v\"\n", url, err)

			return
		}
	}
}
