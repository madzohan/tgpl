package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/madzohan/tgpl/ch1/1_56_fetchSeries/fetch1"
)

func main() {
	urls := os.Args[1:]
	if len(urls) == 0 {
		fmt.Println("Please type in urls in args!")
		return
	}
	fetch1.NewURLsFetcher(urls, http.Get, os.Stdout, os.Stderr, nil, nil).Fetch()
}
