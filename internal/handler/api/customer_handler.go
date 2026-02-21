package api

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/domain"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/usecase"
)

type CustomerHandler struct {
	usecase usecase.CustomerUsecase
}

func NewCustomerHandler(u usecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{usecase: u}
}

func (h *CustomerHandler) ListCustomers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit format"})
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset format"})
		return
	}

	customers, err := h.usecase.ListCustomers(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	customer, err := h.usecase.GetCustomer(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var customer domain.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.CreateCustomer(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer domain.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customer.ID = id

	if err := h.usecase.UpdateCustomer(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.DeleteCustomer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
