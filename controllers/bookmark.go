package controllers

import (
	"encoding/json"
	"log"

	jwt "github.com/appleboy/gin-jwt"
	dbpkg "github.com/filiponegrao/desafio-globo.com/db"
	"github.com/filiponegrao/desafio-globo.com/helper"
	"github.com/filiponegrao/desafio-globo.com/models"
	"github.com/filiponegrao/desafio-globo.com/version"

	"github.com/gin-gonic/gin"
)

func GetBookmarks(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Bookmark{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db, err = parameter.Paginate(db)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	db = parameter.SortRecords(db)
	db = parameter.FilterFields(db)
	bookmarks := []models.Bookmark{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Bookmark{}, fields)

	if err := db.Select(queryFields).Find(&bookmarks).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(bookmarks) > 0 {
		index = int(bookmarks[len(bookmarks)-1].ID)
	}

	if err := parameter.SetHeaderLink(c, index); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	if _, ok := c.GetQuery("stream"); ok {
		enc := json.NewEncoder(c.Writer)
		c.Status(200)

		for _, bookmark := range bookmarks {
			fieldMap, err := helper.FieldToMap(bookmark, fields)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			if err := enc.Encode(fieldMap); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
		}
	} else {
		fieldMaps := []map[string]interface{}{}

		for _, bookmark := range bookmarks {
			fieldMap, err := helper.FieldToMap(bookmark, fields)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			fieldMaps = append(fieldMaps, fieldMap)
		}

		if _, ok := c.GetQuery("pretty"); ok {
			c.IndentedJSON(200, fieldMaps)
		} else {
			c.JSON(200, fieldMaps)
		}
	}
}

func GetBookmark(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Bookmark{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	bookmark := models.Bookmark{}
	id := c.Params.ByName("id")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Bookmark{}, fields)

	if err := db.Select(queryFields).First(&bookmark, id).Error; err != nil {
		content := gin.H{"error": "bookmark with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	fieldMap, err := helper.FieldToMap(bookmark, fields)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	if _, ok := c.GetQuery("pretty"); ok {
		c.IndentedJSON(200, fieldMap)
	} else {
		c.JSON(200, fieldMap)
	}
}

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

	c.Redirect(303, "/bookmarks")
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
	log.Println("PASSANDO AQUI")
	ShowBookMarksPage(c, userID, "URL Excluída com sucesso.")
}
