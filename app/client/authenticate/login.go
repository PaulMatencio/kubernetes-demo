package main

import (
	"crypto/tls"
	json "encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	gopass "github.com/howeyc/gopass"
	//terminal "golang.org/x/crypto/ssh/terminal"
)

var (
	myurl, userid, password string
	token, pass             []byte
	err                     error
	resetc                  string
)

type MyToken struct {
	Token string
}

func usage() {
	usage := "\nUsage: authenticate  -www  <url> -u  <userid> " +
		"\n\nDefault Options:\n"
	fmt.Println(usage)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.StringVar(&myurl, "www", "", "Web server url")
	flag.StringVar(&userid, "u", "", "User identification ")
	flag.StringVar(&resetc, "r", "", "Reset jwt token")
	flag.Parse()

	if len(resetc) > 0 {
		os.Remove("/tmp/$$$mytoken")
		os.Exit(0)
	}
	if len(myurl) == 0 || len(userid) == 0 {
		usage()
	}

	// prompt for password
	fmt.Printf("Enter your password:")
	if pass, err = gopass.GetPasswd(); err == nil {
		password = string(pass)
	} else {
		fmt.Println("error reading password", err)
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// get the proxy environment variable
	proxy := os.Getenv("HTTPS_PROXY")
	if len(proxy) == 0 {
		proxy = os.Getenv("https_proxy")
	}

	if len(proxy) != 0 {
		proxyUrl, _ := url.Parse(proxy)
		if proxyUrl != nil {
			// update the transport to add Insecure connection and proxy
			tr.Proxy = http.ProxyURL(proxyUrl)
		}
	}

	myurl = myurl + "/login"

	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", myurl, nil)
	req.SetBasicAuth(userid, password)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("login error:", err)
	} else {
		if token, err = ioutil.ReadAll(resp.Body); err == nil {
			mytok := MyToken{}
			if err = json.Unmarshal(token, &mytok); err == nil {
				fmt.Println(string(mytok.Token))
				if err = ioutil.WriteFile("/tmp/$$$mytoken", []byte(mytok.Token), 0700); err != nil {
					fmt.Println("Error saving token", err)
				}
			} else {
				fmt.Println("Invalid token data", err)
			}
		} else {
			fmt.Println("Error reading token", err)
		}
	}

}
