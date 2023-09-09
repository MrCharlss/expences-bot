package main

import (
	"encoding/json"
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

type Period struct {
	gorm.Model
	ID        uint `gorm:"primaryKey;autoIncrement;"`
	IncomeQ1  string
	IncomeQ2  string
	Name      string
	ClosingQ1 int32
	ClosingQ2 int32
}

type Expence struct {
	gorm.Model
	ID         uint `gorm:"primaryKey;autoIncrement;"`
	CategoryID int
	Category   Category
	Name       string
	Amount     int
	PeriodID   int
	Period     Period
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
	PeriodID   int
	Period     Period
}

type Month struct {
	Month string
	ID    uint `gorm:"primaryKey;autoIncrement"`
}

var otherKeyboard = tgbot.NewInlineKeyboardMarkup(
	tgbot.NewInlineKeyboardRow(
		// tgbot.NewInlineKeyboardButtonURL("Categories"),
		tgbot.NewInlineKeyboardButtonSwitch("2sw", "open 2"),
		tgbot.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbot.NewInlineKeyboardRow(
		tgbot.NewInlineKeyboardButtonData("4", "4"),
		tgbot.NewInlineKeyboardButtonData("5", "5"),
		tgbot.NewInlineKeyboardButtonData("6", "6"),
	),
)
var numericKeyBoard = tgbot.NewReplyKeyboard(tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton("WOW")))
var inlineKeyboard = tgbot.NewInlineKeyboardMarkup(tgbot.NewInlineKeyboardRow(tgbot.NewInlineKeyboardButtonData("N", "5")))

type Seed struct {
	Categories []string `json:"categories"`
	Months     []string `json:"months"`
	Period     []struct {
		Name      string
		IncomeQ1  int32
		IncomeQ2  int32
		ClosingQ1 int32
		ClosingQ2 int32
	} `json:"period"`
}

func main() {
	fmt.Println("Works...")
	// db, error := gorm.Open(sqlite.Open("finance.db"), &gorm.Config{})
	if len(os.Args) > 1 {
		featFlag := os.Args[1]
		switch featFlag {
		case "this-is-args":
			fmt.Println("...")
			file, error := os.ReadFile("./seed.json")
			if error != nil {
				panic(error)
			}
			v := map[string]Seed{}
			error = json.Unmarshal(file, &v)
			if error != nil {
				panic(error)
			}
			fmt.Printf("%+v", v)

			return
		default:
			fmt.Println("default")

		}

		// error = db.AutoMigrate(&Category{}, &Expence{}, &BudgetItem{}, &Month{}, &Period{})
		// if error != nil {
		// 	panic(error)
		// }
	}

	// if error != nil {
	// 	panic(error)
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
		if update.CallbackQuery != nil {
			callback := tgbot.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			res, err := bot.Request(callback)
			if err != nil {
				panic(err)
			}
			fmt.Printf("\n this is callback %+v \n", res)
			msg := tgbot.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}

		}

		msg := tgbot.NewMessage(update.Message.Chat.ID, "")
		switch update.Message.Command() {
		case "help":
			msg.Text = "Nuevo Gasto: /gasto \nVer budget: /budget"
		case "budget":
			msg.Text = "botones..."
			msg.ReplyMarkup = inlineKeyboard
		case "gasto":

			msg.Text = "botones...2"
			msg.ReplyMarkup = otherKeyboard
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
