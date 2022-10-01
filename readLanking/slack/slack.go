package slack

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

type slack struct {
	params params
	url    string
}

type params struct {
	Text      string `json:"text"`
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	IconURL   string `json:"icon_url"`
	Channel   string `json:"channel"`
}

func NewSlack(url string, text string, username string, iconEmoji string, iconURL string, channel string) *slack {
	p := params{
		Text:      text,
		Username:  username,
		IconEmoji: iconEmoji,
		IconURL:   iconURL,
		Channel:   channel}

	return &slack{
		params: p,
		url:    url}
}

func (s *slack) Send() error {
	params, _ := json.Marshal(s.params)
	resp, err := http.PostForm(
		s.url,
		url.Values{"payload": {string(params)}},
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func Excute(message string) error {
	var slackURL = os.Getenv("SLACK_URL")
	var channel = os.Getenv("CHANNEL")
	var userName = os.Getenv("USER_NAME")
	var iconEmoji = os.Getenv("ICON_EMOJI")
	var iconURL = os.Getenv("ICON_URL")

	sl := NewSlack(slackURL, message, userName, iconEmoji, iconURL, channel)
	err := sl.Send()
	if err != nil {
		return err
	}
	return nil
}
