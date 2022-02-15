package fetch

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Result struct {
	Text  string
	Error error
}

type PageGetter func(url string) (resp *http.Response, err error)
type Http struct {
	getResponse PageGetter
}

func NewHttp(getFunc PageGetter) *Http {
	return &Http{getResponse: getFunc}
}

func (h *Http) Fetch(urls []string) <-chan Result {
	var (
		resp *http.Response
		err  error
	)
	ch := make(chan Result)
	for _, url := range urls {
		resp, err = h.getResponse(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch(url=%v) GET: %s\n", url, err)
			go func() {
				ch <- Result{Text: "", Error: err}
			}()

			continue
		}
		respBody := ioutil.NopCloser(resp.Body)
		scanner := bufio.NewScanner(respBody)
		go func() {
			buf := bytes.NewBufferString("")
			for scanner.Scan() {
				buf.WriteString(scanner.Text() + "\n")
			}
			ch <- Result{Text: buf.String(), Error: err}
		}()
	}
	return ch
}
