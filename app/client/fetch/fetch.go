package main

import (
	"crypto/tls"
	"flag"
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
	myurl, contentc        string
	concurrent, loop, wait int
	nbytes                 int64
	mbytes                 []byte
	content                bool
)

func usage() {
	usage := "\nUsage: fetch  -wwww  string -c  integer -l integer -W integer -content boolean" +
		"\n\nDefault Options:\n"
	fmt.Println(usage)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {

	flag.Usage = usage
	flag.StringVar(&myurl, "www", "", "web page to be fetched")
	flag.IntVar(&concurrent, "c", 1, "number of concurrent requests")
	flag.IntVar(&loop, "l", 1, "number of times of concurrent requests")
	flag.IntVar(&wait, "W", 50, "Wait time (ms) between 2 concurrent requests")
	flag.StringVar(&contentc, "content", "false", "return the content if true otherwise the length of the web page")

	flag.Parse()
	if len(myurl) == 0 {
		usage()
	}
	content, _ = strconv.ParseBool(contentc)
	fmt.Println(content)
	waitTime := time.Duration(wait) * time.Millisecond

	// get the proxy environment variable
	proxy := os.Getenv("HTTPS_PROXY")
	if len(proxy) == 0 {
		proxy = os.Getenv("https_proxy")
	}

	proxyUrl, _ := url.Parse(proxy)
	if proxyUrl == nil {
		fmt.Println("set HTTPS_PROXY or https_proxy environment variable")
		os.Exit(1)
	}
	// update the transport to add Insecure connection and proxy
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyURL(proxyUrl),
	}

	// start := time.Now()
	ch := make(chan string)
	for k := 0; k < loop; k++ {
		start := time.Now()
		for i := 0; i < concurrent; i++ {
			go fetch(myurl, tr, ch) // start a go routine
		}

		for i := 0; i < concurrent; i++ {
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
	if !content {
		nbytes, err = io.Copy(ioutil.Discard, resp.Body)
	} else {
		mbytes, err = ioutil.ReadAll(resp.Body)
	}
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	if !content {
		ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
	} else {
		ch <- fmt.Sprintf("%.2fs %7d %s", secs, string(mbytes), url)
	}
}
