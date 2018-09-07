//DAO.go

package main

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
)

//Stock data
type Stock struct {
	Name  string `json:"name"`
	Price string
}

// StockProprety struct to store information of data
type StockProprety struct {
	StockCode          string `json:"stockcode"`
	StockName          string `json:"stockname"`
	StockCurrentPrice  string `json:"stockcurrentprice"`
	StockChange        string `json:"stockchange"`
	StockPrecentChange string `json:"stockprecentchange"`
	StockDayHighPrice  string `json:"stockdayhighprice"`
	StockDayLowPrice   string `json:"stockdaylowprice"`
	StockVolume        string `json:"stockvolume"`
}

// StockInformation is a map that contains StockNumber as key and Stockproprety map as value
type StockInformation struct {
	StockWhole map[string]StockProprety
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
	/*
		fmt.Println(stock1)
		// Get all the rows in the Sheet1.
		rows := xlsx.GetRows("工作表1")
		for _, row := range rows {
			for _, colCell := range row {
				fmt.Print(colCell, "\t")
			}
			fmt.Println()
		}
	*/
	return stock1
}

// ReadData from xlsx
func ReadData(xlsxName string, sheetName string) [][]string {

	xlsx, err := excelize.OpenFile(xlsxName)
	if err != nil {
		fmt.Println("Cannot find the file")
	}

	// check where's header
	var a = false

	for a != true {
		if xlsx.GetCellValue(sheetName, "A1") != "CALC_DATE" {
			xlsx.RemoveRow(sheetName, 0)
		} else {
			a = true
			break
		}
	}
	fmt.Println("H")
	rows := xlsx.GetRows(sheetName)

	return rows
}

// SetData1 is about
func SetData1() []StockProprety {

	rows := ReadData("RMIT 1.0.xlsx", "asset_universe")

	var stock []StockProprety

	for i := range rows {

		// header location
		if i == 0 {
			continue
		}

		if rows[i][0] != "" {
			stock = append(stock, StockProprety{
				StockCode:          rows[i][0],
				StockName:          rows[i][1],
				StockCurrentPrice:  rows[i][2],
				StockChange:        rows[i][3],
				StockPrecentChange: rows[i][4],
				StockDayHighPrice:  rows[i][5],
				StockDayLowPrice:   rows[i][6],
				StockVolume:        rows[i][7],
			})
		}
	}
	return stock
}

// SetData2 is about
func SetData2() StockInformation {

	rows := ReadData("mockdata.xlsx", "Sheet1")

	stockWhole := make(map[string]StockProprety)
	stockInformation := StockInformation{stockWhole}

	for i := range rows {

		if i == 0 {
			continue
		}
		if rows[i][0] != "" {
			var key = rows[i][1]
			stockWhole[key] = StockProprety{
				rows[i][0], rows[i][1], rows[i][2],
				rows[i][3], rows[i][4], rows[i][5],
				rows[i][6], rows[i][7],
			}
		}
	}
	return stockInformation
}
