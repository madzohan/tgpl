package tests

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/madzohan/tgpl/ch1/1_56_fetchSeries/fetch"
)

func mock_get_page(url string) (resp *http.Response, err error) {
	if url == "error" {
		err = errors.New("some error")
	}
	resp = &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString("body")),
	}
	return
}

func TestFetch(t *testing.T) {
	type args struct {
		urls []string
	}

	tests := []struct {
		name         string
		args         args
		wantRespBody string
	}{
		// {name: "ok", args: args{urls: []string{"https://ifconfig.co/ip"}}, wantRespBody: "188.163.15.24"},
		{name: "ok", args: args{urls: []string{"test"}}, wantRespBody: "body\n"},
		{name: "error", args: args{urls: []string{"error"}}, wantRespBody: ""},
		{name: "multiple", args: args{urls: []string{"test1, test2"}}, wantRespBody: "body\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := fetch.NewHttp(mock_get_page)
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
