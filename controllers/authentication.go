package controllers

import (
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/badoux/checkmail"
	dbpkg "github.com/filiponegrao/desafio-globo.com/db"
	"github.com/filiponegrao/desafio-globo.com/models"
	"github.com/filiponegrao/desafio-globo.com/tools"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password"`
}

type newPassword struct {
	OldPassowrd     string `form:"oldPassword" json:"oldPassword"`
	NewPassword     string `form:"newPassword" json:"newPassword"`
	ConfirmPassowrd string `form:"confirmPassowrd" json:"confirmPassowrd"`
}

func AuthorizationPayload(data interface{}) jwt.MapClaims {
	if user, ok := data.(*models.User); ok {
		return jwt.MapClaims{
			"id": user.ID,
		}
	}
	return jwt.MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &models.User{
		ID: int64(claims["id"].(float64)),
	}
}

// Falha na autênticação
func UserUnauthorized(c *gin.Context, code int, message string) {
	err := ""
	if strings.Contains(message, "missing") {
		err = "Faltando email ou senha"
	} else if strings.Contains(message, "incorrect") {
		err = "Email ou senha incorreta"
	} else if strings.Contains(message, "cookie token is empty") {
		err = "Faltando HEADER de autenticação!"
		c.Redirect(303, "/login")
	} else {
		err = message
	}
	c.JSON(code, gin.H{"error": err})
}

func UserAuthentication(c *gin.Context) (interface{}, error) {
	token := csrf.GetToken(c)

	var loginVals login

	if err := c.Bind(&loginVals); err != nil {
		return nil, err
	}

	email := loginVals.Username
	password := loginVals.Password

	db := dbpkg.DBInstance(c)

	if email == "" {
		message := "Faltando email"
		// c.JSON(400, gin.H{"error": message})
		c.HTML(200, "login.html", gin.H{"message": message, "token": token})
		return nil, nil
	}

	if password == "" {
		message := "Faltando senha (password)"
		c.HTML(200, "login.html", gin.H{"message": message, "token": token})
		return nil, nil
	}

	var user models.User

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		message := "Usuário ou senha incorretos"
		c.HTML(200, "login.html", gin.H{"message": message, "token": token})
		return nil, nil
	}

	encPassword := tools.EncryptTextSHA512(password)

	if encPassword != user.Password {
		message := "Usuário ou senha incorretos"
		c.HTML(200, "login.html", gin.H{"message": message})
		return nil, nil
	}

	user.Password = ""

	return &user, nil
}

func UserAuthorization(user interface{}, c *gin.Context) bool {
	return true
}

func ForgotPassword(c *gin.Context) {
	token := csrf.GetToken(c)

	db := dbpkg.DBInstance(c)

	email := c.PostForm("email")

	var user models.User

	if email == "" {
		c.JSON(400, gin.H{"error": "Faltando parametro de email", "token": token})
		return
	}

	err := checkmail.ValidateFormat(email)
	if err != nil {
		message := err.Error()
		c.HTML(200, "forgot-password.html", gin.H{"message": message, "token": token})
		return
	}

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		message := err.Error()
		c.HTML(200, "forgot-password.html", gin.H{"message": message, "token": token})
		return
	}

	contentHash := user.Email + time.Now().String()
	encodedContent := tools.EncryptTextSHA512(contentHash)
	var recover models.PasswordRecover
	recover.UserID = user.ID
	recover.Hash = encodedContent

	if err := db.Create(&recover).Error; err != nil {
		message := err.Error()
		c.HTML(200, "forgot-password.html", gin.H{"message": message, "token": token})
		return
	}

	path := "/new-password/" + encodedContent

	EmailChangedPassword(email, path)

	c.HTML(200, "login.html", gin.H{"message": "Instruções enviadas com sucesso!", "token": token})

}

func NewPassword(c *gin.Context) {
	token := csrf.GetToken(c)

	db := dbpkg.DBInstance(c)

	hash := c.Param("hash")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirmPassword")

	if password != confirmPassword {
		message := "Senhas precisam ser iguais"
		c.HTML(200, "new-password.html", gin.H{"message": message, "token": token})
		return
	}

	if !checkPassword(password) {
		message := "Senha não confere com o padrão desejado!"
		c.HTML(200, "new-password.html", gin.H{"message": message, "token": token})
		return
	}

	var recover models.PasswordRecover
	if err := db.Where("hash = ?", hash).First(&recover).Error; err != nil {
		message := "Hash inválido!"
		c.HTML(200, "new-password.html", gin.H{"message": message, "token": token})
		return
	}

	var user models.User
	if err := db.First(&user, recover.UserID).Error; err != nil {
		message := err.Error()
		c.HTML(200, "new-password.html", gin.H{"message": message, "token": token})
		return
	}

	encoded := tools.EncryptTextSHA512(password)

	tx := db.Begin()
	user.Password = encoded
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		message := err.Error()
		c.HTML(200, "new-password.html", gin.H{"message": message, "token": token})
		return
	}

	if err := tx.Delete(&recover).Error; err != nil {
		tx.Rollback()
		message := err.Error()
		c.HTML(200, "new-password.html", gin.H{"message": message, "token": token})
		return
	}

	tx.Commit()

	c.Redirect(303, "/login")
}
