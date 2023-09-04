package main

import (
	"fmt"
	"os"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	fmt.Println("Works...")
    token := os.Getenv("TEL_TOKEN")
    fmt.Println(token)
	bot, err := tgbot.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	updateConfig := tgbot.NewUpdate(0)
	updateConfig.Timeout = 30
	updatesChan := bot.GetUpdatesChan(updateConfig)

	for update := range updatesChan {
		if update.Message == nil {
			fmt.Printf("Message nil - %+v\n", update)
			continue
		}
		msg := tgbot.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err = bot.Send(msg); err != nil {
			panic(err)
		}
	}
	fmt.Printf("%+v", bot)

}
