// routes.go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var a = SetStockData()
var b = SetBMData()

func initializeRoutes() {
	//pages

	// Handle the index route
	router.GET("/", func(c *gin.Context) {
		stock := a
		// Call the HTML method of the Context to render a template
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"index.html",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{"text": stock[0].Date},
		)
	})

	router.GET("/bar/:id", func(c *gin.Context) {
		id := c.Param("id")
		address := "bar" + id + ".html"

		c.HTML(
			http.StatusOK,
			address,
			gin.H{"id": id},
		)
	})

	router.GET("/details/:id", func(c *gin.Context) {
		id := c.Param("id")
		address := "details.html"
		c.HTML(
			http.StatusOK,
			address,
			gin.H{"id": id},
		)
	})

	//json api
	router.GET("/line", func(c *gin.Context) {

		stock := a
		c.JSON(http.StatusOK, stock)
	})

	router.GET("/VMQ", func(c *gin.Context) {

		VMQ := SetVMQScore()
		c.JSON(http.StatusOK, VMQ)
	})

	router.GET("/barChart1", func(c *gin.Context) {
		stock := CalGICS(a, b)
		c.JSON(http.StatusOK, stock)
	})
	router.GET("/barChart2", func(c *gin.Context) {
		stock := CalRegion(a, b)
		c.JSON(http.StatusOK, stock)
	})
	router.GET("/barChart3", func(c *gin.Context) {
		stock := CalCountry(a, b)
		c.JSON(http.StatusOK, stock)
	})

	router.GET("/stock/:id", func(c *gin.Context) {
		id := c.Param("id")
		stock := FindID(id, a)
		c.JSON(http.StatusOK, stock)
	})

}
