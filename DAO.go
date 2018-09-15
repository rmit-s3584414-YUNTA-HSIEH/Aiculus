//DAO.go

package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

//Stock data
type Stock struct {
	Name  string
	Price string
}

// Type list
type (
	// StockProprety struct to store information of data
	StockProprety struct {
		Date        string `json:"date"`
		Isin        string
		Ric         string
		Name        string
		IsoCty      string
		Gics        string
		FloatMktCap string
	}

	// BenchMarkProprety struct to store information of Benchmark
	BenchMarkProprety struct {
		Date      string
		Isin      string
		Ric       string
		IsoCty    string
		IDXMktCap string
		Gicses    string
	}

	// StockCalculation struct to store all stock calculation tables
	StockCalculation struct {
		Name       string `json:"name"`
		Value      float64
		Presentage float64 `json:"presentage"`
	}
)

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
	rows := xlsx.GetRows(sheetName)

	return rows
}

// SetStockData function use to read data from excel and return the stock struct
func SetStockData() []StockProprety {

	// Read data from excel, pass xlsx filename and spreadsheet name
	rows := ReadData("Summary Data.xlsx", "Universe")

	var stock []StockProprety

	// Add data into struct
	for i := range rows {

		// header location
		if i == 0 {
			continue
		}

		// Check first element is not empty to add
		if rows[i][0] != "" {
			stock = append(stock, StockProprety{
				Date:        rows[i][0],
				Isin:        rows[i][1],
				Ric:         rows[i][2],
				Name:        rows[i][3],
				IsoCty:      rows[i][4],
				Gics:        rows[i][5],
				FloatMktCap: rows[i][9],
			})
		}
	}

	return stock

}

// SetBMData function use to read data from excel and return the benchmark struct
func SetBMData() []BenchMarkProprety {

	// Read data from excel, pass xlsx filename and spreadsheet name
	rows := ReadData("Benchmark.xlsx", "Sheet1")

	var bench []BenchMarkProprety

	// Add data into struct
	for i := range rows {

		// header location
		if i == 0 {
			continue
		}

		// Check first element is not empty to add
		if rows[i][0] != "" {
			bench = append(bench, BenchMarkProprety{
				Date:      rows[i][0],
				Isin:      rows[i][1],
				Ric:       rows[i][2],
				IsoCty:    rows[i][3],
				IDXMktCap: rows[i][4],
				Gicses:    rows[i][22],
			})
		}
	}

	return bench

}

