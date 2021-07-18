package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/MatrixDeity/gocat/internal/bot"
)

type args struct {
	TelegramToken string
}

func parseArgs() (*args, error) {
	args := &args{}
	flag.StringVar(&args.TelegramToken, "token", "", "Token of Telegram bot")
	flag.Parse()

	if len(args.TelegramToken) == 0 {
		return nil, errors.New("pass Telegram token to -token argument")
	}
	return args, nil
}

func main() {
	args, err := parseArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		flag.Usage()
		os.Exit(1)
	}
	bot := bot.NewBot(args.TelegramToken)
	bot.Run()
}
