package tests

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	fetch0 "github.com/madzohan/tgpl/ch1/1_56_fetchSeries/fetch0"
	fetch1 "github.com/madzohan/tgpl/ch1/1_56_fetchSeries/fetch1"
	"github.com/madzohan/tgpl/ch1/1_56_fetchSeries/fetch2"
	"github.com/madzohan/tgpl/ch1/1_56_fetchSeries/fetch3"
	"github.com/madzohan/tgpl/ch1/1_56_fetchSeries/fetchall"
)

var getterErrorStr, readerErrorStr string = "getter error", "reader error"

func mockGetPage(url string) (resp *http.Response, err error) {
	bodyStr := "body"
	if url == getterErrorStr {
		err = fmt.Errorf(getterErrorStr)
	} else if url == readerErrorStr {
		bodyStr = readerErrorStr
	}
	resp = &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(bodyStr)),
	}
	return resp, err
}

func mockPrintResponse(response *fetch1.Response, outWrite io.Writer) (err error) {
	b, _ := ioutil.ReadAll(io.NopCloser(response.Body))
	err = fmt.Errorf(string(b))
	return
}

func mockMakeURL(url string) string {
	return url
}

func TestFetch0(t *testing.T) {
	t.Parallel()
	type args struct {
		urls []string
	}
	tests := []struct {
		name         string
		args         args
		wantRespBody string
	}{
		{name: "fetch0-ok", args: args{urls: []string{"test"}}, wantRespBody: "body\n"},
		{name: "fetch0-error", args: args{urls: []string{getterErrorStr}}, wantRespBody: ""},
		{name: "fetch0-multiple-urls", args: args{urls: []string{"test1, test2"}}, wantRespBody: "body\n"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := fetch0.NewHttp(mockGetPage)
			respGen := h.Fetch(tc.args.urls)
			for _, url := range tc.args.urls {
				respBody := <-respGen
				if url == "error" && respBody.Error == nil {
					t.Errorf("Fetch(%s) respBody.text = \"%s\" && respBody.err = %v, wantRespBody = \"%s\" && wantError",
						url, respBody.Text, respBody.Error, tc.wantRespBody)
				}
				if respBody.Text != tc.wantRespBody {
					t.Errorf("Fetch(%s) respBody.text = \"%s\" && respBody.err = %v, wantRespBody = \"%s\"",
						url, respBody.Text, respBody.Error, tc.wantRespBody)
					return
				}
			}

		})
	}
}

// chapter 1.5
func TestFetch1(t *testing.T) {
	type args struct {
		urls []string
	}
	type testArgs struct {
		name          string
		args          args
		wantRespBody  string
		wantErrorText string
	}
	tests := []testArgs{
		{"fetch1-ok", args{[]string{"test1"}}, "body\n", ""},
		{"fetch1-error-while-getting-url", args{[]string{getterErrorStr}}, "",
			fmt.Sprint("Fetch: while getting url=\"", getterErrorStr, "\" error occurs: \"", getterErrorStr, "\"\n")},
		{"fetch1-multiple-urls", args{[]string{"test1", "test2"}}, "body\nbody\n", ""},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			outReader, outWriter, errReader, errWriter := SetUp([]string{})
			fetch1.NewURLsFetcher(tc.args.urls, mockGetPage, outWriter, errWriter, nil, mockMakeURL).Fetch()
			TearDown(t, outReader, outWriter, errReader, errWriter, tc.wantRespBody, tc.wantErrorText)
		})
		t.Run(strings.ReplaceAll(tc.name, "1", "2"), func(t *testing.T) {
			outReader, outWriter, errReader, errWriter := SetUp([]string{})
			fetch1.NewURLsFetcher(tc.args.urls, mockGetPage, outWriter, errWriter, fetch2.PrintResponse, mockMakeURL).Fetch()
			TearDown(t, outReader, outWriter, errReader, errWriter, tc.wantRespBody, tc.wantErrorText)
		})
	}

	t.Run("fetch1-error-while-reading-response", func(t *testing.T) {
		tt := testArgs{"", args{[]string{readerErrorStr}}, "",
			fmt.Sprint("Fetch: while reading response from url=\"", readerErrorStr, "\" error occurs: \"", readerErrorStr, "\"\n")}
		outReader, outWriter, errReader, errWriter := SetUp([]string{})
		fetch1.NewURLsFetcher(tt.args.urls, mockGetPage, outWriter, errWriter, mockPrintResponse, mockMakeURL).Fetch()
		TearDown(t, outReader, outWriter, errReader, errWriter, tt.wantRespBody, tt.wantErrorText)
	})

	t.Run("fetch1-make-url", func(t *testing.T) {
		tt := testArgs{"", args{[]string{"test.com"}}, "body\nHTTP status code=0\n", ""}
		outReader, outWriter, errReader, errWriter := SetUp([]string{})
		fetch1.NewURLsFetcher(tt.args.urls, mockGetPage, outWriter, errWriter, fetch3.PrintResponse, nil).Fetch()
		TearDown(t, outReader, outWriter, errReader, errWriter, tt.wantRespBody, tt.wantErrorText)
	})
}

