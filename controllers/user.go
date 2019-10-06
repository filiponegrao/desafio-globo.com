package controllers

import (
	"strings"

	"github.com/badoux/checkmail"
	dbpkg "github.com/filiponegrao/desafio-globo.com/db"
	"github.com/filiponegrao/desafio-globo.com/models"
	"github.com/filiponegrao/desafio-globo.com/tools"
	csrf "github.com/utrack/gin-csrf"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func Config(r *gin.Engine) {
	Router = r
}

func CreateUser(c *gin.Context) {
	token := csrf.GetToken(c)

	db := dbpkg.DBInstance(c)
	user := models.User{}

	if err := c.Bind(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	missing := user.MissingFields()
	if missing != "" {
		message := "Faltando campo de " + missing
		c.HTML(200, "register.html", gin.H{"message": message, "token": token})
		return
	}

	// Valida o email
	err := checkmail.ValidateFormat(user.Email)
	if err != nil {
		message := "E-mail não possui um formato valido"
		c.HTML(200, "register.html", gin.H{"message": message, "token": token})
		return
	}

	if user.Password != user.ConfirmPassowrd {
		message := "Senhas precisam ser iguais"
		c.HTML(200, "register.html", gin.H{"message": message, "token": token})
		return
	}

	if !checkPassword(user.Password) {
		message := "Senha não confere com o padrão desejado!"
		c.HTML(200, "register.html", gin.H{"message": message, "token": token})
		return
	}

	passwordEncoded := tools.EncryptTextSHA512(user.Password)
	user.Password = passwordEncoded

	if err := db.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.HTML(200, "login.html", gin.H{"message": "Usuário cadastrado!", "token": token})
}

func checkPassword(password string) bool {
	if len(password) < 6 {
		return false
	} else if !strings.ContainsAny(password, "0123456789") {
		return false
	} else if !strings.ContainsAny(password, "!@#$%*()-=+<>;:|/\\") {
		return false
	}
	return true
}
