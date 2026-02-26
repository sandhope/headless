package controller

import (
	"app/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleCookie(c *gin.Context) {
	var domain = c.Query("domain")
	var key = c.Query("key")

	if domain == "" {
		c.AbortWithStatus(400)
		return
	}

	result, err := service.GetCookie(domain)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if key != "" {
		c.String(http.StatusOK, result[key])
		return
	}

	c.JSON(http.StatusOK, result)
}
