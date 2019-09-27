package controllers

import (
	"github.com/gin-gonic/gin"
)

func GetLoginPage(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

func GetRegsisterPage(c *gin.Context) {
	c.HTML(200, "register.html", nil)
}

func GetBookmarksPage(c *gin.Context) {
	c.HTML(200, "bookmarks.html", nil)
}
