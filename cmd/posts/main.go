package main

import (
	"blog-example/cmd"
	"blog-example/cmd/posts/handlers"
	"blog-example/internal/platform/database"
	"blog-example/log"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
)

func main() {
	log.Init()
	database.Init()

	if err := setupServer().Run(":" + os.Getenv("APP_PORT")); err != nil {
		log.Fatal(err)
	}
}

func setupServer() *gin.Engine {
	gin.EnableJsonDecoderUseNumber()
	g := gin.New()

	// Used middlewares
	g.Use(gin.Recovery())

	g.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"statusMessage": "OK", "serviceName": "posts"})
	})

	// version 1 async routes
	v1 := g.Group("posts/v1")
	{
		v1.POST("store", func(c *gin.Context) {
			handlers.StorePosts(&cmd.Context{Ctx: c})
		})
		v1.GET("all", func(c *gin.Context) {
			handlers.GetAllPosts(&cmd.Context{Ctx: c})
		})
	}

	return g
}
