package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mrcharlss/finance-bot/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Category struct {
	Name string
	ID   uint `gorm:"primaryKey;autoIncrement"`
}

type Period struct {
	gorm.Model
	ID        uint `gorm:"primaryKey;autoIncrement;"`
	IncomeQ1  int32
	IncomeQ2  int32
	Name      string
	ClosingQ1 int32
	ClosingQ2 int32
	MonthID   int
	Month     Month
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
	MonthID    int
	Month      Month
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
	MonthID    int
	PeriodID   int
	Period     Period
}

type Month struct {
	Month string
	ID    uint `gorm:"primaryKey;autoIncrement"`
}

var otherKeyboard = tgbot.NewInlineKeyboardMarkup(
	tgbot.NewInlineKeyboardRow(
		tgbot.NewInlineKeyboardButtonData("Buget", "budget"),
		tgbot.NewInlineKeyboardButtonData("Gastos", "gastos"),
	),
)
var numericKeyBoard = tgbot.NewReplyKeyboard(tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton("WOW")))
var inlineKeyboard = tgbot.NewInlineKeyboardMarkup(tgbot.NewInlineKeyboardRow(tgbot.NewInlineKeyboardButtonData("N", "5")))

func main() {
	fmt.Println("Running")
	db, error := gorm.Open(sqlite.Open("finance.db"), &gorm.Config{})
	if len(os.Args) > 1 {
		featFlag := os.Args[1]
		switch featFlag {
		case "migrate":
			error = db.AutoMigrate(&Category{}, &Expence{}, &BudgetItem{}, &Month{}, &Period{}, &Month{})
			if error != nil {
				panic(error)
			}
			seed, err := utils.GetSeed()
			if err != nil {
				panic(err)
			}
			for _, item := range seed.BudgetItems {
				if err := db.Create(&BudgetItem{
					Name:       item.Name,
					Amount:     item.Amount,
					CategoryID: item.CatedoryId,
					Fortnight:  item.Fortnight,
					MonthID:    item.MonthId,
					PeriodID:   item.PeriodId,
				}).Error; err != nil {
					panic(err)
				}

			}
			for _, item := range seed.Months {
				if err := db.Create(&Month{
					Month: item,
				}).Error; err != nil {
					panic(err)
				}
			}
			for _, item := range seed.Expenses {
				if err := db.Create(&Expence{
					MonthID:    item.MonthId,
					Name:       item.Name,
					CategoryID: item.CatedoryId,
					PeriodID:   item.PeriodId,
					Amount:     item.Amount,
				}).Error; err != nil {
					panic(err)
				}
			}
			for _, item := range seed.Periods {
				if err := db.Create(&Period{
					MonthID:   item.MonthId,
					Name:      item.Name,
					IncomeQ1:  item.Incomeq1,
					IncomeQ2:  item.Incomeq2,
					ClosingQ1: item.Closingq1,
					ClosingQ2: item.Closingq2,
				}).Error; err != nil {
					panic(err)
				}
			}
			for _, item := range seed.Categories {
				if err := db.Create(&Category{
					Name: item,
				}).Error; err != nil {
					panic(err)
				}
			}
			fmt.Print("done migrating")

			// fmt.Println(seed.BudgetItems)
			// fmt.Println(seed.Months)
			// fmt.Println(seed.Categories)
			// fmt.Println(seed.Expenses)
			// fmt.Println(seed.Periods)
			if err != nil {
				panic(err)
			}
			// var cats []Category
			// for _, cat := range seed.Categories {
			//              cats = append(cats, Category{
			// 		Name: cat,
			// 	})
			// }
			// rows := db.Table("categories").Create(&cats).RowsAffected
			// fmt.Printf("%v - %+v", rows, cats)

			return
		default:
			fmt.Println("default")

		}

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
	// bot.Debug = true
	updateConfig := tgbot.NewUpdate(0)
	updateConfig.Timeout = 30
	updatesChan := bot.GetUpdatesChan(updateConfig)

	for update := range updatesChan {
		// fmt.Printf("%+v", update.CallbackQuery)
		if update.Message == nil {
			fmt.Printf("Message nil - %+v\n", update)
		}

		if update.CallbackQuery != nil {
			fmt.Println("CALLBACK")
			msg := tgbot.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
			switch update.CallbackData() {
			case "gastos":
				msg.Text = "This is gastos"
				msg.ReplyMarkup = inlineKeyboard

			default:
				msg.Text = "No valid option"
			}

			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
			continue

		}

		msg := tgbot.NewMessage(update.Message.Chat.ID, "")
		switch update.Message.Command() {
		case "help":
			msg.Text = "Nuevo Gasto: /gasto \nVer budget: /budget"
		case "budget":
			msg.Text = "botones..."
			msg.ReplyMarkup = numericKeyBoard
		case "gasto":

			msg.Text = "Acciones:"
			msg.ReplyMarkup = otherKeyboard
		case "close":
			msg.Text = "closing"
			msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
		case "consultar":
			fmt.Println("consultar")
			data := utils.GetExpences(db)
			var resString strings.Builder
			var total int = 0
			for _, item := range data {
				resString.WriteString(fmt.Sprintf("Name: %v\nAmount: $%v\n- - -\n", item.Name, item.Amount))

				total += item.Amount
			}
			resString.WriteString(fmt.Sprintf("Total: $%v", total))

			msg.Text = resString.String()
		default:
			msg.Text = "Comando desconocido"

		}

		if _, err = bot.Send(msg); err != nil {
			panic(err)
		}
	}
	fmt.Printf("%+v", bot)

}
