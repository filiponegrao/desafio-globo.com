package controllers

import (
	"log"
	"strings"

	jwt "github.com/appleboy/gin-jwt"
	dbpkg "github.com/filiponegrao/desafio-globo.com/db"
	"github.com/filiponegrao/desafio-globo.com/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func GetLoginPage(c *gin.Context) {
	token := csrf.GetToken(c)
	c.HTML(200, "login.html", gin.H{"token": token})
}

func GetRegsisterPage(c *gin.Context) {
	token := csrf.GetToken(c)
	c.HTML(200, "register.html", gin.H{"token": token})
}

func GetForgotPasswordPage(c *gin.Context) {
	token := csrf.GetToken(c)

	c.HTML(200, "forgot-password.html", gin.H{"token": token})
}

func GetNewPasswordPage(c *gin.Context) {
	token := csrf.GetToken(c)
	c.HTML(200, "new-password.html", gin.H{"token": token})
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
	token := csrf.GetToken(c)

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
		"token":     token,
		"username":  user.Name,
		"bookmarks": bookmarks,
		"message":   m,
	})
}

func GetCreateBookmarkPage(c *gin.Context) {
	token := csrf.GetToken(c)

	c.HTML(200, "create-bookmark.html", gin.H{"token": token})
}

func CheckIncorrectInput(input string) bool {
	if strings.ContainsAny(input, "'\\-\"") {
		return true
	}

	return false
}
