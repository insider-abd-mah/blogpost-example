package handlers

import (
	"blog-example/cmd"
	"blog-example/cmd/posts/internal"
	"blog-example/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

// StorePosts ...
//
// Store posts in DB
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Responses:
//  400: bad-request
// 	200: ok
func StorePosts(context cmd.ContextInterface) {
	var p internal.Post
	ctx := context.GetGinCtx()

	if err := ctx.ShouldBindBodyWith(&p, binding.JSON); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"statusMessage": "validation error"})

		return
	}

	db := context.GetDatabaseConnection()

	if err := p.Insert(db); err != nil {
		log.ErrChan <- fmt.Errorf("could not insert post, error givin is: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Internal server error"})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"statusMessage": "ok"})
}

// GetAllPosts ...
//
// Get all the posts from DB
//
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Responses:
//  400: bad-request
// 	200: ok
func GetAllPosts(context cmd.ContextInterface) {
	ctx := context.GetGinCtx()
	db := context.GetDatabaseConnection()
	posts, err := internal.GetAllPosts(db)

	if err != nil {
		log.ErrChan <- fmt.Errorf("could not retrive posts, error givin is: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Internal server error"})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"statusMessage": "ok", "data": posts})
}
