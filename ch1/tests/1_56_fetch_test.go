package tests

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	fetch0 "github.com/madzohan/tgpl/ch1/1_56_fetchSeries/fetch0"
	fetch1 "github.com/madzohan/tgpl/ch1/1_56_fetchSeries/fetch1"
)

var getterErrorStr, readerErrorStr string = "getter error", "reader error"

func mockGetPage(url string) (resp *http.Response, err error) {
	bodyStr := "body"
	if url == getterErrorStr {
		err = errors.New(getterErrorStr)
	} else if url == readerErrorStr {
		bodyStr = readerErrorStr
	}
	resp = &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(bodyStr)),
	}
	return resp, err
}

func mockReadResponseBody(body io.Reader) (buf []byte, err error) {
	buf, err = io.ReadAll(body)
	if string(buf) == readerErrorStr {
		err = errors.New(readerErrorStr)
		return []byte{}, err
	}
	return
}

func TestFetch0(t *testing.T) {
	type args struct {
		urls []string
	}

	tests := []struct {
		name         string
		args         args
		wantRespBody string
	}{
		// {name: "ok", args: args{urls: []string{"https://ifconfig.co/ip"}}, wantRespBody: "XXX.XXX.XXX.XXX"},
		{name: "fetch0-ok", args: args{urls: []string{"test"}}, wantRespBody: "body\n"},
		{name: "fetch0-error", args: args{urls: []string{getterErrorStr}}, wantRespBody: ""},
		{name: "fetch0-multiple-urls", args: args{urls: []string{"test1, test2"}}, wantRespBody: "body\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := fetch0.NewHttp(mockGetPage)
			// h := NewHttp(http.Get)
			respGen := h.Fetch(tt.args.urls)
			for _, url := range tt.args.urls {
				respBody := <-respGen
				if url == "error" && respBody.Error == nil {
					t.Errorf("Fetch(%s) respBody.text = \"%s\" && respBody.err = %v, wantRespBody = \"%s\" && wantError", url, respBody.Text, respBody.Error, tt.wantRespBody)
				}
				if respBody.Text != tt.wantRespBody {
					t.Errorf("Fetch(%s) respBody.text = \"%s\" && respBody.err = %v, wantRespBody = \"%s\"", url, respBody.Text, respBody.Error, tt.wantRespBody)
					return
				}
			}

		})
	}
}

func TestFetch1(t *testing.T) {
	type args struct {
		urls []string
	}
	tests := []struct {
		name          string
		args          args
		wantRespBody  string
		wantErrorText string
	}{
		{"fetch1-ok", args{[]string{"test1"}}, "body\n", ""},
		{"fetch1-error-while-getting-url", args{[]string{getterErrorStr}}, "",
			strings.Join([]string{"Fetch: while getting url=\"", getterErrorStr, "\" error occurs: \"", getterErrorStr, "\"\n"}, "")},
		{"fetch1-error-while-reading-response", args{[]string{readerErrorStr}}, "",
			strings.Join([]string{"Fetch: while reading response from url=\"", readerErrorStr, "\" error occurs: \"", readerErrorStr, "\"\n"}, "")},
		{"fetch1-multiple-urls", args{[]string{"test1", "test2"}}, "body\nbody\n", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outReader, outWriter, errReader, errWriter := SetUp([]string{})
			fetch1.Fetch(tt.args.urls, mockGetPage, mockReadResponseBody, outWriter)
			TearDown(t, outReader, outWriter, errReader, errWriter, tt.wantRespBody, tt.wantErrorText)
		})
	}

}
