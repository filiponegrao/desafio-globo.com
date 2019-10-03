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

	r.Use(LoginInterceptor())

	// Metodos sem autorizacao
	r.GET("/login", controllers.GetLoginPage)
	r.POST("/users", controllers.CreateUser)
	r.GET("/register", controllers.GetRegsisterPage)
	r.GET("/forgot-password", controllers.GetForgotPasswordPage)

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
		LoginResponse:   LoginResponse,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	api := r.Group("")
	// api.POST("/login", authMiddleware.LoginHandler)
	api.POST("/login", authMiddleware.LoginHandler)
	api.POST("/forgot-password", controllers.ForgotPassword)

	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/bookmarks", controllers.GetBookmarksPage)
		api.GET("/create-bookmark", controllers.GetCreateBookmarkPage)
		api.POST("/bookmarks", controllers.CreateBookmark)

		api.GET("/bookmarks/:id", controllers.GetBookmark)
		// api.PUT("/bookmarks/:id", controllers.UpdateBookmark)
		api.POST("/delete-bookmark/:id", controllers.DeleteBookmark)

		// api.GET("/users", controllers.GetUsers)
		// api.GET("/users/:id", controllers.GetUser)
		// api.PUT("/users/:id", controllers.UpdateUser)
		// api.DELETE("/users/:id", controllers.DeleteUser)

	}
}

func RedirectToLogin(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/login")
}

func LoginResponse(c *gin.Context, n int, token string, time time.Time) {
	tokenString := "Bearer " + token
	c.SetCookie("authorization", tokenString, 3600, "", "localhost:8080", false, false)
	c.Set("authorization", tokenString)
	c.Redirect(303, "/bookmarks")
}

func LoginInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := c.Cookie("authorization")
		if err != nil {
			log.Println(err)
		} else {
			c.Request.Header.Set("Authorization", resp)
		}
		c.Next()
	}
}
