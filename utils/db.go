package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"gorm.io/gorm"
)

type Period struct {
	Name      string `json:"name"`
	Incomeq1  int32  `json:"incomeQ1"`
	Incomeq2  int32  `json:"IncomeQ2"`
	Closingq1 int32  `json:"ClosingQ1"`
	Closingq2 int32  `json:"ClosingQ2"`
	MonthId   int    `json:"month_id"`
}

type Expense struct {
	Name       string `json:"name"`
	CatedoryId int    `json:"category_id"`
	Amount     int    `json:"amount"`
	PeriodId   int    `json:"period_id"`
	MonthId    int    `json:"month_id"`
}

type BudgetItem struct {
	Name       string `json:"name"`
	CatedoryId int    `json:"category_id"`
	Amount     int    `json:"amount"`
	PeriodId   int    `json:"period_id"`
	Fortnight  int    `json:"fortnight"`
	MonthId      int    `json:"month_id"`
}
type Seed struct {
	Categories  []string     `json:"categories"`
	Months      []string     `json:"months"`
	Periods     []Period     `json:"periods"`
	Expenses    []Expense    `json:"expenses"`
	BudgetItems []BudgetItem `json:"budgetItems"`
}

func GetSeed() (Seed, error) {

	file, error := os.ReadFile("./seed.json")
	if error != nil {
		panic(error)
	}
	v := Seed{}
	error = json.Unmarshal(file, &v)
	if error != nil {
		panic(error)
	}
	return v, nil
}

type Budget struct {
    ID int
    category_id int
    name string
}
var budgetItems = "budget_items"
func GetBudget(db *gorm.DB)  {
    fmt.Print("GetBudget")

    var budget Budget

    error  := db.Raw("select category_id, sum(amount), name from budget_items where month_id = '9' group by category_id;").Scan(&budget).Error
    if error != nil {
        panic(error)
    }
    fmt.Print(budget)
    
}
