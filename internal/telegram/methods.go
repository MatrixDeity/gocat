package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func (c *Client) GetUpdates(offset int64) ([]*Update, error) {
	body, _ := json.Marshal(getUpdatesData{offset})

	response, err := c.makeRequest(http.MethodPost, "/getUpdates", bytes.NewReader(body), makeHeader("application/json"))
	if err != nil {
		return nil, err
	}

	updates := []*Update{}
	json.Unmarshal(response.Raw, &updates)
	return updates, nil
}

func (c *Client) SendMessage(chatID int64, text string) (*Message, error) {
	body, _ := json.Marshal(sendMessageData{chatID, text})

	response, err := c.makeRequest(http.MethodPost, "/sendMessage", bytes.NewReader(body), makeHeader("application/json"))
	if err != nil {
		return nil, err
	}

	message := &Message{}
	json.Unmarshal(response.Raw, message)
	return message, nil
}

func (c *Client) SendPhoto(chatID int64, photo []byte, buttons []string) (*Message, error) {
	body := &bytes.Buffer{}
	mpwriter := multipart.NewWriter(body)

	mpwriter.WriteField("chat_id", fmt.Sprint(chatID))

	formFile, _ := mpwriter.CreateFormFile("photo", "cat.png")
	formFile.Write(photo)

	formButton, _ := mpwriter.CreateFormField("reply_markup")
	keyboard := makeKeyboardMarkup(buttons)
	buttonsBytes, _ := json.Marshal(keyboard)
	formButton.Write(buttonsBytes)

	mpwriter.Close()

	response, err := c.makeRequest(http.MethodPost, "/sendPhoto", body, makeHeader(mpwriter.FormDataContentType()))
	if err != nil {
		return nil, err
	}

	message := &Message{}
	json.Unmarshal(response.Raw, message)
	return message, nil
}

func (c *Client) SendChatAction(chatID int64, action string) error {
	body, _ := json.Marshal(sendChatActionData{chatID, action})

	_, err := c.makeRequest(http.MethodPost, "/sendChatAction", bytes.NewReader(body), makeHeader("application/json"))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) makeRequest(method string, path string, body io.Reader, header http.Header) (*Result, error) {
	url := c.url + path
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("telegram request error: %w", err)
	}
	request.Header = header

	response, err := c.rawClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("telegram API error: %w", err)
	}
	defer response.Body.Close()

	result := &Result{}
	err = json.NewDecoder(response.Body).Decode(result)
	if err != nil {
		return nil, fmt.Errorf("telegram response decode error: %w", err)
	}
	if !result.Ok {
		return nil, fmt.Errorf("telegram API error: %s - %d - %s", url, result.ErrorCode, result.Description)
	}
	return result, nil
}
