package router

import (
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/filiponegrao/desafio-globo.com/controllers"

	"github.com/gin-gonic/gin"
)

func Initialize(r *gin.Engine) {

	controllers.Config(r)

	r.LoadHTMLGlob("view/*")

	// Metodos sem autorizacao
	r.GET("/login", controllers.GetLoginPage)
	r.POST("/users", controllers.CreateUser)
	r.GET("/register", controllers.GetRegsisterPage)

	r.GET("", RedirectToLogin)

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "test zone",
		Key:             []byte("secret key"),
		Timeout:         time.Hour * 24 * 7,
		MaxRefresh:      time.Hour,
		IdentityKey:     "id",
		PayloadFunc:     controllers.AuthorizationPayload,
		IdentityHandler: controllers.IdentityHandler,
		Authenticator:   controllers.UserAuthentication,
		Authorizator:    controllers.UserAuthorization,
		Unauthorized:    controllers.UserUnauthorized,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	api := r.Group("")
	// api.POST("/login", authMiddleware.LoginHandler)
	api.POST("/login", func(c *gin.Context) {
		token := authMiddleware.LoginHandler(c)
		authToken := "Bearer " + token
		c.SetCookie("Authorization", authToken, 3600, "/", "localhost", false, true)
		c.Request.URL.Path = "/bookmarks"
		c.Request.Method = "GET"
		log.Println(authToken)
		c.Request.Header.Add("Authorization", authToken)
		r.HandleContext(c)
	})

	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/bookmarks", controllers.GetBookmarksPage)
		// api.GET("/bookmarks/:id", controllers.GetBookmark)
		// api.POST("/bookmarks", controllers.CreateBookmark)
		// api.PUT("/bookmarks/:id", controllers.UpdateBookmark)
		// api.DELETE("/bookmarks/:id", controllers.DeleteBookmark)

		// api.GET("/users", controllers.GetUsers)
		// api.GET("/users/:id", controllers.GetUser)
		// api.PUT("/users/:id", controllers.UpdateUser)
		// api.DELETE("/users/:id", controllers.DeleteUser)

	}
}

func RedirectToLogin(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/login")
}
