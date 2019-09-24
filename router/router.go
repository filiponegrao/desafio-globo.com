package router

import (
	"github.com/filiponegrao/desafio-globo.com/controllers"

	"github.com/gin-gonic/gin"
)

func Initialize(r *gin.Engine) {

	r.LoadHTMLGlob("view/*")

	// r.GET("/", controllers.APIEndpoi		nts)
	r.GET("/login", controllers.GetLoginPage)

	api := r.Group("")
	{
		api.GET("/bookmarks", controllers.GetBookmarks)
		api.GET("/bookmarks/:id", controllers.GetBookmark)
		api.POST("/bookmarks", controllers.CreateBookmark)
		api.PUT("/bookmarks/:id", controllers.UpdateBookmark)
		api.DELETE("/bookmarks/:id", controllers.DeleteBookmark)

		api.GET("/users", controllers.GetUsers)
		api.GET("/users/:id", controllers.GetUser)
		api.POST("/users", controllers.CreateUser)
		api.PUT("/users/:id", controllers.UpdateUser)
		api.DELETE("/users/:id", controllers.DeleteUser)

	}
}
