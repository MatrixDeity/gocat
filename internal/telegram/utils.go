package telegram

import "net/http"

func makeHeader(contentType string) http.Header {
	header := http.Header{}
	header.Set("Content-Type", contentType)
	header.Set("Accept", "application/json")
	return header
}

func makeKeyboardMarkup(buttons []string) *ReplyKeyboardMarkup {
	keyboard := make([][]*KeyboardButton, len(buttons))
	for index, buttonLabel := range buttons {
		button := &KeyboardButton{Text: buttonLabel}
		keyboard[index] = []*KeyboardButton{button}
	}
	return &ReplyKeyboardMarkup{InlineKeyboard: keyboard, ResizeKeyboard: true}
}
