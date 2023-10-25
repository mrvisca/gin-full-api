package main

import (
	"gin-full-api/config"
	"gin-full-api/middleware"
	"gin-full-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

func main() {
	// Setup database
	config.InitDB()
	defer config.DB.Close()
	gotenv.Load()

	// Setup Routing
	router := gin.Default()

	v1 := router.Group("/api/v1/")
	{
		v1.GET("/auth/:provider", routes.RedirectHandler)
		v1.GET("/auth/:provider/callback", routes.CallbackHandler)

		// Route testing token
		v1.GET("/check", middleware.IsAuth(), routes.CheckToken)

		v1.GET("/article/:slug", routes.GetArticle)
		artikel := v1.Group("/artikel")
		{
			artikel.GET("/", routes.GetHome)
			artikel.GET("/tag/:tag", routes.GetArticleByTag)
			artikel.POST("/articles", middleware.IsAuth(), routes.PostArticle)
			artikel.PUT("/update/:id", middleware.IsAuth(), routes.UpdateArticle)
		}
	}

	router.Run()
}
