package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	os.Setenv("SLACK_BOT_TOKEN", "xoxb-6461516691651-6465403985029-h37b2sFfhD76tLJae8JI57Il")
	os.Setenv("CHANNEL_ID", "C06DKF6M7QB")

	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	channelArr := []string{os.Getenv("CHANNEL_ID")}
	fileArr := []string{"SRT.pdf"}

	for _, file := range fileArr {
		params := slack.FileUploadParameters{
			Channels: channelArr,
			File:     file,
		}

		dt, err := api.UploadFile(params)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Name: %s, URL: %s\n", dt.Name, dt.URLPrivate)
	}
}
