package delivery

import (
	"UserService/internal/users"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(
	router *gin.RouterGroup,
	useCase users.Service,
	log *slog.Logger,
) {
	h := newHandler(useCase, log)

	users := router.Group("users")
	{
		users.POST("/create", h.create)
		users.PUT("/edit", h.edit)
		users.DELETE("/:id", h.delete)
		users.GET("/", h.get)
	}
}
