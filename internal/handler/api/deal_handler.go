package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/domain"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/usecase"
)

type DealHandler struct {
	usecase usecase.DealUsecase
}

func NewDealHandler(u usecase.DealUsecase) *DealHandler {
	return &DealHandler{usecase: u}
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
	var deal domain.Deal
	if err := c.ShouldBindJSON(&deal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.CreateDeal(&deal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, deal)
}

func (h *DealHandler) UpdateDeal(c *gin.Context) {
	id := c.Param("id")
	var deal domain.Deal
	if err := c.ShouldBindJSON(&deal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	deal.ID = id

	updated, err := h.usecase.UpdateDeal(&deal)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// UpdateDealStatus represents the kanban movement
func (h *DealHandler) UpdateDealStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status     string  `json:"status"`
		AssigneeID *string `json:"assignee_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.usecase.UpdateDealStatus(id, req.Status, req.AssigneeID); err != nil {
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
