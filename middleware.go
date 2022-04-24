package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		api_key := c.Query("api_key")
		if api_key == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "please provide api_key"})
			c.Abort()
			return
		}
		c.Next()
	}
}
