package xsearch

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Head struct {
	TotalResultsAvailable int `json:"totalResultsAvailable"`
	TotalResultsReturned  int `json:"totalResultsReturned"`
}

type Badge struct {
	Show  bool   `json:"show"`
	Type  string `json:"type"`
	Color string `json:"color"`
}

type Entry struct {
	ID                   string   `json:"id"`
	URL                  string   `json:"url"`
	DetailURL            string   `json:"detailUrl"`
	DetailQuoteURL       string   `json:"detailQuoteUrl"`
	Badge                Badge    `json:"badge"`
	DisplayText          string   `json:"displayText"`
	DisplayTextBody      string   `json:"displayTextBody"`
	DisplayTextFragments string   `json:"displayTextFragments"`
	DisplayTextEntities  string   `json:"displayTextEntities"`
	Urls                 []any    `json:"urls"`
	Hashtags             []any    `json:"hashtags"`
	HashtagUrls          struct{} `json:"hashtagUrls"`
	Mentions             []any    `json:"mentions"`
	MentionUrls          struct{} `json:"mentionUrls"`
	ReplyMentions        []any    `json:"replyMentions"`
	ReplyMentionUrls     struct{} `json:"replyMentionUrls"`
	CreatedAt            int      `json:"createdAt"`
	ReplyCount           int      `json:"replyCount"`
	ReplyURL             string   `json:"replyUrl"`
	RtCount              int      `json:"rtCount"`
	RtURL                string   `json:"rtUrl"`
	QtCount              int      `json:"qtCount"`
	LikesCount           int      `json:"likesCount"`
	LikesURL             string   `json:"likesUrl"`
	UserID               string   `json:"userId"`
	UserURL              string   `json:"userUrl"`
	Name                 string   `json:"name"`
	ScreenName           string   `json:"screenName"`
	ProfileImage         string   `json:"profileImage"`
	MediaType            []any    `json:"mediaType"`
	Media                []any    `json:"media"`
	PossiblySensitive    bool     `json:"possiblySensitive"`
	TweetThemeNormal     []any    `json:"tweetThemeNormal"`
	UserThemeNormal      []string `json:"userThemeNormal"`
	TwitterContextID     []any    `json:"twitterContextID"`
	VideoClassifyID      []any    `json:"videoClassifyId"`
	InReplyTo            string   `json:"inReplyTo"`
}

type Timeline struct {
	Head       Head    `json:"head"`
	Entry      []Entry `json:"entry"`
	MediaTweet bool    `json:"mediaTweet"`
}

type payload struct {
	Timeline Timeline `json:"timeline"`
}

type opt struct {
	latestTweetId string
}

type option func(*opt)

type Options []option

func WithLastTweetId(id string) option {
	return func(o *opt) {
		o.latestTweetId = id
	}
}

func Search(word string, options ...option) ([]Entry, error) {
	var o opt
	for _, option := range options {
		option(&o)
	}
	urlstring := `https://search.yahoo.co.jp/realtime/api/v1/pagination?p=` + url.QueryEscape(word)
	if o.latestTweetId != "" {
		urlstring += `&latestTweetId=` + url.QueryEscape(o.latestTweetId)
	}
	req, err := http.NewRequest(http.MethodGet, urlstring, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "curl/8.9.1")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	var p payload
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}
	return p.Timeline.Entry, nil
}
