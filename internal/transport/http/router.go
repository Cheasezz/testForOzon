package httpHandlers

import (
	"net/http"

	"github.com/Cheasezz/testForOzon/internal/app"
	"github.com/gin-gonic/gin"
)

func New(env *app.Env, port string) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(NewDataLoadersInjector(env))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Page not found"})
	})

	// Из документации gqlgen:
	// If you are using an external router, remember to send ALL requests go to your handler (at /query), not just POST requests!
	router.Any("/query", NewGraphQLHandler(env))

	router.GET("/", playgroundHandler())

	env.Logger.Info("connect to http://localhost:%s/ for GraphQL playground", port)
	return router
}
