package main

import (
	"amplify_bot/pkg/config"
	"amplify_bot/pkg/telegram"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cnf, err := config.CreateConfig()
	if err != nil {
		log.Fatalf("error in config loading: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(cnf.TgToken)
	if err != nil {
		log.Panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func(bot *tgbotapi.BotAPI) {
		defer wg.Done()
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates := bot.GetUpdatesChan(u)

		tg := telegram.NewTelegram(bot)

		for update := range updates {
			if update.Message != nil {
				tg.Process(&update)
			}
		}
	}(bot)

	wg.Wait()

}
