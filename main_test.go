package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestCheckURL(t *testing.T) {
	urls := []string{
		"google.com",
		"https://amazon.com",
		"http://www.adjust.com",
		"www.facebook.com",
	}

	want := []string{
		"http://google.com",
		"https://amazon.com",
		"http://www.adjust.com",
		"http://www.facebook.com",
	}

	for i := 0; i < len(urls); i++ {
		got := checkURL(urls[i])
		if want[i] != got {
			t.Errorf("Wanted %s, got %s", want[i], got)
		}
	}
}

func TestHashResponse(t *testing.T) {
	fakeRequester := func(url string) []byte {
		return []byte(url)
	}

	urls := []string{
		"http://dux.dev",
		"htop.exe",
	}
	got := &bytes.Buffer{}

	hashResponse(urls, fakeRequester, got, 3)

	for _, url := range urls {
		hashedBody := getHash([]byte(url))
		if !strings.Contains(got.String(), fmt.Sprintf("%s %s", checkURL(url), hashedBody)) {
			log.Fatalf("got %s", got.String())
		}
	}
}
