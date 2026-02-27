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

type DealHandler struct {
	usecase   usecase.DealUsecase
	validator *validator.Validate
}

func NewDealHandler(u usecase.DealUsecase, v *validator.Validate) *DealHandler {
	return &DealHandler{usecase: u, validator: v}
}

func (h *DealHandler) ListDeals(c *gin.Context) {
	deals, err := h.usecase.ListDeals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deals)
}

func (h *DealHandler) GetDeal(c *gin.Context) {
	id := c.Param("id")
	deal, err := h.usecase.GetDeal(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if deal == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deal not found"})
		return
	}

	c.JSON(http.StatusOK, deal)
}

func (h *DealHandler) CreateDeal(c *gin.Context) {
	var req request.CreateDealRequest

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

	propertyID := req.PropertyID
	assigneeID := req.AssigneeID
	deal := &domain.Deal{
		CustomerID: req.CustomerID,
		PropertyID: &propertyID,
		AssigneeID: &assigneeID,
		Status:     req.Status,
		MoveInDate: req.MoveInDate,
	}

	if err := h.usecase.CreateDeal(deal); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    500,
			Message: "サーバー内部エラーが発生しました",
		})
		return
	}

	c.JSON(http.StatusCreated, deal)
}

func (h *DealHandler) UpdateDeal(c *gin.Context) {
	var pathID request.PathID
	if err := c.ShouldBindUri(&pathID); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    400,
			Message: "入力内容にエラーがあります",
			Errors:  response.FormatValidationErrors(err),
		})
		return
	}

	var req request.UpdateDealRequest
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

	deal := &domain.Deal{
		ID:         pathID.ID,
		CustomerID: req.CustomerID,
		PropertyID: req.PropertyID,
		AssigneeID: req.AssigneeID,
		Status:     req.Status,
		MoveInDate: req.MoveInDate,
	}

	updated, err := h.usecase.UpdateDeal(deal)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// UpdateDealStatus represents the kanban movement
func (h *DealHandler) UpdateDealStatus(c *gin.Context) {
	var pathID request.PathID
	if err := c.ShouldBindUri(&pathID); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    400,
			Message: "入力内容にエラーがあります",
			Errors:  response.FormatValidationErrors(err),
		})
		return
	}

	var req request.UpdateDealStatusRequest
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

	if _, err := h.usecase.UpdateDealStatus(pathID.ID, req.Status, req.AssigneeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *DealHandler) DeleteDeal(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.DeleteDeal(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
