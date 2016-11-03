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
	myurl, contentc, securec string
	concurrent, loop, wait   int
	nbytes                   int64
	mbytes                   []byte
	mytoken                  []byte
	content, secure          bool
	err                      error
)

func usage() {
	usage := "\nUsage: fetch -www <url> -S <true/false> -c <concurrent request> -l <number of loops> -W <wait time between 2 loops> -C <true/false>" +
		"\n\nDefault Options:\n"
	fmt.Println(usage)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {

	flag.Usage = usage
	flag.StringVar(&myurl, "www", "", "web page to be fetched")
	flag.StringVar(&securec, "S", "false", "secured web page")
	flag.IntVar(&concurrent, "c", 1, "number of concurrent requests")
	flag.IntVar(&loop, "l", 1, "number of times of concurrent requests")
	flag.IntVar(&wait, "W", 50, "Wait time (ms) between 2 concurrent requests")
	flag.StringVar(&contentc, "C", "false", "return the content if true otherwise the length of the web page")

	flag.Parse()
	if len(myurl) == 0 {
		usage()
	}
	content, _ = strconv.ParseBool(contentc)
	secure, _ = strconv.ParseBool(securec)
	waitTime := time.Duration(wait) * time.Millisecond
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//Proxy:           http.ProxyURL(proxyUrl),
	}

	// get the proxy environment variable
	proxy := os.Getenv("HTTPS_PROXY")
	if len(proxy) == 0 {
		proxy = os.Getenv("https_proxy")
	}
	if len(proxy) != 0 {
		if proxyUrl, _ := url.Parse(proxy); proxyUrl != nil {
			// update the transport to add Insecure connection and proxy
			tr.Proxy = http.ProxyURL(proxyUrl)
		}
	}

	if secure {
		if mytoken, err = ioutil.ReadFile(("/tmp/$$$mytoken")); err != nil {
			fmt.Println("Please login first", err)
			fmt.Println("Switch to unsecured request")
		} else {
			// fmt.Println(string(mytoken))
		}
	}

	// start := time.Now()
	ch := make(chan string)
	for k := 0; k < loop; k++ {
		start := time.Now()
		for i := 0; i < concurrent; i++ {
			go fetch(myurl, mytoken, tr, ch) // start a go routine
		}
		for i := 0; i < concurrent; i++ {
			fmt.Println(<-ch)
		}
		fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
		time.Sleep(waitTime)
	}
}

func fetch(url string, mytoken []byte, tr *http.Transport, ch chan<- string) {

	start := time.Now()
	client := &http.Client{Transport: tr}
	if len(mytoken) > 0 {
		url = url + "/secure"
	}
	req, _ := http.NewRequest("GET", url, nil)
	if len(mytoken) > 0 {
		req.Header.Add("Authorization", "Bearer "+string(mytoken))
		// fmt.Println(req.Header)
	}
	resp, err := client.Do(req)

	if err != nil {
		// request error then  immediat return
		ch <- fmt.Sprint(err) // send to channel ch  and return
		return
	}

	defer resp.Body.Close() // don't leak resources

	// read the resp Body
	if !content {
		nbytes, err = io.Copy(ioutil.Discard, resp.Body)
	} else {
		mbytes, err = ioutil.ReadAll(resp.Body)
	}

	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}

	secs := time.Since(start).Seconds()
	if !content {
		ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
	} else {
		ch <- fmt.Sprintf("%.2fs %s %s", secs, string(mbytes), url)
	}
}
