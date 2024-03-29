package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIEndpoints(c *gin.Context) {
	reqScheme := "http"

	if c.Request.TLS != nil {
		reqScheme = "https"
	}

	reqHost := c.Request.Host
	baseURL := fmt.Sprintf("%s://%s", reqScheme, reqHost)

	resources := map[string]string{
		"bookmarks_url":         baseURL + "/bookmarks",
		"bookmark_url":          baseURL + "/bookmarks/{id}",
		"password_recovers_url": baseURL + "/password_recovers",
		"password_recover_url":  baseURL + "/password_recovers/{id}",
		"users_url":             baseURL + "/users",
		"user_url":              baseURL + "/users/{id}",
	}

	c.IndentedJSON(http.StatusOK, resources)
}
