// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geddit

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	BASE_URL = "http://www.reddit.com"
)

func rurl(format string, args ...interface{}) string {
	return fmt.Sprintf(BASE_URL+format, args...)
}

type request struct {
	url       string
	values    *url.Values
	cookie    *http.Cookie
	useragent string
}

func (r request) getResponse() (*bytes.Buffer, error) {
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
	if r.cookie != nil {
		req.AddCookie(r.cookie)
	}
	req.Header.Set("User-Agent", r.useragent)

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
