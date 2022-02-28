package main

import (
	"fmt"
	"os"

	"github.com/madzohan/tgpl/ch1/1_56_fetchSeries/fetchall"
)

func main() {
	urls := os.Args[1:]
	if len(urls) == 0 {
		fmt.Println("Please type in urls in args!")
		return
	}
	f := fetchall.NewURLsFetcher(fetchall.DefaultTime, fetchall.DefaultHttp, nil, fetchall.DefaultOutWrite, fetchall.DefaultErrWrite)
	f.FetchAll(urls)
}
