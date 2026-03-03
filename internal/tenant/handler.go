package tenant

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
	r.POST("/tenants", h.CreateTenant)
	r.GET("/tenants/:id", h.GetTenant)
}

type createTenantRequest struct {
	Name string `json:"name"`
}

func (h *Handler) CreateTenant(c *gin.Context) {
	var req createTenantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	tenant, err := h.service.CreateTenant(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, tenant)
}

func (h *Handler) GetTenant(c *gin.Context) {
	tenantID := c.Param("id")

	tenant, err := h.service.GetTenant(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if tenant == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "tenant not found",
		})
		return
	}

	c.JSON(http.StatusOK, tenant)
}