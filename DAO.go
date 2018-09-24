//DAO.go

package main

import (
	"fmt"
	"log"
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
		Isin        string `json:"isin"`
		Ric         string `json:"ric"`
		Name        string `json:"name"`
		IsoCty      string `json:"isocty"`
		Gics        string `json:"gics"`
		FloatMktCap string `json:"float"`
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
		Code       string `json:"code"`
		Name       string `json:"name"`
		Value      float64
		Presentage float64 `json:"presentage"`
	}

	// BenchMarkCalculation struct to store all benchmark calculation tables
	BenchMarkCalculation struct {
		Code       string
		Name       string
		Value      float64
		Presentage float64
	}

	// StockVMQ struct to store information of VMQ score
	StockVMQ struct {
		Name     string
		VScore   float64
		MScore   float64
		QScore   float64
		VMQScore float64
	}
)

// ReadData from xlsx
func ReadData(xlsxName string, sheetName string, header string) [][]string {

	xlsx, err := excelize.OpenFile(xlsxName)
	if err != nil {
		fmt.Println("Cannot find the file")
	}

	// check where's header
	var a = false

	for a != true {
		if xlsx.GetCellValue(sheetName, "A1") != header {
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
	rows := ReadData("Summary Data.xlsx", "Universe", "CALC_DATE")

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
	rows := ReadData("Benchmark.xlsx", "Sheet1", "CALC_DATE")

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

// SetVMQScore function to calculate VMQ score
func SetVMQScore() []StockVMQ {
	// Read data from excel, pass xlsx filename and spreadsheet name
	rows := ReadData("Summary Data.xlsx", "VMQ Scores", "ISIN")

	var vmq []StockVMQ

	// Add data into struct
	for i := range rows {

		// header location
		if i == 0 {
			continue
		}

		// Check first element is not empty to add
		if rows[i][0] != "" {
			v := StringToFloat(rows[i][21])
			m := StringToFloat(rows[i][22])
			q := StringToFloat(rows[i][23])
			vmq = append(vmq, StockVMQ{
				Name:     rows[i][2],
				VScore:   v,
				MScore:   m,
				QScore:   q,
				VMQScore: v + m + q,
			})
		}
	}

	return vmq

}

// CalStock function to calculate sum&presentage of stock in order to display
func CalStock(s []StockProprety) []StockCalculation {
	// Set variable
	var (
		stockPrecent []StockCalculation
		// Total sum of float market values
		totalSum float64
	)
	// Get Region map
	regionMap := BuildRegionMap()

	// Set GICS struct
	stockGICSCode := []string{"10", "15", "20", "25", "30", "35",
		"40", "45", "50", "55", "60"}
	stockGICSName := []string{"Energy", "Materials", "Industrials", "Consumer Discretionary",
		"Consumer Staples", "Health Care", "Financials", "Information Technology",
		"Telecommunication Services", "Utilities", "Real Estate"}

	for i := 0; i < len(stockGICSName); i++ {
		stockPrecent = append(stockPrecent, StockCalculation{
			Code:       stockGICSCode[i],
			Name:       stockGICSName[i],
			Value:      0,
			Presentage: 0,
		})
	}

	// Set Region struct
	stockRegionCode := []string{"NA", "EURXUK", "GB", "APXJP", "JP"}
	stockRegionName := []string{"North America", "Europe ex UK", "United Kingdom",
		"Asia Pacific ex Japan", "Japan"}

	for i := 0; i < len(stockRegionName); i++ {
		stockPrecent = append(stockPrecent, StockCalculation{
			Code:       stockRegionCode[i],
			Name:       stockRegionName[i],
			Value:      0,
			Presentage: 0,
		})
	}

	// Get every FloatMktCap from struct, convert them to float64, and sum up
	for i := range s {

		number := StringToFloat(s[i].FloatMktCap)
		code := s[i].Gics[:2]
		region := s[i].IsoCty
		regionCode := CheckRegion(region, regionMap)

		// Calculate value base by GICS
		for j := 0; j < len(stockGICSCode); j++ {
			stockPrecent[j].SetValue(code, number)
		}

		// Calculate value base by Region
		for k := len(stockGICSCode); k < (len(stockGICSCode) + len(stockRegionCode)); k++ {
			stockPrecent[k].SetValue(regionCode, number)
		}

		// Always add to totalsum
		totalSum += number
	}

	// Set presentage of each item
	for i := 0; i < (len(stockGICSCode) + len(stockRegionCode)); i++ {
		stockPrecent[i].SetPresentage(totalSum)
	}

	return stockPrecent
}

// CalBench function to calculate sum&presentage of benchmark in order to display
func CalBench(b []BenchMarkProprety) []BenchMarkCalculation {
	// Set variable
	var (
		bench []BenchMarkCalculation
		// Total sum of float market values
		totalSum float64
	)
	// Get Region map
	regionMap := BuildRegionMap()

	// Set GICS struct
	benchGICSCode, benchGICSName := BuildGICSList()

	for i := 0; i < len(benchGICSName); i++ {
		bench = append(bench, BenchMarkCalculation{
			Code:       benchGICSCode[i],
			Name:       benchGICSName[i],
			Value:      0,
			Presentage: 0,
		})
	}

	// Set Region struct
	benchRegionCode, benchRegionName := BuildRegionList()

	for i := 0; i < len(benchRegionName); i++ {
		bench = append(bench, BenchMarkCalculation{
			Code:       benchRegionCode[i],
			Name:       benchRegionName[i],
			Value:      0,
			Presentage: 0,
		})
	}

	// Get every FloatMktCap from struct, convert them to float64, and sum up
	for i := range b {

		number := StringToFloat(b[i].IDXMktCap)
		code := b[i].Gicses[:2]
		region := b[i].IsoCty
		regionCode := CheckRegion(region, regionMap)

		// Calculate value base by GICS
		for j := 0; j < len(benchGICSCode); j++ {
			bench[j].SetBValue(code, number)
		}

		// Calculate value base by Region
		for k := len(benchGICSCode); k < (len(benchGICSCode) + len(benchRegionCode)); k++ {
			bench[k].SetBValue(regionCode, number)
		}

		// Always add to totalsum
		totalSum += number
	}

	// Set presentage of each item
	for i := 0; i < (len(benchGICSCode) + len(benchRegionCode)); i++ {
		bench[i].SetBPresentage(totalSum)
	}

	return bench
}

// BuildGICSList function to provide the list of GICS name&code
func BuildGICSList() ([]string, []string) {

	a := []string{"10", "15", "20", "25", "30", "35",
		"40", "45", "50", "55", "60"}
	b := []string{"Energy", "Materials", "Industrials", "Consumer Discretionary",
		"Consumer Staples", "Health Care", "Financials", "Information Technology",
		"Telecommunication Services", "Utilities", "Real Estate"}

	return a, b

}

// BuildRegionList function to provide the list of region name&code
func BuildRegionList() ([]string, []string) {

	a := []string{"NA", "EURXUK", "GB", "APXJP", "JP"}
	b := []string{"North America", "Europe ex UK", "United Kingdom",
		"Asia Pacific ex Japan", "Japan"}

	return a, b

}

// BuildRegionMap function to build the region map, which key is region and values is code
func BuildRegionMap() map[string]([]string) {

	// Make map
	var m map[string]([]string)
	m = make(map[string]([]string))
	// Current iso code we are using
	m["NA"] = []string{"CA", "US", "MX", "NA"}
	m["GB"] = []string{"GB"}
	m["JP"] = []string{"JP"}
	m["APXJP"] = []string{"AU", "HK", "NZ", "SG", "CN", "KR", "TW", "APXJP"}
	m["EURXUK"] = []string{"AT", "BE", "CH", "DE", "DK", "ES", "FI",
		"FR", "IE", "IL", "IT", "NL", "NO", "PT", "CZ", "GR", "HU",
		"PL", "SE", "EURXUK"}

	return m
}

// CheckRegion function to check region code
func CheckRegion(s string, m map[string]([]string)) string {

	// Check if value exist in map
	for key, value := range m {
		for i := range value {
			if value[i] == s {
				return key
			}
		}
	}

	return ""
}

// SetValue function to add value by getting correct code
func (c *StockCalculation) SetValue(s string, a float64) {
	if s == c.Code {
		c.Value = c.Value + a
	}
}

// SetPresentage function use to set presentage of each data struct
func (c *StockCalculation) SetPresentage(a float64) {
	c.Presentage = (c.Value / a) * 100
	c.Presentage = math.Round(c.Presentage*100) / 100
}

// SetBValue function to add value by getting correct code
func (c *BenchMarkCalculation) SetBValue(s string, a float64) {
	if s == c.Code {
		c.Value = c.Value + a
	}
}

// SetBPresentage function use to set presentage of each data struct
func (c *BenchMarkCalculation) SetBPresentage(a float64) {
	c.Presentage = (c.Value / a) * 100
	c.Presentage = math.Round(c.Presentage*100) / 100
}

// StringToFloat function to convert string to float for further use
func StringToFloat(s string) float64 {
	var n float64
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal(err)
	}
	n = math.Round(n*1000) / 1000
	return n
}
