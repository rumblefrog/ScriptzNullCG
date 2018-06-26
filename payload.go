package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Client - HTTP client for sending payload
var Client = &http.Client{
	Timeout: 30000000000,
}

// CreatePayload - Constructs payload for POST form query
func CreatePayload(t *Thread, r *Reply) bool {
	if r.Liked {
		return false
	}

	form := url.Values{}

	form.Add("_xfToken", t.XFToken)
	form.Add("_xfResponseType", "json")

	request, err := http.NewRequest("POST", Target+r.LikeHref, strings.NewReader(form.Encode()))

	log.Println("To " + Target + r.LikeHref)

	request.Header = Headers

	if err != nil {
		return false
	}

	res, err := Client.Do(request)

	if err != nil || res.StatusCode != 200 {
		return false
	}

	body, _ := ioutil.ReadAll(res.Body)

	log.Println(string(body))

	return true
}
