package fetchall

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type HttpIface interface {
	Get(url string) (resp *http.Response, err error)
	Size(responseBody io.ReadCloser) (numBytes int64, err error)
}
type TimeIface interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

type Result struct {
	Text  string
	Error error
}

type ResultHandler func(Result)

type URLsFetcher struct {
	Time TimeIface
	Http HttpIface

	handleResult ResultHandler

	outWrite io.Writer
	errWrite io.Writer
}

// TimeIface real (default) implementation
type RealTime struct{}

func (RealTime) Now() time.Time                  { return time.Now() }
func (RealTime) Since(s time.Time) time.Duration { return time.Since(s) }

// HttpIface real (default) implementation
type RealHttp struct{}

func (RealHttp) Get(url string) (resp *http.Response, err error) { return http.Get(url) }
func (RealHttp) Size(responseBody io.ReadCloser) (numBytes int64, err error) {
	numBytes, err = io.Copy(io.Discard, responseBody)
	return
}

var ( // default <nil> values
	DefaultTime     TimeIface
	DefaultHttp     HttpIface
	DefaultOutWrite *os.File
	DefaultErrWrite *os.File
)

func NewURLsFetcher(Time TimeIface, Http HttpIface, handleResult ResultHandler, outWrite io.Writer, errWrite io.Writer) *URLsFetcher {
	fetcher := URLsFetcher{Time, Http, handleResult, outWrite, errWrite}

	if Time == DefaultTime {
		fetcher.Time = RealTime{}
	}
	if Http == DefaultHttp {
		fetcher.Http = RealHttp{}
	}
	if outWrite == DefaultOutWrite {
		fetcher.outWrite = os.Stdout
	}
	if errWrite == DefaultErrWrite {
		fetcher.errWrite = os.Stderr
	}
	if handleResult == nil {
		fetcher.handleResult = func(res Result) {
			if res.Error != nil {
				fmt.Fprintln(fetcher.errWrite, res.Error)
			} else {
				fmt.Fprintln(fetcher.outWrite, res.Text)
			}
		}
	}

	return &fetcher
}
func (f *URLsFetcher) Fetch(url string, chResult chan<- Result) {
	start := f.Time.Now()
	resp, err := f.Http.Get(url)
	if err != nil {
		chResult <- Result{Text: "", Error: fmt.Errorf("Fetch: while getting url=\"%s\" error occurs: \"%v\"", url, err)}
		return
	}
	nbytes, err := f.Http.Size(ioutil.NopCloser(resp.Body))
	if err != nil {
		chResult <- Result{Text: "", Error: fmt.Errorf("Fetch: while reading response from url=\"%s\" error occurs: \"%v\"", url, err)}
		return
	}
	secs := f.Time.Since(start).Seconds()
	chResult <- Result{Text: fmt.Sprintf("Fetch: got response size=%d from url=\"%s\" which took=%.2f seconds", nbytes, url, secs), Error: nil}
}
func (f *URLsFetcher) FetchAll(urls []string) {
	start := f.Time.Now()
	chResult := make(chan Result)
	for _, url := range urls {
		go f.Fetch(url, chResult)
	}

	for range urls {
		f.handleResult(<-chResult)
	}

	fmt.Fprintf(f.outWrite, "Fetch All: %x urls took %.2f seconds\n", len(urls), f.Time.Since(start).Seconds())
}
