package bot

import (
	"time"

	"github.com/MatrixDeity/gocat/internal/cataas"
	"github.com/MatrixDeity/gocat/internal/telegram"
)

const (
	TelegramTokenVar = "GOCAT_TELEGRAM_TOKEN"
	StartCommand     = "/start"
	GoCommand        = "GO CAT"
)

type Bot struct {
	telegramClient *telegram.Client
	cataasClient   *cataas.Client
	lastUpdateID   int
}

func NewBot() *Bot {
	telegramToken := getEnvVarChecked(TelegramTokenVar)
	return &Bot{
		telegramClient: telegram.NewClient(telegramToken),
		cataasClient:   cataas.NewClient(),
		lastUpdateID:   -1,
	}
}

func (b *Bot) Run() {
	updatesChan := b.pollUpdates()
	for update := range updatesChan {
		message := update.Message
		if message == nil {
			continue
		}
		go b.processMessage(message)
	}
}

func (b *Bot) pollUpdates() chan *telegram.Update {
	updatesChan := make(chan *telegram.Update)

	go func() {
		b.fetchUpdates()
		for {
			for _, update := range b.fetchUpdates() {
				updatesChan <- update
			}
			time.Sleep(time.Second)
		}
	}()

	return updatesChan
}

func (b *Bot) fetchUpdates() []*telegram.Update {
	updates, err := b.telegramClient.GetUpdates(b.lastUpdateID)
	if err != nil {
		return []*telegram.Update{}
	}
	result := make([]*telegram.Update, 0, len(updates))
	for _, update := range updates {
		if b.lastUpdateID < update.ID {
			b.lastUpdateID = update.ID
			result = append(result, update)
		}
	}
	return result
}

func (b *Bot) processMessage(message *telegram.Message) {
	switch message.Text {
	case StartCommand, GoCommand:
		b.sendCatPhoto(message.From.ID)
	}
}

func (b *Bot) sendCatPhoto(chatID int) {
	b.telegramClient.SendChatAction(chatID, telegram.UploadPhotoChatAction)
	photo, _ := b.cataasClient.GetCatPhoto()
	b.telegramClient.SendPhoto(chatID, photo, []string{"GO CAT"})
}
