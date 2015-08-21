package notifications

import (
	"encoding/json"
	"net/http"
	"net/url"
)

var (
	WebhookURL      string = "https://hooks.slack.com/services/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"
	DefaultUserName string = "risu"
	DefaultChannel  string = "#infrastructure"
)

type Slack struct {
	Text     string `json:"text"`
	Username string `json:"username"`
	Channel  string `json:"channel"`
}

func Push(message string) {
	params, _ := json.Marshal(Slack{
		message,
		DefaultUserName,
		DefaultChannel,
	})

	http.PostForm(
		WebhookURL,
		url.Values{"payload": {string(params)}},
	)
}