// CalStock function to calculate sum&presentage of stock in order to display
func CalStock(s []StockProprety) []StockCalculation {
	// Set variable
	var (
		stockPrecent []StockCalculation
		// Total sum of float market values
		totalSum float64
		// GICS vars
		eSum  float64
		mSum  float64
		iSum  float64
		cdSum float64
		csSum float64
		hcSum float64
		fSum  float64
		itSum float64
		tsSum float64
		uSum  float64
		reSum float64
		// Region vars
		naSum     float64
		eurxukSum float64
		gbSum     float64
		apxjpSum  float64
		jpSum     float64
	)
	// Get every FloatMktCap from struct, convert them to float64, and sum up
	for i := range s {

		// Convert sting to float64
		number, err := strconv.ParseFloat(s[i].FloatMktCap, 64)
		if err != nil {
			break
		}
		// Round to integer
		number = math.Round(number)

		// Check Gics to define 11 classification
		// Energy 10
		if s[i].Gics[:2] == "10" {
			eSum += number
		}
		// Materials 15
		if s[i].Gics[:2] == "15" {
			mSum += number
		}
		// Industrials 20
		if s[i].Gics[:2] == "20" {
			iSum += number
		}
		// Consumer Discretionary 25
		if s[i].Gics[:2] == "25" {
			cdSum += number
		}
		// Consumer Staples 30
		if s[i].Gics[:2] == "30" {
			csSum += number
		}
		// Health Care 35
		if s[i].Gics[:2] == "35" {
			hcSum += number
		}
		// Financials 40
		if s[i].Gics[:2] == "40" {
			fSum += number
		}
		// Information Technology 45
		if s[i].Gics[:2] == "45" {
			itSum += number
		}
		// Telecommunication Services 50
		if s[i].Gics[:2] == "50" {
			tsSum += number
		}
		// Utilities 55
		if s[i].Gics[:2] == "55" {
			uSum += number
		}
		// Real Estate 60
		if s[i].Gics[:2] == "60" {
			reSum += number
		}

		// Check ISO city code to define region
		if s[i].IsoCty == "CA" || s[i].IsoCty == "US" ||
			s[i].IsoCty == "MX" {
			naSum += number
		}
		if s[i].IsoCty == "GB" {
			gbSum += number
		}
		if s[i].IsoCty == "JP" {
			jpSum += number
		}
		if s[i].IsoCty == "AU" || s[i].IsoCty == "HK" ||
			s[i].IsoCty == "NZ" || s[i].IsoCty == "SG" ||
			s[i].IsoCty == "CN" || s[i].IsoCty == "KR" ||
			s[i].IsoCty == "TW" {
			apxjpSum += number
		}
		if s[i].IsoCty == "AT" || s[i].IsoCty == "BE" ||
			s[i].IsoCty == "CH" || s[i].IsoCty == "DE" ||
			s[i].IsoCty == "DK" || s[i].IsoCty == "ES" ||
			s[i].IsoCty == "FI" || s[i].IsoCty == "FR" ||
			s[i].IsoCty == "IE" || s[i].IsoCty == "IL" ||
			s[i].IsoCty == "IT" || s[i].IsoCty == "NL" ||
			s[i].IsoCty == "NO" || s[i].IsoCty == "PT" ||
			s[i].IsoCty == "CZ" || s[i].IsoCty == "GR" ||
			s[i].IsoCty == "HU" || s[i].IsoCty == "PL" ||
			s[i].IsoCty == "SE" {
			eurxukSum += number
		}

		// Always add to totalsum
		totalSum += number
	}

	// Add these variables to slice for further use
	stockGICSName := []string{"Energy", "Materials", "Industrials", "Consumer Discretionary",
		"Consumer Staples", "Health Care", "Financials", "Information Technology",
		"Telecommunication Services", "Utilities", "Real Estate"}
	stockGICSValue := []float64{eSum, mSum, iSum, cdSum, csSum, hcSum, fSum,
		itSum, tsSum, uSum, reSum}

	stockRegionName := []string{"North America", "Europe ex UK", "United Kingdom",
		"Asia Pacific ex Japan", "Japan"}
	stockRegionValue := []float64{naSum, eurxukSum, gbSum, apxjpSum, jpSum}

	// get precentage
	gics := CalculatePresentage(totalSum, stockGICSValue)
	region := CalculatePresentage(totalSum, stockRegionValue)

	for i := 0; i < len(stockGICSName); i++ {
		stockPrecent = append(stockPrecent, StockCalculation{
			Name:       stockGICSName[i],
			Value:      stockGICSValue[i],
			Presentage: gics[i],
		})
	}

	for i := 0; i < len(stockRegionName); i++ {
		stockPrecent = append(stockPrecent, StockCalculation{
			Name:       stockRegionName[i],
			Value:      stockRegionValue[i],
			Presentage: region[i],
		})
		fmt.Println(stockRegionName[i], region[i])
	}

	return stockPrecent
}

// CalculatePresentage function use to get totalsum&values to calculate their precentage, then return
func CalculatePresentage(a float64, b []float64) []float64 {
	var presentage []float64
	var number float64
	for i := 0; i < len(b); i++ {

		// Float to 2 decimals
		n := (b[i] / a) * 100
		n = math.Round(n*100) / 100
		presentage = append(presentage, n)

		number += n
	}

	return presentage
}
