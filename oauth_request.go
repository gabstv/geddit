package geddit

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type oauthRequest struct {
	accessToken string
	url         string
	useragent   string
	values      *url.Values
}

func (r oauthRequest) getResponse() (*bytes.Buffer, error) {
	// Determine the HTTP action.
	var action, finalurl string
	if r.values == nil {
		action = "GET"
		finalurl = r.url
	} else {
		action = "POST"
		finalurl = r.url + "?" + r.values.Encode()
	}

	// Create a request and add the proper headers.
	req, err := http.NewRequest(action, finalurl, nil)
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
