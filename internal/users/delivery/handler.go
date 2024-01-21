package delivery

import (
	"UserService/internal/models"
	"UserService/internal/users"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	STATUS_OK    = "ok"
	STATUS_ERROR = "error"
)

type response struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func newResponse(status, msg string) *response {
	return &response{
		Status: status,
		Msg:    msg,
	}
}

type handler struct {
	service users.Service
	log     *slog.Logger
}

func newHandler(userService users.Service, log *slog.Logger) *handler {
	return &handler{
		service: userService,
		log:     log,
	}
}

func (h *handler) create(c *gin.Context) {
	inp := new(models.User)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	newUser, err := h.service.Create(*inp)
	if err != nil {
		h.log.Debug("ошибка при создании пользователя", slog.Any("error", err.Error()), slog.Any("request", inp))
		c.AbortWithStatusJSON(http.StatusBadRequest, newResponse(STATUS_ERROR, err.Error()))
		return
	}
	c.JSON(http.StatusOK, newUser)
}

func (h *handler) edit(c *gin.Context) {
	inp := new(users.EditUserRequest)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	newUser, err := h.service.Edit(*inp)
	if err != nil {
		h.log.Debug("ошибка при редактировании пользователя", slog.Any("error", err.Error()), slog.Any("request", inp))
		c.AbortWithStatusJSON(http.StatusBadRequest, newResponse(STATUS_ERROR, err.Error()))
		return
	}
	c.JSON(http.StatusOK, newUser)
}

func (h *handler) delete(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Debug("id пользователя имеет не верный формат")
		c.AbortWithStatusJSON(http.StatusBadRequest, newResponse(STATUS_ERROR, "id пользователя имеет не верный формат"))
		return
	}
	err = h.service.Delete(userId)
	if err != nil {
		h.log.Debug("ошибка при удалении пользователя пользователя", slog.Any("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, newResponse(STATUS_ERROR, err.Error()))
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func (h *handler) get(c *gin.Context) {
	filters := make(map[string]string)
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		pageSize = 10
	}
	filters["age"] = c.Query("age")
	filters["name"] = c.Query("name")
	filters["surname"] = c.Query("surname")
	filters["patronymic"] = c.Query("patronymic")
	filters["gender"] = c.Query("gender")
	filters["nationality"] = c.Query("nationality")

	users, err := h.service.GetWithFilter(filters, page, pageSize)

	if err != nil {
		h.log.Debug("ошибка при получении пользователей", slog.Any("error", err.Error()), slog.Any("filters", filters), slog.Any("page", page), slog.Any("pageSize", pageSize))
		c.AbortWithStatusJSON(http.StatusBadRequest, newResponse(STATUS_ERROR, err.Error()))
		return
	}
	c.JSON(http.StatusOK, users)
}
