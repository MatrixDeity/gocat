package telegram

import (
	"net/http"
)

const BaseURL = "https://api.telegram.org"

const (
	TypingChatAction      = "typing"
	UploadPhotoChatAction = "upload_photo"
)

type Client struct {
	url       string
	token     string
	rawClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token:     token,
		url:       BaseURL + "/bot" + token,
		rawClient: &http.Client{},
	}
}
