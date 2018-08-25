//DAO.go

package main

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
)

//Stock data
type Stock struct {
	name  string
	price string
}

var stockSlice = make([]Stock, 2)

func readExcel() Stock {
	xlsx, err := excelize.OpenFile("./dashboardData.xlsx")
	if err != nil {
		fmt.Println(err)
		return Stock{"1", "2"}
	}
	// Get value from cell by given worksheet name and axis.
	stock1 := Stock{xlsx.GetCellValue("工作表1", "A1"), xlsx.GetCellValue("工作表1", "A2")}
	fmt.Println(stock1)
	// Get all the rows in the Sheet1.
	rows := xlsx.GetRows("工作表1")
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
	return stock1
}
