package httpHandlers

import (
	"time"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Cheasezz/testForOzon/internal/app"
	"github.com/Cheasezz/testForOzon/internal/gql/resolvers"
	"github.com/Cheasezz/testForOzon/internal/gql/runtime"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
)

const websocketKeepAlivePingInterval = 5 * time.Second
const maxUploadSize = 30 * 1024 * 1024
const queryCacheLRUSize = 1000
const automaticPersistedQueryCacheLRUSize = 100
const complexityLimit = 1000

func NewGraphQLHandler(env *app.Env) gin.HandlerFunc {
	handler := gqlhandler.New(
		runtime.NewExecutableSchema(
			runtime.Config{Resolvers: resolvers.NewResolver(env)},
		),
	)

	handler.AddTransport(transport.Websocket{
		KeepAlivePingInterval: websocketKeepAlivePingInterval,
	})
	handler.AddTransport(transport.Options{})
	handler.AddTransport(transport.GET{})
	handler.AddTransport(transport.POST{})

	handler.SetQueryCache(lru.New[*ast.QueryDocument](queryCacheLRUSize))

	handler.Use(extension.Introspection{})

	handler.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](automaticPersistedQueryCacheLRUSize)})

	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
