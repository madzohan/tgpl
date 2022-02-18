package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/madzohan/tgpl/ch1/1_56_fetchSeries/fetch0"
)

func main() {
	urls := os.Args[1:]
	if len(urls) == 0 {
		fmt.Println("Please type in urls in args!")
		return
	}
	h := fetch0.NewHttp(http.Get)
	respGen := h.Fetch(urls)
	for _, url := range urls {
		respBody := <-respGen
		fmt.Printf("url=%s body=\n%s\n", url, respBody.Text)
	}
}