// chapter 1.6
// TimeIface test implementation
type MockTime struct{}

func (MockTime) Now() time.Time                  { return time.Unix(0, 0) }                 // 1970-01-01 00:00:00 +0000 UTC
func (MockTime) Since(s time.Time) time.Duration { return time.Duration(10) * time.Second } // 10s

// HttpIface test implementation
type MockHTTP struct {
	errGet  bool
	errSize bool
}

func (h MockHTTP) Get(url string) (resp *http.Response, err error) {
	resp, err = mockGetPage(url)
	if h.errGet {
		err = errors.New(getterErrorStr)
	}
	return
}
func (h MockHTTP) Size(responseBody io.ReadCloser) (numBytes int64, err error) {
	if h.errSize {
		err = fmt.Errorf(readerErrorStr)
	}
	return int64(10), err
}
func TestFetchAll(t *testing.T) {
	t.Parallel()
	res1 := "Fetch: got response size=10 from url=\"https://ifconfig.co/json\" which took=10.00 seconds\n"
	res2 := "Fetch: got response size=10 from url=\"https://reqres.in/api/products/3\" which took=10.00 seconds\n"
	getResSum := func(n int) string { return fmt.Sprintf("Fetch All: %d urls took 10.00 seconds\n", n) }
	type fields struct {
		Time         fetchall.TimeIface
		Http         fetchall.HttpIface
		handleResult fetchall.ResultHandler
		outWrite     io.Writer
		errWrite     io.Writer
	}
	type args struct {
		urls []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args

		wantText  string
		wantError string
	}{
		{"fetchall-real",
			fields{fetchall.DefaultTime, fetchall.DefaultHttp, nil, fetchall.DefaultOutWrite, fetchall.DefaultErrWrite},
			args{[]string{"fooBar"}},
			"",
			""},
		{"fetchall-ok",
			fields{MockTime{}, MockHTTP{}, nil, nil, nil},
			args{[]string{"https://ifconfig.co/json"}},
			res1 + getResSum(1),
			""},
		{"fetchall-ok-multiple",
			fields{MockTime{}, MockHTTP{}, nil, nil, nil},
			args{[]string{"https://ifconfig.co/json", "https://reqres.in/api/products/3"}},
			res1 + res2 + getResSum(2),
			""},
		{"fetchall-error-while-getting-response",
			fields{MockTime{}, MockHTTP{errGet: true}, nil, nil, nil},
			args{[]string{"https://ifconfig.co/json", "https://reqres.in/api/products/3"}},
			"Fetch All: 2 urls took 10.00 seconds\n",
			"Fetch: while getting url=\"https://reqres.in/api/products/3\" error occurs: \"getter error\"\n" +
				"Fetch: while getting url=\"https://ifconfig.co/json\" error occurs: \"getter error\"\n"},
		{"fetchall-error-while-reading-response",
			fields{MockTime{}, MockHTTP{errSize: true}, nil, nil, nil},
			args{[]string{"https://ifconfig.co/json", "https://reqres.in/api/products/3"}},
			"Fetch All: 2 urls took 10.00 seconds\n",
			"Fetch: while reading response from url=\"https://reqres.in/api/products/3\" error occurs: \"reader error\"\n" +
				"Fetch: while reading response from url=\"https://ifconfig.co/json\" error occurs: \"reader error\"\n"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			isRealRun := strings.HasSuffix(tc.name, "real")
			var (
				outReader *os.File
				outWriter *os.File
				errReader *os.File
				errWriter *os.File
			)
			if isRealRun {
				outReader, outWriter, errReader, errWriter = nil, nil, nil, nil
			} else {
				outReader, outWriter, errReader, errWriter = SetUp([]string{})
			}
			fetcher := fetchall.NewURLsFetcher(
				tc.fields.Time,
				tc.fields.Http,
				tc.fields.handleResult,
				outWriter,
				errWriter,
			)
			fetcher.FetchAll(tc.args.urls)
			if !isRealRun {
				TearDown(t, outReader, outWriter, errReader, errWriter, tc.wantText, tc.wantError)
			} else {
				resp := http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("")),
				}
				_, _ = fetcher.Http.Size(ioutil.NopCloser(resp.Body))
			}
		})
	}
}
