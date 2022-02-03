package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zzz4zzz/go-playground/db/models"
	"github.com/zzz4zzz/go-playground/http/middlewares"
	"github.com/zzz4zzz/go-playground/trace"
)

type userParams struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func AddUserTraced(provider *trace.Provider) gin.HandlerFunc {
	tracer := provider.Provider.Tracer("user-trace")

	return func(c *gin.Context) {
		var userParams = new(userParams)

		if err := c.BindJSON(userParams); err != nil {
			return
		}
		DB := middlewares.GetDBFromContext(c)
		user := &models.User{FirstName: userParams.FirstName, LastName: userParams.LastName}

		_, span := tracer.Start(c.Request.Context(), "creating-user")
		DB.Create(user)
		span.End()

		c.JSON(200, user)
	}
}
