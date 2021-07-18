package telegram

import "encoding/json"

type Result struct {
	Ok          bool            `json:"ok"`
	Raw         json.RawMessage `json:"result"`
	ErrorCode   int             `json:"error_code"`
	Description string          `json:"description"`
}

type Update struct {
	Message *Message `json:"message"`
	ID      int      `json:"update_id"`
}

type Message struct {
	Date int    `json:"date"`
	From *User  `json:"from"`
	ID   int    `json:"message_id"`
	Text string `json:"text"`
}

type User struct {
	ID       int    `json:"id"`
	IsBot    bool   `json:"is_bot"`
	Username string `json:"username"`
}

type ReplyKeyboardMarkup struct {
	InlineKeyboard [][]*KeyboardButton `json:"keyboard"`
	ResizeKeyboard bool                `json:"resize_keyboard"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

type getUpdatesData struct {
	Offset int `json:"offset"`
}

type sendMessageData struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

type sendChatActionData struct {
	ChatID int    `json:"chat_id"`
	Action string `json:"action"`
}
