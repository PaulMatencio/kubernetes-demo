package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var (
	myurl string
	wait  time.Duration
)

func main() {
	// get the proxy environment variable
	proxy := os.Getenv("HTTPS_PROXY")
	if len(proxy) == 0 {
		proxy = os.Getenv("https_proxy")
	}

	proxyUrl, _ := url.Parse(proxy)

	// update the transport to add Insecure connection and proxy
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyURL(proxyUrl),
	}
	argnum := len(os.Args) - 1
	if argnum == 0 {
		fmt.Println("client url  [concurrent_requests] [number_of_times]  [waittime_between_times]")
		os.Exit(1)
	}
	// myurl := "https://golang.org"
	mynum := 1 // concurrent go routines
	times := 1 // number of times
	wait := 100
	if argnum >= 1 {
		myurl = os.Args[1]
	}
	if argnum >= 2 {
		mynum, _ = strconv.Atoi(os.Args[2])
	}
	if argnum >= 3 {
		times, _ = strconv.Atoi(os.Args[3])
	}
	if argnum >= 4 {
		wait, _ = strconv.Atoi(os.Args[4])
	}
	waitTime := time.Duration(wait) * time.Millisecond

	// start := time.Now()
	ch := make(chan string)
	for k := 0; k < times; k++ {
		start := time.Now()
		for i := 0; i < mynum; i++ {
			go fetch(myurl, tr, ch) // start a go routine
		}

		for i := 0; i < mynum; i++ {
			fmt.Println(<-ch)
		}

		fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
		time.Sleep(waitTime)
	}

}

func fetch(url string, tr *http.Transport, ch chan<- string) {
	start := time.Now()

	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)

	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch return
	}
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}
