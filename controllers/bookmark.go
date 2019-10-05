package controllers

import (
	jwt "github.com/appleboy/gin-jwt"
	dbpkg "github.com/filiponegrao/desafio-globo.com/db"
	"github.com/filiponegrao/desafio-globo.com/models"

	"github.com/gin-gonic/gin"
)

func CreateBookmark(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID := int64(claims["id"].(float64))

	db := dbpkg.DBInstance(c)
	bookmark := models.Bookmark{}

	if err := c.Bind(&bookmark); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	bookmark.UserID = userID
	if err := db.Find(&bookmark.Owner, userID).Error; err != nil {
		message := "Houve um problema com suas credenciais. Atuentique-se novamente."
		c.HTML(200, "bookmarks.html", gin.H{"message": message})
		return
	}

	if err := db.Create(&bookmark).Error; err != nil {
		message := err.Error()
		c.HTML(200, "bookmarks.html", gin.H{"message": message})
		return
	}

	ShowBookMarksPage(c, userID, "Criado com sucesso!")
	// GetBookmarksPage(c)
	// c.HTML(200, "bookmarks.html", nil)
}

func DeleteBookmark(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	userID := int64(claims["id"].(float64))

	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	bookmark := models.Bookmark{}

	if db.First(&bookmark, id).Error != nil {
		message := "URL não existe mais"
		c.HTML(200, "bookmarks.html", message)
		return
	}

	if bookmark.UserID != userID {
		message := "Você nãoé o dono desta URL"
		c.HTML(200, "bookmarks.html", message)
		return
	}

	if err := db.Delete(&bookmark).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ShowBookMarksPage(c, userID, "URL Excluída com sucesso.")
}
