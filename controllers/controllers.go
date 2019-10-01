package controllers

import (
	"log"

	jwt "github.com/appleboy/gin-jwt"
	dbpkg "github.com/filiponegrao/desafio-globo.com/db"
	"github.com/filiponegrao/desafio-globo.com/models"
	"github.com/gin-gonic/gin"
)

func GetLoginPage(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

func GetRegsisterPage(c *gin.Context) {
	c.HTML(200, "register.html", nil)
}

func GetBookmarksPage(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID := int64(claims["id"].(float64))
	db := dbpkg.DBInstance(c)

	message := ""
	var bookmarks []models.Bookmark
	if err := db.Where("user_id = ?", userID).Find(&bookmarks).Error; err != nil {
		message = "Erro ao tentar recuperar suas URLs: " + err.Error()
	}

	ShowBookMarksPage(c, userID, message)
}

func ShowBookMarksPage(c *gin.Context, userId int64, message string) {
	db := dbpkg.DBInstance(c)

	var user models.User
	if err := db.Find(&user, userId).Error; err != nil {
		log.Println(err)
	}
	m := message
	var bookmarks []models.Bookmark
	if err := db.Where("user_id = ?", userId).Find(&bookmarks).Error; err != nil {
		m = "Erro ao tentar recuperar suas URLs: " + err.Error()
	}

	c.HTML(200, "bookmarks.html", gin.H{
		"username":  user.Name,
		"bookmarks": bookmarks,
		"message":   m,
	})
}

func GetCreateBookmarkPage(c *gin.Context) {

	// claims := jwt.ExtractClaims(c)
	// userID := int64(claims["id"].(float64))

	c.HTML(200, "create-bookmark.html", nil)
}
