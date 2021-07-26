package bot

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

const (
	formatString = "time=%s\tlogger_id=%d\tchat_id=%d\tmessage=%s\n"
	timeLayout   = "2006-01-02T15:04:05Z-0700"
	noChatID     = -1
)

type logger struct {
	id     int32
	chatID int64
}

func newLogger(chatID int64) *logger {
	return &logger{
		id:     rand.Int31(),
		chatID: chatID,
	}
}

func (l *logger) info(printable interface{}) {
	l.print(os.Stdout, printable)
}

func (l *logger) error(printable interface{}) {
	l.print(os.Stderr, printable)
}

func (l *logger) print(out io.Writer, printable interface{}) {
	if printable == nil {
		return
	}

	now := time.Now().Format(timeLayout)
	fmt.Fprintf(out, formatString, now, l.id, l.chatID, printable)
}
