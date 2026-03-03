package document

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
	r.POST("/tenants/:id/documents", h.UploadDocument)
	r.GET("/tenants/:id/documents", h.GetDocuments)
}

func (h *Handler) UploadDocument(c *gin.Context) {
	tenantID := c.Param("id")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	doc, err := h.service.UploadDocument(c.Request.Context(), tenantID, file)
	if err != nil {
		if err.Error() == "tenant not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, doc)
}

func (h *Handler) GetDocuments(c *gin.Context) {
	tenantID := c.Param("id")

	docs, err := h.service.GetDocuments(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, docs)
}