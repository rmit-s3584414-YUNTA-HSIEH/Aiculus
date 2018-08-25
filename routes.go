// routes.go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func initializeRoutes() {

	// Handle the index route
	router.GET("/", func(c *gin.Context) {
		stock := readExcel()
		// Call the HTML method of the Context to render a template
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"index.html",
			// Pass the data that the page uses (in this case, 'title')

			gin.H{
				"message": stock,
			},
		)
	})

	router.GET("/json", func(c *gin.Context) {
		stock := readExcel()
		c.JSON(http.StatusOK, gin.H{"stock": stock.name, "price": stock.price})

	})
}
