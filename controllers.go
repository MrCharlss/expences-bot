package main

import (
	"fmt"

	"gorm.io/gorm"
)

type Budget struct {
    ID int
    CategoryID []int
    Name []string
    Total []int
}
var budgetItems = "budget_items"
func GetBudget(db *gorm.DB) Budget  {
    fmt.Print("GetBudget")

    var budget Budget

    db.Raw("select category_id, sum(amount), name from budget_items where month_id = '9' group_by category_id").Scan(&budget)
    return budget
    
}
