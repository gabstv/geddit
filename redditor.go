// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geddit

import (
	"encoding/json"
	"fmt"
	//"log"
	"net/url"
	"strconv"
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

func (r *OARedditor) Submitted(hideVotedLinks bool, limit int, popsort popularitySort, agesort ageSort, params ...Param) ([]Submission, error) {
	vals := url.Values{}
	if !hideVotedLinks {
		vals.Set("show", "all")
	}
	vals.Set("limit", strconv.Itoa(limit))
	vals.Set("sort", string(popsort))
	vals.Set("t", string(agesort))
	vals.Set("username", r.Name)
	for _, v := range params {
		vals.Set(v.Key, v.Value)
	}
	body, err := r.session.Get(&vals, "/user/%s/submitted", r.Name)
	if err != nil {
		return nil, err
	}

	//log.Println(body.String())

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
