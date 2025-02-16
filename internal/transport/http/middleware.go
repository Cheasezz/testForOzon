package httpHandlers

import (
	"context"

	"github.com/Cheasezz/testForOzon/internal/app"
	"github.com/Cheasezz/testForOzon/internal/repositories/loaders"
	"github.com/gin-gonic/gin"
)

func NewDataLoadersInjector(env *app.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(
			c.Request.Context(),
			loaders.DataLoadersContextKey,
			loaders.NewDataLoaders(env.Repositories),
		)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
