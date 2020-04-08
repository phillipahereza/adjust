package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type requester func(url string) []byte

func main() {
	var parallelRequests int
	flag.IntVar(&parallelRequests, "parallel", 10, "Number of requests to make in parallel")
	flag.Parse()

	urls := flag.Args()

	if len(os.Args[1:]) == 0 {
		fmt.Printf("Usage:\n\nmyhttp -parallel 4 http://amazon.com google.com\n")
	}
	hashResponse(urls, getRequestBody, os.Stdout, parallelRequests)
}

func hashResponse(urls []string, getter requester, output io.Writer, nRequests int) {
	guard := make(chan struct{}, nRequests)
	resultsChan := make(chan string)

	for _, url := range urls {
		go func(u string) {
			guard <- struct{}{}
			responseBody := getter(u)
			if len(responseBody) == 0 {
				resultsChan <- fmt.Sprintf("%s %s", u, "Invalid Response")
				<-guard
				return
			}
			hashedBody := getHash(responseBody)
			resultsChan <- fmt.Sprintf("%s %s", checkURL(u), hashedBody)
			<-guard
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		fmt.Fprintf(output, "%s\n", <-resultsChan)
	}
}

func getHash(p []byte) string {
	hash := md5.New()
	hash.Write(p)
	hashedBytes := hash.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}

func getRequestBody(url string) []byte {
	url = checkURL(url)
	res, err := http.Get(url)
	if err != nil {
		return []byte{}
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}
	}
	defer res.Body.Close()
	return body
}

func checkURL(url string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	return "http://" + url
}
