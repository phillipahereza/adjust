package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	args := os.Args[1:]
	for _, u := range args {
		result := makeRequest(u)
		fmt.Println(result)
	}

}

func makeRequest(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	hash := md5.Sum(body)
	return fmt.Sprintf("%s %x", url, hash)
}
