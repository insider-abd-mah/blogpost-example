package cmd

import (
	"blog-example/internal/platform/database"
	"github.com/gin-gonic/gin"
)

// ContextInterface ..
type ContextInterface interface {
	GetGinCtx() *gin.Context
	GetDatabaseConnection() database.DBConnectionInterface
}

// Context is a custom context
type Context struct {
	Ctx *gin.Context
}

// GetGinCtx ..
func (c *Context) GetGinCtx() *gin.Context {
	return c.Ctx
}

// GetDatabaseConnection ...
func (c *Context) GetDatabaseConnection() database.DBConnectionInterface {
	return database.GetConnection()
}
