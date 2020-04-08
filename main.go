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
	var parallelJobs int
	flag.IntVar(&parallelJobs, "parallel", 10, "Number of jobs to run in parallel")
	flag.Parse()

	urls := flag.Args()
	hashResponse(urls, getRequestBody, os.Stdout, parallelJobs)
}

func hashResponse(urls []string, getter requester, output io.Writer, nJobs int) {
	guard := make(chan struct{}, nJobs)
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
