package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/fatih/color"

	"golang.org/x/net/proxy"
)

// Specify Tor proxy ip and port
var torProxy string = "socks5://127.0.0.1:9150" // 9150 w/ Tor Browser

func checkurl(webUrl string) Result {
	fmt.Printf("Checking %+v \n", webUrl)
	// Parse Tor proxy URL string to a URL type
	torProxyUrl, err := url.Parse(torProxy)
	if err != nil {
		return Result{webUrl, err, -1}
	}

	// Create proxy dialer using Tor SOCKS proxy
	torDialer, err := proxy.FromURL(torProxyUrl, proxy.Direct)
	if err != nil {
		return Result{webUrl, err, -1}
	}

	// Set up a custom HTTP transport to use the proxy and create the client
	torTransport := &http.Transport{Dial: torDialer.Dial}
	client := &http.Client{Transport: torTransport, Timeout: time.Second * 5}

	// Make request
	resp, err := client.Head(webUrl)
	if err != nil {
		return Result{webUrl, err, -1}
	}
	defer resp.Body.Close()

	return Result{webUrl, nil, resp.StatusCode}
}

func checkURL(urlvalue string, c chan Result) {
	c <- checkurl(urlvalue)
}

type Result struct {
	URL  string
	Err  error
	Code int
}

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Missing URL \n")
		os.Exit(-1)
	}
	file, err := os.Open("dir")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	webUrl, err := url.Parse(os.Args[1])

	if err != nil {
		fmt.Printf("Invalid URL %s \n", webUrl)
		os.Exit(-1)
	}

	okColor := color.New(color.FgGreen)
	errColor := color.New(color.FgRed)
	for scanner.Scan() {
		webUrl.Path = scanner.Text()
		if result := checkurl(webUrl.String()); result.Err != nil {
			errColor.Printf("[ERROR] on %s \n: %v \n", result.URL, result.Err)
		} else if result.Code < 400 {
			okColor.Printf("[OK - %d] Could be %s \n", result.Code, result.URL)
		} else {
			errColor.Printf("[NOT OK - %d] %s \n", result.Code, result.URL)
		}
	}

}
