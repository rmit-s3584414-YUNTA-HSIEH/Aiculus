// routes.go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var a = SetStockData()

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

}
func getLineData() {
	router.GET("/line", func(c *gin.Context) {

		stock := a
		c.JSON(http.StatusOK, stock)
	})
}
func getOtherGraphs() {

	router.GET("/OtherGraphs", func(c *gin.Context) {
		stock := CalStock(a)
		c.JSON(http.StatusOK, stock)
	})
}
