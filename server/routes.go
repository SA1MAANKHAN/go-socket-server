package server

import (
	"socket_server/server/handler"
	"socket_server/server/provider"
	"socket_server/server/repository"
	"socket_server/server/service"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConfigureRoutes(server *Server, socketServer *socketio.Server) {
	homeHandler := handler.HomeHandler{}
	postHandler := handler.PostHandler{DB: server.db}
	registerHandler := handler.NewRegisterHandler()

	jwtAuth := provider.NewJwtAuth(server.db)

	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.engine.GET("/socket.io/*any", gin.WrapH(socketServer))

	server.engine.POST(
		"/users",
		registerHandler.RegisterUser(service.NewUserService(repository.NewUsersRepository(server.db))),
	)

	server.engine.POST("/login", jwtAuth.Middleware().LoginHandler)

	needsAuth := server.engine.Group("/").Use(jwtAuth.Middleware().MiddlewareFunc())
	needsAuth.GET("/", homeHandler.Index())
	needsAuth.GET("/refresh", jwtAuth.Middleware().RefreshHandler)
	needsAuth.POST("/posts", postHandler.SavePost)
	needsAuth.GET("/posts", postHandler.GetPosts)
	needsAuth.GET("/post/:id", postHandler.GetPostByID)
	needsAuth.PUT("/post/:id", postHandler.UpdatePost)
	needsAuth.DELETE("/post/:id", postHandler.DeletePost)

}
