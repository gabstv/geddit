package geddit

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type method string

const (
	OAUTH_BASE_URL        = "https://oauth.reddit.com"
	GET            method = "GET"
	POST           method = "POST"
	PATCH          method = "PATCH"
)

type oauthRequest struct {
	accessToken string
	url         string
	useragent   string
	values      *url.Values
	action      method
}

func ourl(format string, args ...interface{}) string {
	return fmt.Sprintf(OAUTH_BASE_URL+format, args...)
}

func (r oauthRequest) getResponse() (*bytes.Buffer, error) {
	// Determine the HTTP action.
	var buffer bytes.Buffer
	var action, finalurl string
	if r.action == GET {
		action = "GET"
		finalurl = r.url
		if r.values != nil {
			finalurl = r.url + "?" + r.values.Encode()
		}
	} else if r.action == POST {
		action = "POST"
		finalurl = r.url
		if r.values != nil {
			buffer.WriteString(r.values.Encode())
		}
	} else {
		action = "PATCH"
		finalurl = r.url
		if r.values != nil {
			buffer.WriteString(r.values.Encode())
		}
	}

	log.Println("finalurl", finalurl)

	// Create a request and add the proper headers.
	req, err := http.NewRequest(action, finalurl, &buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", r.useragent)
	req.Header.Set("Authorization", "bearer "+r.accessToken)

	cl := &http.Client{
		Timeout: time.Second * 30,
	}

	// Handle the request
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	respbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(respbytes), nil
}
