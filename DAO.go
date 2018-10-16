package main

import (
	"fmt"
	"io"
	"math"
	"os"
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

	// SecruityData struct to store information of secruity data
	SecruityData struct {
		Name   string  `json:"name"`
		IsoCty string  `json:"isocty"`
		Sector string  `json:"sector"`
		Weight float64 `json:"weight"`
	}

	// StockVMQ struct to store information of VMQ score
	StockVMQ struct {
		Name     string `json:"name"`
		Date     []string
		VScore   []float64 `json:"v"`
		MScore   []float64 `json:"m"`
		QScore   []float64 `json:"q"`
		VMQScore []float64 `json:"vmq"`
	}

	// GICSCalculation struct to store all stock calculation tables
	GICSCalculation struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		SValue      float64
		BValue      float64
		SPercentage float64 `json:"spercentage"`
		BPercentage float64 `json:"bpercentage"`
		Diff        float64 `json:"diff"`
	}

	// RegionCalculation struct to store all benchmark calculation tables
	RegionCalculation struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		SValue      float64
		BValue      float64
		SPercentage float64 `json:"spercentage"`
		BPercentage float64 `json:"bpercentage"`
		Diff        float64 `json:"diff"`
	}

	//CountryCalculation struct to store all calculation based by country
	CountryCalculation struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		RegionCode  string `json:"region"`
		SValue      float64
		BValue      float64
		SPercentage float64 `json:"spercentage"`
		BPercentage float64 `json:"bpercentage"`
		Diff        float64 `json:"diff"`
	}
)

var sLog = "log/s_reports.log"
var bmLog = "log/bm_reports.log"

//
func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

//create file is not exists
func createFile(path string) {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if isError(err) {
			return
		}
		defer file.Close()
	}

	// fmt.Println("==> done creating file", path)
}

func writeLogFile(path string, rows []int, number int) {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) {
		return
	}
	defer file.Close()

	// write some text line-by-line to file
	if len(rows) > 0 {
		if rows[0] == 0 {
			_, err = file.WriteString("Colomuns missing.\r\n")
		}
		file.WriteString(strconv.Itoa(number) + " records are imported.\r\n")
		file.WriteString(strconv.Itoa(len(rows)) + " records are failed.\r\n")
		for i := range rows {
			_, err = file.WriteString("row" + " " + strconv.Itoa(rows[i]) + "is missing.\r\n")
		}
	} else {
		file.WriteString(strconv.Itoa(number) + " records are imported.\r\n")
		file.WriteString("No error.")
	}
	// save changes
	err = file.Sync()
	if isError(err) {
		return
	}

	// fmt.Println("==> done writing to file")
}

func readFile(path string) {
	// re-open file
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) {
		return
	}
	defer file.Close()

	// read file, line by line
	var text = make([]byte, 1024)
	for {
		_, err = file.Read(text)

		// break if finally arrived at end of file
		if err == io.EOF {
			break
		}

		// break if error occured
		if err != nil && err != io.EOF {
			isError(err)
			break
		}
	}

	fmt.Println("==> done reading from file")
	fmt.Println(string(text))
}

func deleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	if isError(err) {
		return
	}

	fmt.Println("==> done deleting file")
}

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

//validate excel file before loading
func validationSummary(rows [][]string) (bool, []int) {
	//city := []string{"CA", "US", "MX", "NA", "AU", "HK", "NZ",
	//	"SG", "CN", "KR", "TW", "APXJP", "AT", "BE", "CH", "DE",
	//	"DK", "ES", "FI", "FR", "IE", "IL", "IT", "NL", "NO", "PT",
	//	"CZ", "GR", "HU", "PL", "SE", "EURXUK", "GB", "JP"}
	//gics := []string{}
	errorRows := []int{}
	errorCols := false
	pointer := []int{}
	counter := 0
	success := 0

	//var errorLogs = ""
	//var correctRows int
	for j := range rows[0] {
		if rows[0][j] == "GICS_SI" {
			counter++
			pointer = append(pointer, j)
		}
		if rows[0][j] == "ISO_CTY_CODE" {
			counter++
			pointer = append(pointer, j)
		}
		if rows[0][j] == "FFLOAT_MKTCAP_USD" {
			counter++
			pointer = append(pointer, j)
		}

	}
	if counter != 3 {
		//table columns are missing
		errorRows = append(errorRows, 0)
		errorCols := true
		return errorCols, errorRows
	}
	for i := range rows {

		if i == 0 {
			continue
		}

		if rows[i][0] != "" {
			for j := range pointer {
				if len(rows[i][pointer[j]]) == 0 {
					errorRows = append(errorRows, i)
				}
			}
			success++
		}
	}
	createFile(sLog)
	writeLogFile(sLog, errorRows, success)
	return errorCols, errorRows
}

