package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/tenants/:id/users", h.CreateUser)
	r.GET("/tenants/:id/users", h.GetUsers)
}

type createUserRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *Handler) CreateUser(c *gin.Context) {
	tenantID := c.Param("id")

	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), tenantID, req.Email)
	if err != nil {
		if err.Error() == "tenant not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *Handler) GetUsers(c *gin.Context) {
	tenantID := c.Param("id")

	users, err := h.service.GetUsers(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}