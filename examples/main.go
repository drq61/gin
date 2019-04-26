package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()


	router.GET("/welcome", func(c *gin.Context) {

		if err := c.Request.Check("id","tag_id");err != nil {
			c.Error(err)
			return
		}

		c.String(http.StatusOK, "Hello %v",c.Request.GetInt64("tag_id"))
	})
	router.Run(":8080")
}