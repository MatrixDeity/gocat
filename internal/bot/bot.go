package bot

import (
	"log"
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
	lastUpdateID   int
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
	log.Println("GoCat is running")
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

	log.Println("GoCat stopped")
}

func (b *Bot) pollUpdates() chan *telegram.Update {
	updatesChan := make(chan *telegram.Update)

	go func() {
		b.fetchUpdates()
		for {
			select {
			case <-b.stopChan:
				close(updatesChan)
				return
			default:
			}

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
	chatID := message.From.ID
	logForChat(chatID, "process request")
	switch message.Text {
	case StartCommand, GoCommand:
		b.sendCatPhoto(chatID)
	default:
		logForChat(chatID, "bad command: "+message.Text)
	}
	logForChat(chatID, "done")
}

func (b *Bot) sendCatPhoto(chatID int) {
	defer b.goroutines.Done()

	err := b.telegramClient.SendChatAction(chatID, telegram.UploadPhotoChatAction)
	logForChat(chatID, err)

	photo, err := b.cataasClient.GetCatPhoto()
	logForChat(chatID, err)

	_, err = b.telegramClient.SendPhoto(chatID, photo, []string{GoCommand})
	logForChat(chatID, err)
}

func logForChat(chatId int, printable interface{}) {
	if printable != nil {
		log.Printf("chat:%d %s", chatId, printable)
	}
}
