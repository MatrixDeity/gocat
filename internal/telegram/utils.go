package telegram

import "net/http"

func makeHeader(contentType string) http.Header {
	header := http.Header{}
	header.Set("Content-Type", contentType)
	header.Set("Accept", "application/json")
	return header
}

func makeKeyboardMarkup(buttons []string) *ReplyKeyboardMarkup {
	keyboard := make([][]*KeyboardButton, 1)
	keyboard[0] = make([]*KeyboardButton, 0, len(buttons))
	for _, buttonLabel := range buttons {
		button := &KeyboardButton{Text: buttonLabel}
		keyboard[0] = append(keyboard[0], button)
	}
	return &ReplyKeyboardMarkup{InlineKeyboard: keyboard, ResizeKeyboard: true}
}
