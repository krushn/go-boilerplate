package main

import (
	"go-boilerplate/controllers"

	"github.com/gin-gonic/gin"
)

// RouteHandler Create add /example route to gin engine
func RouteHandler(app *gin.Engine) *gin.Engine {

	// Add route to the gin engine
	app.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := app.Group("/v1")
	{
		/*** START USER ***/
		user := new(controllers.UserController)

		v1.POST("/user/login", user.Login)
		v1.POST("/user/register", user.Register)
		v1.GET("/user/logout", user.Logout)

		/*** START AUTH ***/
		auth := new(controllers.AuthController)

		//Refresh the token when needed to generate new access_token and refresh_token for the user
		v1.POST("/token/refresh", auth.Refresh)

		v1.GET("/token/validate", auth.TokenValid)

		/*** START Article
		article := new(controllers.ArticleController)

		v1.POST("/article", TokenAuthMiddleware(), article.Create)
		v1.GET("/articles", TokenAuthMiddleware(), article.All)
		v1.GET("/article/:id", TokenAuthMiddleware(), article.One)
		v1.PUT("/article/:id", TokenAuthMiddleware(), article.Update)
		v1.DELETE("/article/:id", TokenAuthMiddleware(), article.Delete)***/
	}

	//app.LoadHTMLGlob("./public/html/*")

	//app.Static("/public", "./public")

	// return gin engine with newly added route
	return app
}