//Benchmark validation
func validationBenchmark(rows [][]string) (bool, []int) {
	errorRows := []int{}
	errorCols := false
	pointer := []int{}
	counter := 0
	success := 0

	for j := range rows[0] {
		if rows[0][j] == "GICS_IG" {
			counter++
			pointer = append(pointer, j)
		}
		if rows[0][j] == "GICS_ES" {
			counter++
			pointer = append(pointer, j)
		}
		if rows[0][j] == "ISO_CTY_CODE2" {
			counter++
			pointer = append(pointer, j)
		}
		if rows[0][j] == "REGION" {
			counter++
			pointer = append(pointer, j)
		}
	}
	if counter != 4 {
		//table columns are missing
		errorRows = append(errorRows, 0)
		errorCols = true
		return errorCols, errorRows
	}

	for i := range rows {

		if i == 0 {
			continue
		}
		if rows[i][0] != "" {
			for j := range pointer {
				if len(rows[i][pointer[j]]) == 0 {
					errorRows = append(errorRows, i)
				}
			}
			success++
		}

	}

	createFile(bmLog)
	writeLogFile(bmLog, errorRows, success)
	return errorCols, errorRows
}

// SetStockData function use to read data from excel and return the stock struct
func SetStockData() []StockProprety {

	// Read data from excel, pass xlsx filename and spreadsheet name
	rows := ReadData("data/Summary Data.xlsx", "Universe", "CALC_DATE")

	dataFail, errorRows := validationSummary(rows)
	fmt.Println(dataFail)
	var stock []StockProprety
	// Add data into struct
	for i := range rows {

		A := false
		// header location
		if i == 0 {
			continue
		}

		for j := range errorRows {
			if i == errorRows[j] {
				A = true
			}
		}
		if A == true {
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

//FindID :get id from table and return the data
func FindID(id string, stock []StockProprety) []StockProprety {

	var theStock []StockProprety

	// Check the type of id
	idtype := CheckID(id)

	if idtype == "gics" {
		for i := 0; i < len(stock); i++ {
			var stringID = stock[i].Gics
			if strings.HasPrefix(stringID, id) {
				theStock = append(theStock, stock[i])
			}
		}
		if len(theStock) == 0 {
			return stock
		}

		return theStock
	}

	if idtype == "country" {
		for i := 0; i < len(stock); i++ {
			var stringID = stock[i].IsoCty
			if strings.HasPrefix(stringID, id) {
				theStock = append(theStock, stock[i])
			}
		}
		if len(theStock) == 0 {
			return stock
		}

		return theStock
	}

	if idtype == "region" {

		// Get regionmap
		regionmap := BuildRegionMap()

		idlist := regionmap[id]

		// Check every isocty according to the map
		for j := range regionmap[id] {

			for i := 0; i < len(stock); i++ {
				var stringID = stock[i].IsoCty

				if strings.HasPrefix(stringID, idlist[j]) {
					theStock = append(theStock, stock[i])
				}
			}
		}
		if len(theStock) == 0 {
			return stock
		}
		return theStock
	}

	return nil
}

// CheckID function to check the type of return id
func CheckID(id string) string {

	gics := []string{"10", "15", "20", "25", "30", "35",
		"40", "45", "50", "55", "60"}

	region := []string{"NA", "EURXUK", "GB", "APXJP", "JP"}

	country := []string{"US", "SG", "SE", "PT", "NZ",
		"NO", "NL", "MX", "KR", "JP",
		"IT", "IL", "IE", "HK", "GB",
		"FR", "FI", "ES", "DK", "DE",
		"CN", "CH", "CA", "BE", "AU", "AT"}

	for i := range gics {
		if id == gics[i] {
			return "gics"
		}
	}
	for i := range region {
		if id == region[i] {
			return "region"
		}
	}
	for i := range country {
		if id == country[i] {
			return "country"
		}
	}
	return ""
}

// SetBMData function use to read data from excel and return the benchmark struct
func SetBMData() []BenchMarkProprety {

	// Read data from excel, pass xlsx filename and spreadsheet name
	rows := ReadData("data/Benchmark.xlsx", "Sheet1", "CALC_DATE")

	dataFail, errorRows := validationBenchmark(rows)
	fmt.Println(dataFail)
	var bench []BenchMarkProprety

	// Add data into struct
	for i := range rows {

		A := false
		// header location
		if i == 0 {
			continue
		}

		for j := range errorRows {
			if i == errorRows[j] {
				A = true
			}
		}
		if A == true {
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

// SetSecruityData function use to read data from excel and return the secruity struct
func SetSecruityData() []SecruityData {

	// Read data from excel, pass xlsx filename and spreadsheet name
	rows := ReadData("data/Summary Data.xlsx", "Portfolio", "prev_date_d0")

	var (
		secruity []SecruityData
		header   []int
	)

	for i := range rows[0] {
		if rows[0][i] == "iso_cty" {
			header = append(header, i)
		}
		if rows[0][i] == "gics_ind" {
			header = append(header, i)
		}
		if rows[0][i] == "name" {
			header = append(header, i)
		}
		if rows[0][i] == "stock_only_wgt" {
			header = append(header, i)
		}
		if rows[0][i] == "GICS_ES" {
			header = append(header, i)
		}
	}
	// Add data into struct
	for i := range rows {

		// header location
		if i == 0 {
			continue
		}

		// Get full gics code
		s := []string{rows[i][header[4]], rows[i][header[1]]}
		gics := strings.Join(s, "")
		sector := GetGICSName(gics)

		// Convert % number into float64
		weightS := rows[i][header[3]]
		weightS = strings.TrimRight(weightS, "%")
		weight := StringToFloat(weightS)

		// Check isocty code
		if len(rows[i][header[0]]) == 2 {
			secruity = append(secruity, SecruityData{
				Name:   rows[i][header[2]],
				IsoCty: rows[i][header[0]],
				Sector: sector,
				Weight: weight,
			})
		}
	}

	return secruity

}

// GetGICSName function is to define the gics name by code
func GetGICSName(s string) string {

	// Set GICS struct
	stockGICSCode := []string{"10", "15", "20", "25", "30", "35",
		"40", "45", "50", "55", "60"}
	stockGICSName := []string{"ENE", "MAT", "IND", "CSD", "CSS",
		"HLC", "FIN", "IFT", "TEL", "UTI", "REL"}
	for i := range stockGICSCode {
		if s[:2] == stockGICSCode[i] {
			return stockGICSName[i]
		}
	}
	return ""
}

// SetVMQScore function to calculate VMQ score
func SetVMQScore() []StockVMQ {
	// Read data from excel, pass xlsx filename and spreadsheet name
	rows := ReadData("data/Book1.xlsx", "VMQ Scores", "DATE")

	var (
		vmq     []StockVMQ
		date    []string
		pointer []int
	)

	// Define company name
	for i := range rows[1] {
		if rows[1][i] == "SECURITY_NAME" {
			var A bool
			for j := 1; j < 4; j++ {
				if rows[1][i+j] == "" {
					A = false
				} else {
					A = true
				}
			}
			if A != false {
				pointer = append(pointer, i)
				date = append(date, rows[0][i+4])
			}
		}
	}

	// // Define date
	// for i := 1; i < len(rows[0]); i++ {
	// 	if rows[0][i] != "" {
	// 		date = append(date, rows[0][i])
	// 	}
	// }

	// Setup struct
	for i := 2; i < len(rows); i++ {
		if rows[i][0] != "" {
			vmq = append(vmq, StockVMQ{
				Name:     rows[i][0],
				Date:     date,
				VScore:   nil,
				MScore:   nil,
				QScore:   nil,
				VMQScore: nil,
			})
		}
	}

	// Add data into struct
	for i := range rows {

		// Skip first two rows
		if i > 1 {

			// Check the name pointer
			for j := range pointer {

				// Ensure insert non-nil data

				if rows[i][pointer[j]] != "" {
					name := rows[i][pointer[j]]
					v := StringToFloat(rows[i][pointer[j]+1])
					m := StringToFloat(rows[i][pointer[j]+2])
					q := StringToFloat(rows[i][pointer[j]+3])
					vmqs := StringToFloat(rows[i][pointer[j]+4])

					// Check the struct to add data
					for i := 0; i < len(vmq); i++ {
						vmq[i].SetVMQ(name, v, m, q, vmqs)
					}
				}
			}
		}
	}
	if len(pointer) > 0 {
		for i := 0; i < len(vmq)-1; i++ {
			for j := 0; j < len(vmq)-1-i; j++ {
				if vmq[j].VMQScore[0] < vmq[j+1].VMQScore[0] {
					temp := vmq[j]
					vmq[j] = vmq[j+1]
					vmq[j+1] = temp
				}

			}

		}
	}

	return vmq

}

// SetVMQ function to set vmq score, base on the name
func (c *StockVMQ) SetVMQ(s string, v float64, m float64, q float64, vmqs float64) {
	if s == c.Name {
		c.VScore = append(c.VScore, v)
		c.MScore = append(c.MScore, m)
		c.QScore = append(c.QScore, q)
		c.VMQScore = append(c.VMQScore, vmqs)
	}
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
	if err == nil {
		n = math.Round(n*1000) / 1000
	}
	return n
}
