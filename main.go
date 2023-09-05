package main

import (
	"fmt"
	"os"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type Category struct {
	Name string
	ID   uint `gorm:"primaryKey;autoIncrement"`
}

type Expence struct {
	gorm.Model
	ID         uint `gorm:"primaryKey;autoIncrement;"`
	CategoryID int
	Category   Category
	Name       string
	Amount     int
}

type BudgetItem struct {
	gorm.Model
	ID         uint `gorm:"primaryKey;autoIncrement;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Name       string
	Amount     int
	CategoryID int
	Category   Category
	Fortnight  int
	Month      Month
	MonthId    int
}

type Month struct {
	Month string
	ID    uint `gorm:"primaryKey;autoIncrement"`
}

var numericKeyBoard = tgbot.NewReplyKeyboard(tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton("WOW")))
var inlineKeyboard = tgbot.NewInlineKeyboardMarkup(tgbot.NewInlineKeyboardRow(tgbot.NewInlineKeyboardButtonData("N", "5")))

func main() {
	fmt.Println("Works...")
	// db, error := gorm.Open(sqlite.Open("finance.db"), &gorm.Config{})
	// error = db.AutoMigrate(&Category{}, &Expence{}, &BudgetItem{}, &Month{})
	// if error != nil {
	//     panic(error)
	// }

	// if error != nil {
	//     panic(error)
	// }
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
			// fmt.Printf("Message nil - %+v\n", update)
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbot.NewMessage(update.Message.Chat.ID, "")
		switch update.Message.Command() {
		case "help":
			msg.Text = "Nuevo Gasto: /gasto \nVer budget: /budget"
		case "budget":
            msg.Text = "botones..."
			msg.ReplyMarkup = inlineKeyboard
		case "close":
            msg.Text = "closing"
			msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
		default:
			msg.Text = "Comando desconocido"

		}

		if _, err = bot.Send(msg); err != nil {
			panic(err)
		}
	}
	fmt.Printf("%+v", bot)

}
