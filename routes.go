// routes.go

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

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

	router.GET("/sectorDetail", func(c *gin.Context) {

		address := "sectorDetail.html"
		c.HTML(
			http.StatusOK,
			address,
			gin.H{},
		)
	})

	router.GET("/regionDetail", func(c *gin.Context) {

		address := "regionDetail.html"
		c.HTML(
			http.StatusOK,
			address,
			gin.H{},
		)
	})

	router.GET("/support", func(c *gin.Context) {

		address := "support.html"
		c.HTML(
			http.StatusOK,
			address,
			gin.H{},
		)
	})

	router.GET("/about", func(c *gin.Context) {

		address := "about.html"
		c.HTML(
			http.StatusOK,
			address,
			gin.H{},
		)
	})

	router.GET("/security", func(c *gin.Context) {

		address := "security.html"
		c.HTML(
			http.StatusOK,
			address,
			gin.H{},
		)
	})

	router.GET("/country", func(c *gin.Context) {

		address := "country.html"
		c.HTML(
			http.StatusOK,
			address,
			gin.H{},
		)
	})

	router.GET("/detailsVMQ", func(c *gin.Context) {

		address := "detailsVMQ.html"
		c.HTML(
			http.StatusOK,
			address,
			gin.H{},
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

	router.GET("/account", func(c *gin.Context) {
		address := "account.html"
		c.HTML(
			http.StatusOK,
			address,
			gin.H{},
		)
	})

	router.GET("/testVMQ", func(c *gin.Context) {
		address := "testVMQTable.html"
		c.HTML(
			http.StatusOK,
			address,
			gin.H{},
		)
	})

	router.POST("/upload", func(c *gin.Context) {
		// name := c.PostForm("name")
		// email := c.PostForm("email")

		// Source
		file, err := c.FormFile("file")
		if err != nil {
			//c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			//return
		}

		// if err := c.SaveUploadedFile(file, file.Filename); err != nil {
		//c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		//return
		// }
		src, _ := file.Open()
		defer src.Close()
		// Destination
		dstname := fmt.Sprintf("./data/%s", file.Filename)
		dst, _ := os.Create(dstname)
		defer dst.Close()
		// Copy
		io.Copy(dst, src)

		c.Redirect(http.StatusMovedPermanently, "/")

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

	router.GET("/sector", func(c *gin.Context) {
		stock := CalGICS(a, b)
		c.JSON(http.StatusOK, stock)
	})
	router.GET("/region", func(c *gin.Context) {
		stock := CalRegion(a, b)
		c.JSON(http.StatusOK, stock)
	})
	router.GET("/countries", func(c *gin.Context) {
		stock := CalCountry(a, b)
		c.JSON(http.StatusOK, stock)
	})

	router.GET("/stock/:id", func(c *gin.Context) {
		id := c.Param("id")
		stock := FindID(id, a)
		c.JSON(http.StatusOK, stock)
	})

	router.GET("/secruity", func(c *gin.Context) {
		sec := SetSecruityData()
		c.JSON(http.StatusOK, sec)
	})

}
