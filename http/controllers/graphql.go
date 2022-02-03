package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	"github.com/zzz4zzz/go-playground/http/graphql"
)

var graphqlHandler = handler.New(&handler.Config{
	Schema:     &graphql.Schema,
	Pretty:     true,
	GraphiQL:   true,
	Playground: true,
})

func GinGraphQLHandler(c *gin.Context) {
	graphqlHandler.ServeHTTP(c.Writer, c.Request)
}
