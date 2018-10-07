package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// Type lists
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

	// GICSCalculation struct to store all stock calculation tables
	GICSCalculation struct {
		Code        string  `json:"code"`
		Name        string  `json:"name"`
		SValue      float64 `json:"svalue"`
		BValue      float64 `json:"bvalue"`
		SPercentage float64 `json:"spercentage"`
		BPercentage float64 `json:"bpercentage"`
		Diff        float64 `json:"diff"`
	}

	// RegionCalculation struct to store all benchmark calculation tables
	RegionCalculation struct {
		Code        string  `json:"code"`
		Name        string  `json:"name"`
		SValue      float64 `json:"svalue"`
		BValue      float64 `json:"bvalue"`
		SPercentage float64 `json:"spercentage"`
		BPercentage float64 `json:"bpercentage"`
		Diff        float64 `json:"diff"`
	}

	//CountryCalculation struct to store all calculation based by country
	CountryCalculation struct {
		Code        string  `json:"code"`
		Name        string  `json:"name"`
		RegionCode  string  `json:"region"`
		SValue      float64 `json:"svalue"`
		BValue      float64 `json:"bvalue"`
		SPercentage float64 `json:"spercentage"`
		BPercentage float64 `json:"bpercentage"`
		Diff        float64 `json:"diff"`
	}

	// StockVMQ struct to store information of VMQ score
	StockVMQ struct {
		Name     string  `json:"name"`
		VScore   float64 `json:"v"`
		MScore   float64 `json:"m"`
		QScore   float64 `json:"q"`
		VMQScore float64 `json:"vmq"`
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

//FindCics :get stock by gics
func FindCics(id string) []StockProprety {
	stock := SetStockData()
	var theGicsStock []StockProprety

	for i := 0; i < len(stock); i++ {
		var stringID = stock[i].Gics
		if strings.HasPrefix(stringID, id) {
			theGicsStock = append(theGicsStock, stock[i])
		}
	}
	if len(theGicsStock) == 0 {
		return stock
	}
	//fmt.Println(theGicsStock)
	return theGicsStock
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

// CalGICS function to calculate sum&Persentage of stock in order to display
func CalGICS(s []StockProprety, b []BenchMarkProprety) []GICSCalculation {
	// Set variable
	var (
		gics []GICSCalculation
		// Total sum of float market values
		totalSSum float64
		totalBSum float64
	)

	// Set GICS struct
	stockGICSCode, stockGICSName := BuildGICSList()

	for i := 0; i < len(stockGICSName); i++ {
		gics = append(gics, GICSCalculation{
			Code:        stockGICSCode[i],
			Name:        stockGICSName[i],
			SValue:      0,
			BValue:      0,
			SPercentage: 0,
			BPercentage: 0,
			Diff:        0,
		})
	}

	// Get every FloatMktCap from struct, convert them to float64, and sum up
	for i := range s {

		number := StringToFloat(s[i].FloatMktCap)
		code := s[i].Gics[:2]

		// Calculate value base by GICS
		for j := 0; j < len(stockGICSCode); j++ {
			gics[j].SetValue(code, number)
		}

		// Always add to totalsum
		totalSSum += number
	}

	// Get every FloatMktCap from struct, convert them to float64, and sum up
	for i := range b {
		number := StringToFloat(b[i].IDXMktCap)
		code := b[i].Gicses[:2]

		// Calculate value base by GICS
		for j := 0; j < len(stockGICSCode); j++ {
			gics[j].SetBValue(code, number)
		}

		// Always add to totalsum
		totalBSum += number
	}

	// Set Persentage of each item
	for i := 0; i < (len(stockGICSCode)); i++ {
		gics[i].SetPercentage(totalSSum)
		gics[i].SetBPercentage(totalBSum)
	}

	//fmt.Println(gics)
	return gics
}

// CalRegion function to calculate sum&Persentage of benchmark in order to display
func CalRegion(s []StockProprety, b []BenchMarkProprety) []RegionCalculation {
	// Set variable
	var (
		regions []RegionCalculation
		// Total sum of float market values
		totalSSum float64
		totalBSum float64
	)
	// Get Region map
	regionMap := BuildRegionMap()

	// Set Region struct
	benchRegionCode, benchRegionName := BuildRegionList()

	for i := 0; i < len(benchRegionName); i++ {
		regions = append(regions, RegionCalculation{
			Code:        benchRegionCode[i],
			Name:        benchRegionName[i],
			SValue:      0,
			BValue:      0,
			SPercentage: 0,
			BPercentage: 0,
			Diff:        0,
		})
	}

	// Get every FloatMktCap from struct, convert them to float64, and sum up
	for i := range s {

		number := StringToFloat(s[i].FloatMktCap)
		region := s[i].IsoCty
		regionCode := CheckRegion(region, regionMap)

		// Calculate value base by GICS
		for j := 0; j < len(benchRegionCode); j++ {
			regions[j].SetValue(regionCode, number)
		}

		// Always add to totalsum
		totalSSum += number
	}

	// Get every FloatMktCap from struct, convert them to float64, and sum up
	for i := range b {

		number := StringToFloat(b[i].IDXMktCap)
		region := b[i].IsoCty
		regionCode := CheckRegion(region, regionMap)

		// Calculate value base by Region
		for k := 0; k < len(benchRegionCode); k++ {
			regions[k].SetBValue(regionCode, number)
		}

		// Always add to totalsum
		totalBSum += number
	}

	// Set Persentage of each item
	for i := 0; i < len(benchRegionCode); i++ {
		regions[i].SetPercentage(totalSSum)
		regions[i].SetBPercentage(totalBSum)
	}

	//fmt.Println(regions)
	return regions
}

//CalCountry data
func CalCountry(s []StockProprety, b []BenchMarkProprety) []CountryCalculation {
	// Set variable
	var (
		countrys []CountryCalculation
		// Total sum of float market values
		totalSSum float64
		totalBSum float64
	)

	// Set Region struct
	stockCountryCode, stockCountryName := BuildCountryList()

	for i := 0; i < len(stockCountryName); i++ {
		countrys = append(countrys, CountryCalculation{
			Code:        stockCountryCode[i],
			Name:        stockCountryName[i],
			RegionCode:  "",
			SValue:      0,
			BValue:      0,
			SPercentage: 0,
			BPercentage: 0,
			Diff:        0,
		})
	}

	// Get every FloatMktCap from struct, convert them to float64, and sum up
	for i := range s {

		number := StringToFloat(s[i].FloatMktCap)
		country := s[i].IsoCty

		// Calculate value base by GICS
		for j := 0; j < len(stockCountryCode); j++ {
			countrys[j].SetSValue(country, number)
		}

		// Always add to totalsum
		totalSSum += number
	}

	// Get every FloatMktCap from struct, convert them to float64, and sum up
	for i := range b {

		number := StringToFloat(b[i].IDXMktCap)
		country := b[i].IsoCty

		// Calculate value base by Region
		for k := 0; k < len(stockCountryCode); k++ {
			countrys[k].SetBValue(country, number)
		}

		// Always add to totalsum
		totalBSum += number
	}

	// Set Persentage of each item
	for i := 0; i < (len(stockCountryCode)); i++ {
		countrys[i].SetPercentage(totalSSum)
		countrys[i].SetBPercentage(totalBSum)
	}

	for i := range countrys {
		countrys[i].SetRegionCode()
	}

	// fmt.Println(countrys)
	return countrys
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

// BuildCountryList function to build the Country map, which key is Country and values is code
func BuildCountryList() ([]string, []string) {

	a := []string{"US", "SG", "SE", "PT", "NZ",
		"NO", "NL", "MX", "KR", "JP",
		"IT", "IL", "IE", "HK", "GB",
		"FR", "FI", "ES", "DK", "DE",
		"CN", "CH", "CA", "BE", "AU", "AT"}
	b := []string{"United States", "Singapore", "Sweden", "Portugal", "New Zealand",
		"Norway", "Netherlands", "Mexico", "South Korea", "Japan",
		"Italy", "Israel", "Ireland", "Hong Kong", "Great Britain",
		"France", "Finland", "Spain", "Denmark", "Germany",
		"China", "Switzerland", "Canada", "Belgium", "Australia", "Austria"}

	return a, b

}

//CheckCountry function is to
func CheckCountry(s string, cm map[string]([]string)) string {

	// Check if value exist in map
	for key, value := range cm {
		for i := range value {
			if value[i] == s {
				return key
			}
		}
	}

	return ""
}

// GICSCalculation pointer functions below

// SetValue function to add value by getting correct code
func (c *GICSCalculation) SetValue(s string, a float64) {
	if s == c.Code {
		c.SValue = c.SValue + a
	}
}

// SetBValue function to
func (c *GICSCalculation) SetBValue(s string, a float64) {
	if s == c.Code {
		c.BValue = c.BValue + a
	}
}

// SetPercentage function use to set Persentage of each data struct
func (c *GICSCalculation) SetPercentage(a float64) {
	c.SPercentage = (c.SValue / a) * 100
	c.SPercentage = math.Round(c.SPercentage*100) / 100
}

// SetBPercentage function to
func (c *GICSCalculation) SetBPercentage(a float64) {
	c.BPercentage = (c.BValue / a) * 100
	c.BPercentage = math.Round(c.BPercentage*100) / 100
	c.Diff = math.Round((c.SPercentage-c.BPercentage)*1000) / 1000
}

// GICSCalculation pointer functions end

// RegionCalculation pointer functions below

// SetValue function to add value by getting correct code
func (c *RegionCalculation) SetValue(s string, a float64) {
	if s == c.Code {
		c.SValue = c.SValue + a
	}
}

// SetBValue function to add value by getting correct code
func (c *RegionCalculation) SetBValue(s string, a float64) {
	if s == c.Code {
		c.BValue = c.BValue + a
	}
}

// SetPercentage function use to set Persentage of each data struct
func (c *RegionCalculation) SetPercentage(a float64) {
	c.SPercentage = (c.SValue / a) * 100
	c.SPercentage = math.Round(c.SPercentage*100) / 100
}

// SetBPercentage function use to set Persentage of each data struct
func (c *RegionCalculation) SetBPercentage(a float64) {
	c.BPercentage = (c.BValue / a) * 100
	c.BPercentage = math.Round(c.BPercentage*100) / 100
	c.Diff = math.Round((c.SPercentage-c.BPercentage)*1000) / 1000
}

// RegionCalculation pointer functions end

// CountryCalculation pointer functions below

// SetSValue function to add value by getting correct code
func (c *CountryCalculation) SetSValue(s string, a float64) {
	if s == c.Code {
		c.SValue = c.SValue + a
	}
}

// SetBValue function to add value by getting correct code
func (c *CountryCalculation) SetBValue(s string, a float64) {
	if s == c.Code {
		c.BValue = c.BValue + a
	}
}

// SetPercentage function use to set Persentage of each data struct
func (c *CountryCalculation) SetPercentage(a float64) {
	c.SPercentage = (c.SValue / a) * 100
	c.SPercentage = math.Round(c.SPercentage*100) / 100
}

// SetBPercentage function use to set Persentage of each data struct
func (c *CountryCalculation) SetBPercentage(a float64) {
	c.BPercentage = (c.BValue / a) * 100
	c.BPercentage = math.Round(c.BPercentage*100) / 100
	c.Diff = math.Round((c.SPercentage-c.BPercentage)*1000) / 1000
}

// SetRegionCode function use to set RegionCode for each data struct
func (c *CountryCalculation) SetRegionCode() {
	regionMap := BuildRegionMap()
	c.RegionCode = CheckRegion(c.Code, regionMap)
}

// CountryCalculation pointer functions end

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
