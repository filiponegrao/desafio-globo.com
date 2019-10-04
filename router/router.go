package router

import (
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/filiponegrao/desafio-globo.com/controllers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	csrf "github.com/utrack/gin-csrf"

	"github.com/gin-gonic/gin"
)

func Initialize(r *gin.Engine) {

	controllers.Config(r)

	r.LoadHTMLGlob("view/*")

	r.Use(LoginInterceptor())

	// Anti CSRF
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.Use(csrf.Middleware(csrf.Options{
		Secret: "DesafioGlobo.Com.MyBookmarks",
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	// Metodos sem autorizacao

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
	api.GET("/login", controllers.GetLoginPage)
	api.POST("/users", controllers.CreateUser)
	api.GET("/register", controllers.GetRegsisterPage)
	api.GET("/forgot-password", controllers.GetForgotPasswordPage)
	api.GET("/new-password/:hash", controllers.GetNewPasswordPage)

	r.GET("", RedirectToLogin)
	api.POST("/login", authMiddleware.LoginHandler)
	api.POST("/forgot-password", controllers.ForgotPassword)
	api.POST("/new-password/:hash", controllers.NewPassword)

	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/bookmarks", controllers.GetBookmarksPage)
		api.GET("/create-bookmark", controllers.GetCreateBookmarkPage)
		api.POST("/bookmarks", controllers.CreateBookmark)

		api.POST("/delete-bookmark/:id", controllers.DeleteBookmark)
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
		// Secure actions
		c.Header("X-Frame-Options", "deny")
		c.Header("X-XSS-Protection", "1")
		c.Header("X-Content-Type-Options", "nosniff")

		resp, err := c.Cookie("authorization")
		if err != nil {
			log.Println(err)
		} else {
			c.Request.Header.Set("Authorization", resp)
		}
		c.Next()
	}
}
