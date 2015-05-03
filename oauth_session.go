// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package reddit implements an abstraction for the reddit.com API.
package geddit

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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
	defer resp.Body.Close()
	if err != nil {
		return err
	}

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

func (s OAuthSession) Me() (*Redditor, error) {
	req := &oauthRequest{
		accessToken: s.accessToken,
		url:         "https://oauth.reddit.com/api/v1/me",
		useragent:   s.useragent,
	}
	body, err := req.getResponse()
	if err != nil {
		return nil, err
	}

	type Response struct {
		Data Redditor
	}
	r := &Response{}
	err = json.NewDecoder(body).Decode(r)
	if err != nil {
		return nil, err
	}

	return &r.Data, nil
}
