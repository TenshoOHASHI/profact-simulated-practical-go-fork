package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/domain"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/handler/request"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/handler/response"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/usecase"
)

type PropertyHandler struct {
	usecase   usecase.PropertyUsecase
	validator *validator.Validate
}

func NewPropertyHandler(u usecase.PropertyUsecase, v *validator.Validate) *PropertyHandler {
	return &PropertyHandler{usecase: u, validator: v}
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
	var req request.CreatePropertyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    400,
			Message: "リクエスト形式が不正です",
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    400,
			Message: "入力内容にエラーがあります",
			Errors:  response.FormatValidationErrors(err),
		})
		return
	}

	property := &domain.Property{
		Name:    req.Name,
		Rent:    req.Rent,
		Address: req.Address,
		Layout:  &req.Layout,
		Status:  req.Status,
	}

	if err := h.usecase.CreateProperty(property); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    500,
			Message: "サーバー内部エラーが発生しました",
		})
		return
	}

	c.JSON(http.StatusCreated, property)
}

func (h *PropertyHandler) UpdateProperty(c *gin.Context) {

	var pathID request.PathID
	if err := c.ShouldBindUri(&pathID); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    400,
			Message: "入力内容にエラーがあります",
			Errors:  response.FormatValidationErrors(err),
		})
		return
	}

	var req request.UpdatePropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    400,
			Message: "入力内容にエラーがあります",
			Errors:  response.FormatValidationErrors(err),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    400,
			Message: "入力内容にエラーがあります",
			Errors:  response.FormatValidationErrors(err),
		})
		return
	}

	property := &domain.Property{
		ID:      pathID.ID,
		Name:    req.Name,
		Rent:    req.Rent,
		Address: req.Address,
		Layout:  req.Layout,
		Status:  req.Status,
	}

	updated, err := h.usecase.UpdateProperty(property)

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
