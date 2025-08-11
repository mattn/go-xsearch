package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-xsearch"
)

func main() {
	var asjson bool
	var latestTweetId string
	flag.BoolVar(&asjson, "json", false, "Output as JSON")
	flag.StringVar(&latestTweetId, "latestTweetId", "", "Latest Tweet ID")
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	var options search.Options
	if latestTweetId != "" {
		options = append(options, search.WithLastTweetId(latestTweetId))
	}
	word := strings.Join(flag.Args(), " ")
	entries, err := search.Search(word)
	if err != nil {
		log.Fatal(err)
	}

	if asjson {
		json.NewEncoder(os.Stdout).Encode(entries)
		return
	}
	re := regexp.MustCompile(`\tSTART\t[^\t]+\tEND\t`)
	for _, entry := range entries {
		text := entry.DisplayTextBody
		text = re.ReplaceAllStringFunc(text, func(s string) string {
			text = strings.TrimSuffix(strings.TrimPrefix(s, "\tSTART\t"), "\tEND\t")
			return color.RedString(text)
		})
		fmt.Fprintln(color.Output, entry.ID, text)
	}
}
