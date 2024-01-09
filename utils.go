package main

import (
	"fmt"
	"net/http"
)

// SearchStringInSlice searches for a string in a slice of strings
func SearchStringInSlice(slice []string, target string) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}

func makeGetRequest(ssl bool, host string, port int) error {
	proto := "http://"
	if ssl {
		proto = "https://"
	}
	url := fmt.Sprintf("%s%s:%d", proto, host, port)
	_, err := http.Get(url)

	if err != nil {
		fmt.Printf("[!] Error Occurred while calling %s: %s\n", url, err)
	}

	return err
}
