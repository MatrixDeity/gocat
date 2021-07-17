package cataas

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) GetCatPhoto() ([]byte, error) {
	response, err := c.makeResponse("GET", "/cat")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	photo, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed picture decode: %w", err)
	}
	return photo, err
}

func (c *Client) makeResponse(method string, path string) (*http.Response, error) {
	url := c.url + path
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("cataas request error: %w", err)
	}
	response, err := c.rawClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("cataas API error: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cataas API error: %s - %s", url, response.Status)
	}
	return response, nil
}
