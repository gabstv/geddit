// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geddit

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

type Redditor struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	LinkKarma    int     `json:"link_karma"`
	CommentKarma int     `json:"comment_karma"`
	Created      float32 `json:"created_utc"`
	Gold         bool    `json:"is_gold"`
	Mod          bool    `json:"is_mod"`
	Mail         *bool   `json:"has_mail"`
	ModMail      *bool   `json:"has_mod_mail"`
	InboxCount   int     `json:"inbox_count"`
}

// String returns the string representation of a reddit user.
func (r *Redditor) String() string {
	return fmt.Sprintf("%s (%d-%d)", r.Name, r.LinkKarma, r.CommentKarma)
}

// OAuth Redditor
type OARedditor struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	LinkKarma    int     `json:"link_karma"`
	CommentKarma int     `json:"comment_karma"`
	Created      float32 `json:"created_utc"`
	Gold         bool    `json:"is_gold"`
	Mod          bool    `json:"is_mod"`
	Mail         *bool   `json:"has_mail"`
	ModMail      *bool   `json:"has_mod_mail"`
	InboxCount   int     `json:"inbox_count"`
	//
	session *OAuthSession `json:"-"`
}

// String returns the string representation of a reddit user.
func (r *OARedditor) String() string {
	return fmt.Sprintf("%s (%d-%d)", r.Name, r.LinkKarma, r.CommentKarma)
}

func (r *OARedditor) Submitted() ([]Submission, error) {
	vals := url.Values{}
	vals.Set("show", "all")
	vals.Set("limit", "3")
	vals.Set("sort", "new")
	vals.Set("t", "all")
	vals.Set("username", r.Name)
	req := &oauthRequest{
		accessToken: r.session.accessToken,
		url:         ourl("/user/%s/submitted", r.Name),
		useragent:   r.session.useragent,
		values:      &vals,
	}
	log.Println(req.url)
	body, err := req.getResponse()
	if err != nil {
		//return nil, err
		log.Println("SUBMITTED ERROR", err)
		return nil, err
	}

	log.Println(body.String())

	type SubContainer struct {
		Data Submission `json:"data"`
	}
	type Resp struct {
		Data struct {
			Children []SubContainer `json:"children"`
		} `json:"data"`
	}

	result := &Resp{}

	err = json.NewDecoder(body).Decode(result)

	if err != nil {
		return nil, err
	}

	rs := make([]Submission, len(result.Data.Children))
	for k, v := range result.Data.Children {
		rs[k] = v.Data
	}
	return rs, nil
}
