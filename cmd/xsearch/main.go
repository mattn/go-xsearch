package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/mattn/go-xsearch"
)

const name = "xsearch"

const version = "0.0.3"

var revision = "HEAD"

func main() {
	var asjson bool
	var latestTweetId string
	var loop bool
	var duration time.Duration
	var showVersion bool
	flag.BoolVar(&asjson, "json", false, "Output as JSON")
	flag.StringVar(&latestTweetId, "latestTweetId", "", "Latest Tweet ID")
	flag.BoolVar(&loop, "loop", false, "Loop")
	flag.DurationVar(&duration, "duration", 20*time.Second, "Duration")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	word := strings.Join(flag.Args(), " ")

	jsonw := json.NewEncoder(os.Stdout)
	first := true
	for {
		var options xsearch.Options
		if latestTweetId != "" {
			options = append(options, xsearch.WithLatestTweetId(latestTweetId))
		}
		if asjson {
			options = append(options, xsearch.WithRemoveMarker(true))
		}
		entries, err := xsearch.Search(word, options...)
		if err != nil {
			log.Fatal(err)
		}

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].CreatedAt < entries[j].CreatedAt
		})

		if !first || !loop {
			if asjson {
				jsonw.Encode(entries)
			} else {
				re := regexp.MustCompile(`\tSTART\t[^\t]+\tEND\t`)
				for _, entry := range entries {
					text := entry.DisplayTextBody
					text = re.ReplaceAllStringFunc(text, func(s string) string {
						text = strings.TrimSuffix(strings.TrimPrefix(s, "\tSTART\t"), "\tEND\t")
						return color.RedString(text)
					})
					fmt.Fprintln(color.Output, entry.ID, text)
					latestTweetId = entry.ID
				}
			}
		}

		if !loop {
			break
		}

		time.Sleep(duration)
		first = false
	}
}
