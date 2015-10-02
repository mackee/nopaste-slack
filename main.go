package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nlopes/slack"
)

const defaultChannel = "general"

func main() {
	var token, message, channel, group, username string
	var err error
	flag.StringVar(&token, "token", "", "token")
	flag.StringVar(&message, "message", "", "post message")
	flag.StringVar(&channel, "channel", "", "posting channel name")
	flag.StringVar(&group, "group", "", "posting group name")
	flag.StringVar(&username, "username", "", "posting username")
	flag.Parse()

	if token == "" {
		fmt.Fprintln(os.Stderr, "token was required")
		os.Exit(1)
	}
	if channel != "" && group != "" {
		fmt.Fprintln(os.Stderr, "Can't specify both the channel and the group at once.")
		os.Exit(1)
	}
	if channel == "" && group == "" {
		channel = defaultChannel
	}

	api := slack.New(token)

	var channelID string = ""
	if channel != "" {
		channels, err := api.GetChannels(true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "get channels error: %s", err)
			os.Exit(1)
		}
		for _, ch := range channels {
			if ch.Name == channel {
				channelID = ch.ID
				break
			}
		}
	}
	if group != "" {
		groups, err := api.GetGroups(true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "get groups error: %s", err)
			os.Exit(1)
		}
		for _, ch := range groups {
			if ch.Name == group {
				channelID = ch.ID
				break
			}
		}
	}

	if channelID == "" {
		fmt.Fprintln(os.Stderr, "channel not found", err)
		os.Exit(1)
	}

	var pretext string
	if hostname, err := os.Hostname(); err == nil {
		if username == "" {
			pretext = fmt.Sprintf("this nopaste posted by %s", hostname)
		} else {
			pretext = fmt.Sprintf("this nopaste posted by %s@%s", username, hostname)
		}
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
		Channels:       []string{channelID},
	}
	file, err := api.UploadFile(fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file upload error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("Message successfully upload snippet: %s at %s", file.URL, file.Timestamp)
}
