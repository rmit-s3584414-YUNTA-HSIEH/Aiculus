// routes.go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var a = SetStockData()
var b = SetBMData()

func initializeRoutes() {

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

	router.GET("/line", func(c *gin.Context) {

		stock := a
		c.JSON(http.StatusOK, stock)
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
		stock := CalGICS(a, b)
		c.JSON(http.StatusOK, stock)
	})

	router.GET("/stock/:id", func(c *gin.Context) {
		id := c.Param("id")
		stock := FindCics(id)
		c.JSON(http.StatusOK, stock)
	})

	router.GET("/bar/:id", func(c *gin.Context) {
		id := c.Param("id")
		address := "bar.html"
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

}
