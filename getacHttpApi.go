package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getacFota(c *gin.Context) {
	log.Println("RunFOTA")
	file, err := c.FormFile("image_file")
	if err != nil {
		// Handle error
		return
	}

	// Save the uploaded file
	err = c.SaveUploadedFile(file, "/tmp/fota.img")
	if err != nil {
		// Handle error
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})

	log.Println("getacFota done")
	//c.HTML(http.StatusOK, "fota.page.gohtml", nil)
}
