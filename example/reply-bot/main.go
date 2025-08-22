package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/dghubble/oauth1"
	twitter "github.com/g8rswimmer/go-twitter/v2"
	"github.com/mattn/go-xsearch"
)

type authorize struct {
}

func (a authorize) Add(req *http.Request) {
}

func main() {
	var verbose bool
	var pattern, filter, reply string
	var from string
	var clientToken, clientSecret, accessToken, accessSecret string
	var duration time.Duration
	var jsonOut bool
	flag.StringVar(&pattern, "pattern", "ぬるぽ", "Pattern")
	flag.StringVar(&filter, "filter", "^ぬるぽ$", "Regexp filter")
	flag.StringVar(&reply, "reply", "ｶﾞｯ", "Reply")
	flag.StringVar(&from, "from", "", "From")
	flag.StringVar(&clientToken, "client-token", os.Getenv("REPLYBOT_CLIENT_TOKEN"), "Twitter ClientToken")
	flag.StringVar(&clientSecret, "client-secret", os.Getenv("REPLYBOT_CLIENT_SECRET"), "Twitter ClientSecret")
	flag.StringVar(&accessToken, "access-token", os.Getenv("REPLYBOT_ACCESS_TOKEN"), "Twitter AccessToken")
	flag.StringVar(&accessSecret, "access-secret", os.Getenv("REPLYBOT_ACCESS_SECRET"), "Twitter AccessSecret")
	flag.DurationVar(&duration, "duration", 20*time.Second, "Duration")
	flag.BoolVar(&jsonOut, "json", false, "Output JSON")
	flag.BoolVar(&verbose, "verbose", false, "Verbose")

	flag.Parse()

	client := &twitter.Client{
		Authorizer: authorize{},
		Client: oauth1.NewConfig(clientToken, clientSecret).Client(oauth1.NoContext, &oauth1.Token{
			Token:       accessToken,
			TokenSecret: accessSecret,
		}),
		Host: "https://api.twitter.com",
	}

	var err error
	var filterRe *regexp.Regexp
	if filter != "" {
		filterRe, err = regexp.Compile(filter)
		if err != nil {
			log.Fatal(err)
		}
	}

	jsonw := json.NewEncoder(os.Stdout)
	first := true
	latestTweetId := ""
	for {
		entries, err := xsearch.Search(filter, xsearch.WithLatestTweetId(latestTweetId), xsearch.WithRemoveMarker(true))
		if err != nil {
			log.Println(err)
			continue
		}

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].CreatedAt < entries[j].CreatedAt
		})
		for _, entry := range entries {
			if filterRe != nil && !filterRe.MatchString(entry.DisplayTextBody) {
				continue
			}
			if jsonOut {
				jsonw.Encode(entry)
			} else {
				fmt.Println(entry.ID, entry.ScreenName, entry.DisplayTextBody)
			}
			latestTweetId = entry.ID

			if first == false && reply != "" {
				if from == "" || entry.ScreenName == from {
					req := twitter.CreateTweetRequest{
						Text: reply,
						Reply: &twitter.CreateTweetReply{
							InReplyToTweetID: entry.ID,
						},
					}
					_, err = client.CreateTweet(context.Background(), req)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
		time.Sleep(duration)
		first = false
	}
}
