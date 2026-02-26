package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/domain"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/usecase"
)

type PropertyHandler struct {
	usecase usecase.PropertyUsecase
}

func NewPropertyHandler(u usecase.PropertyUsecase) *PropertyHandler {
	return &PropertyHandler{usecase: u}
}

func (h *PropertyHandler) ListProperties(c *gin.Context) {
	properties, err := h.usecase.ListProperties()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, properties)
}

func (h *PropertyHandler) GetProperty(c *gin.Context) {
	id := c.Param("id")
	property, err := h.usecase.GetProperty(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if property == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}

	c.JSON(http.StatusOK, property)
}

func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	var property domain.Property
	if err := c.ShouldBindJSON(&property); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.CreateProperty(&property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, property)
}

func (h *PropertyHandler) UpdateProperty(c *gin.Context) {
	id := c.Param("id")
	var property domain.Property
	if err := c.ShouldBindJSON(&property); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	property.ID = id

	updated, err := h.usecase.UpdateProperty(&property)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *PropertyHandler) DeleteProperty(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.DeleteProperty(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
