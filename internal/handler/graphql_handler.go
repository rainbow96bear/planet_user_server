package handler

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_user_server/graph"
	"github.com/rainbow96bear/planet_user_server/internal/resolver"
	"github.com/rainbow96bear/planet_user_server/middleware"
)

type GraphqlHandler struct {
	server *handler.Server
}

func NewGraphqlHandler(r *resolver.Resolver) *GraphqlHandler {
	exec := graph.NewExecutableSchema(graph.Config{
		Resolvers: r,
	})

	return &GraphqlHandler{
		server: handler.NewDefaultServer(exec),
	}
}

func (h *GraphqlHandler) Graphql() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.server.ServeHTTP(c.Writer, c.Request)
	}
}

func (h *GraphqlHandler) Playground() gin.HandlerFunc {
	return gin.WrapH(playground.Handler("Planet User GraphQL", "/graphql"))
}

func (h *GraphqlHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/graphql", middleware.AuthMiddleware(), h.Graphql())
	r.GET("/playground", h.Playground())
}
