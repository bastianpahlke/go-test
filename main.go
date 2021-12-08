package main

import "github.com/gin-gonic/gin"
import "net/http"

func main() {
	r := gin.Default()

	r.GET("/userData/:userId", func(c *gin.Context) {
		userId := c.Param("userId")
		c.String(http.StatusOK, "Data for user %s", userId)
	})

	r.Run()
}
