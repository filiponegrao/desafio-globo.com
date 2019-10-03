package controllers

import (
	"log"
	"strings"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/badoux/checkmail"
	dbpkg "github.com/filiponegrao/desafio-globo.com/db"
	"github.com/filiponegrao/desafio-globo.com/models"
	"github.com/filiponegrao/desafio-globo.com/tools"
	"github.com/gin-gonic/gin"
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
		c.HTML(200, "login.html", gin.H{"message": message})
		return nil, nil
	}

	if password == "" {
		message := "Faltando senha (password)"
		c.HTML(200, "login.html", gin.H{"message": message})
		return nil, nil
	}

	var user models.User

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		//message := "Usuario com email " + email + " nao encontrado."
		//c.JSON(400, gin.H{"error": message})
		c.HTML(200, "login.html", gin.H{"message": "Usuário ou senha incorretos"})
		return nil, nil
	}

	encPassword := tools.EncryptTextSHA512(password)

	if encPassword != user.Password {
		//message := "Senha incorreta"
		//c.JSON(400, gin.H{"error": message})
		c.HTML(200, "login.html", gin.H{"message": "Usuário ou senha incorretos"})

		return nil, nil
	}

	user.Password = ""

	return &user, nil
}

func UserAuthorization(user interface{}, c *gin.Context) bool {
	return true
}

func ForgotPassword(c *gin.Context) {

	db := dbpkg.DBInstance(c)

	body := c.Request.RequestURI
	log.Println(body)
	var user models.User

	parts := strings.Split(body, "email=")
	if len(parts) <= 1 {
		c.JSON(400, gin.H{"error": "Faltando parametro de email"})
		return
	}
	email := parts[1]
	err := checkmail.ValidateFormat(user.Email)
	if err != nil {
		message := "E-mail não possui um formato valido"
		c.HTML(200, "forgot-password.html", gin.H{"message": message})
		return
	}

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// password := tools.RandomString(6)
	// passwordEncode := tools.EncryptTextSHA512(password)

	// user.Password = passwordEncode

	// if err := db.Save(user).Error; err != nil {
	// 	c.JSON(400, gin.H{"error": err.Error()})
	// 	return
	// }

	// EmailPasswordNew(user.Email, password)

	c.JSON(200, "Nova senha enviada para o email!")
}
