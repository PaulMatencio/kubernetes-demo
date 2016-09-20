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

func main() {

	proxyUrl, _ := url.Parse("http://proxylb.internal.epo.org:8080")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyURL(proxyUrl),
	}

	myurl := "https://golang.org"
	loop := 0
	mynum := 1
	if len(os.Args) > 1 {
		myurl = os.Args[1]
	}
	if len(os.Args) > 2 {
		mynum, _ = strconv.Atoi(os.Args[2])
	}
	start := time.Now()
	ch := make(chan string)

	for loop < mynum {
		go fetch(myurl, tr, ch) // start a go routine
		loop++
	}
	loop = 0
	for loop < mynum {
		fmt.Println(<-ch)
		loop++
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())

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
