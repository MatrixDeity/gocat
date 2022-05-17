package bot

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	lastUpdateID   int64
	stopChan       chan os.Signal
	goroutines     sync.WaitGroup
}

func NewBot(telegramToken string) *Bot {
	return &Bot{
		telegramClient: telegram.NewClient(telegramToken),
		cataasClient:   cataas.NewClient(),
		lastUpdateID:   -1,
		stopChan:       make(chan os.Signal, 1),
	}
}

func (b *Bot) Run() {
	log := newLogger(noChatID)
	log.info("GoCat is running")
	signal.Notify(b.stopChan, syscall.SIGINT, syscall.SIGTERM)

	updatesChan := b.pollUpdates()
	for update := range updatesChan {
		message := update.Message
		if message == nil {
			continue
		}
		b.goroutines.Add(1)
		go b.processMessage(message)
	}
	b.goroutines.Wait()

	log.info("GoCat stopped")
}

func (b *Bot) pollUpdates() chan *telegram.Update {
	updatesChan := make(chan *telegram.Update)

	go func(updatesChan chan<- *telegram.Update) {
		b.fetchUpdates()
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			select {
			case <-b.stopChan:
				ticker.Stop()
				close(updatesChan)
				return
			default:
			}

			for _, update := range b.fetchUpdates() {
				updatesChan <- update
			}
		}
	}(updatesChan)

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
	defer b.goroutines.Done()

	chatID := message.From.ID
	log := newLogger(chatID)

	log.info("process request")
	switch message.Text {
	case StartCommand, GoCommand:
		b.sendCatPhoto(chatID, log)
	default:
		log.info("bad command: " + message.Text)
	}
	log.info("done")
}

func (b *Bot) sendCatPhoto(chatID int64, log *logger) {
	err := b.telegramClient.SendChatAction(chatID, telegram.UploadPhotoChatAction)
	log.error(err)

	photo, err := b.cataasClient.GetCatPhoto()
	log.error(err)

	_, err = b.telegramClient.SendPhoto(chatID, photo, []string{GoCommand})
	log.error(err)
}
