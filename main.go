package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nlopes/slack"
)

func main() {
	var token, message, channel, username string
	var err error
	flag.StringVar(&token, "token", "", "token")
	flag.StringVar(&message, "message", "", "post message")
	flag.StringVar(&channel, "channel", "general", "posting channel name")
	flag.StringVar(&username, "username", "nopaste", "posting username")
	flag.Parse()

	if token == "" {
		fmt.Fprintln(os.Stderr, "token was required")
		os.Exit(1)
	}

	api := slack.New(token)

	channels, err := api.GetChannels(true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "get channels error: %s", err)
		os.Exit(1)
	}
	var channelInfo *slack.Channel
	for _, ch := range channels {
		if ch.Name == channel {
			channelInfo = &ch
			break
		}
	}
	if channelInfo == nil {
		fmt.Fprintln(os.Stderr, "channel not found", err)
		os.Exit(1)
	}

	var pretext string
	if hostname, err := os.Hostname(); err == nil {
		pretext = fmt.Sprintf("this nopaste posted by %s", hostname)
	}

	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "stdin read error: %s", err)
		os.Exit(1)
	}
	bs := string(b)
	fp := slack.FileUploadParameters{
		Title:          message,
		InitialComment: pretext,
		Content:        bs,
		Channels:       []string{channelInfo.ID},
	}
	file, err := api.UploadFile(fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file upload error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("Message successfully upload snippet: %s at %s", file.URL, file.Timestamp)
}
