package cataas

import "net/http"

const BaseURL = "https://cataas.com"

type Client struct {
	url       string
	rawClient *http.Client
}

func NewClient() *Client {
	return &Client{
		url:       BaseURL,
		rawClient: &http.Client{},
	}
}
