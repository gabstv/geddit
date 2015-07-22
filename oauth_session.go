// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package reddit implements an abstraction for the reddit.com API.
package geddit

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// OAuthSession represents an OAuth session with reddit.com --
// all authenticated API calls are methods bound to this type.
type OAuthSession struct {
	username     string
	password     string
	clientID     string
	clientSecret string
	accessToken  string
	tokenType    string
	expiresIn    int
	scope        string
	useragent    string
}

// NewLoginSession creates a new session for those who want to log into a
// reddit account via OAuth.
func NewOAuthSession(username, password, useragent, clientID, clientSecret string) (*OAuthSession, error) {
	session := &OAuthSession{
		username:     username,
		password:     password,
		clientID:     clientID,
		clientSecret: clientSecret,
		useragent:    useragent,
	}

	err := session.newToken(&url.Values{
		"username":   {username},
		"password":   {password},
		"grant_type": {"password"},
	})
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *OAuthSession) newToken(postValues *url.Values) error {

	loginURL := "https://www.reddit.com/api/v1/access_token"
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(postValues.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set the auth header
	req.SetBasicAuth(s.clientID, s.clientSecret)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type Response struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope"`
	}

	r := &Response{}
	err = json.Unmarshal(body, r)
	if err != nil {
		return err
	}

	s.accessToken = r.AccessToken
	s.tokenType = r.TokenType
	s.expiresIn = r.ExpiresIn
	s.scope = r.Scope

	log.Println(s.scope)

	return nil
}

func (s OAuthSession) RevokeToken() error {
	revokeURL := "https://www.reddit.com/api/v1/revoke_token"
	postValues := &url.Values{
		"token":           {s.accessToken},
		"token_type_hint": {s.tokenType},
	}
	req, err := http.NewRequest("POST", revokeURL, strings.NewReader(postValues.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.clientID, s.clientSecret)

	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	// 401 returned if basic auth failed
	// 204 is returned even if given token is invalid
	if resp.StatusCode != http.StatusNoContent {
		return errors.New(resp.Status)
	}

	return nil
}

func (s *OAuthSession) Get(params *url.Values, urlformat string, urlvars ...interface{}) (*bytes.Buffer, error) {
	surl := ourl(urlformat, urlvars...)
	req := &oauthRequest{
		accessToken: s.accessToken,
		url:         surl,
		useragent:   s.useragent,
		action:      GET,
		values:      params,
	}
	return req.getResponse()
}

func (s *OAuthSession) Post(params *url.Values, urlformat string, urlvars ...interface{}) (*bytes.Buffer, error) {
	surl := ourl(urlformat, urlvars...)
	req := &oauthRequest{
		accessToken: s.accessToken,
		url:         surl,
		useragent:   s.useragent,
		action:      POST,
		values:      params,
	}
	return req.getResponse()
}

func (s *OAuthSession) Patch(params *url.Values, urlformat string, urlvars ...interface{}) (*bytes.Buffer, error) {
	surl := ourl(urlformat, urlvars...)
	req := &oauthRequest{
		accessToken: s.accessToken,
		url:         surl,
		useragent:   s.useragent,
		action:      PATCH,
		values:      params,
	}
	return req.getResponse()
}

func (s *OAuthSession) Me() (*OARedditor, error) {
	body, err := s.Get(nil, "/api/v1/me")
	if err != nil {
		return nil, err
	}

	oresp := &OARedditor{}
	log.Println(body.String())
	err = json.NewDecoder(body).Decode(oresp)
	if err != nil {
		return nil, err
	}
	// put the session in it
	oresp.session = s

	return oresp, nil
}

func (s *OAuthSession) User(username string) (*OARedditor, error) {
	body, err := s.Get(nil, "/user/%s/about", username)
	if err != nil {
		return nil, err
	}

	type Resp struct {
		Kind string     `json:"kind"`
		Data OARedditor `json:"data"`
	}

	oresp := &Resp{}
	log.Println(body.String())
	err = json.NewDecoder(body).Decode(oresp)
	if err != nil {
		return nil, err
	}
	// put the session in it
	oresp.Data.session = s

	log.Println(oresp.Data.String())

	return &oresp.Data, nil
}
