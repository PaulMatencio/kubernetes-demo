package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var (
	myurl, userid, password string
	token                   []byte
)

func usage() {
	usage := "\nUsage: login  -wwww  string -u  userid -p password" +
		"\n\nDefault Options:\n"
	fmt.Println(usage)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {

	flag.Usage = usage
	flag.StringVar(&myurl, "www", "", "web page to login")
	flag.StringVar(&userid, "c", "paul", "user idrntification ")
	flag.StringVar(&password, "p", "oct@@@1998", "user password")

	flag.Parse()
	if len(myurl) == 0 {
		usage()
	}

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
	myurl = myurl + "/login"
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", myurl, nil)
	req.SetBasicAuth(userid, password)
	/*
		form := url.Values{}
		form.Add("u", userid)
		form.Add("p", password)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	*/
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("login error:", err)
	} else {
		token, err = ioutil.ReadAll(resp.Body)
		fmt.Println(string(token))
	}

}
