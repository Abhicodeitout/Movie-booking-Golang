package handlers

import "github.com/gin-gonic/gin"

func Root(c *gin.Context) {
	c.JSON(200, gin.H{
		"service": "movie-booking-api",
		"status":  "ok",
		"version": "dev",
	})
}
