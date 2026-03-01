package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/domain"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/handler/request"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/handler/response"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/usecase"
)

type EmployeeHandler struct {
	usecase   usecase.EmployeeUsecase
	validator *validator.Validate
}

func NewEmployeeHandler(u usecase.EmployeeUsecase, v *validator.Validate) *EmployeeHandler {
	return &EmployeeHandler{usecase: u, validator: v}
}

func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req request.CreateEmployeeRequest

	if err := c.ShouldBind(&req); err != nil {
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

	employee := &domain.Employee{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password,
		Role:         req.Role,
	}

	if err := h.usecase.CreateEmployee(employee); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    500,
			Message: "サーバー内部エラーが発生しました",
		})
		return
	}
	// PasswordHashは json:"-" なので除外される
	c.JSON(http.StatusCreated, employee)
}

func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	var pathID request.PathID
	if err := c.ShouldBindUri(&pathID); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    400,
			Message: "リクエスト形式が不正です",
		})
		return
	}

	var req request.UpdateEmployeeRequest
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

	if req.Name == nil && req.Email == nil && req.Password == nil && req.Role == nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    400,
			Message: "更新するフィールドを指定してください",
		})
		return
	}

	employee := &domain.Employee{ID: pathID.ID}
	if req.Name != nil {
		employee.Name = *req.Name
	}
	if req.Email != nil {
		employee.Email = *req.Email
	}
	if req.Password != nil {
		employee.PasswordHash = *req.Password
	}
	if req.Role != nil {
		employee.Role = *req.Role
	}

	updated, err := h.usecase.UpdateEmployee(employee)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    500,
			Message: "サーバー内部エラーが発生しました",
		})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.DeleteEmployee(id); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    500,
			Message: "サーバー内部エラーが発生しました",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *EmployeeHandler) GetEmployee(c *gin.Context) {
	id := c.Param("id")
	employee, err := h.usecase.GetEmployee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    500,
			Message: "サーバー内部エラーが発生しました",
		})
		return
	}
	if employee == nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{
			Code:    404,
			Message: "ユーザーが見つかりません",
		})
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    400,
			Message: "リクエスト形式が不正です",
		})
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    400,
			Message: "リクエスト形式が不正です",
		})
		return
	}

	employees, err := h.usecase.ListEmployees(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    500,
			Message: "サーバー内部エラーが発生しました",
		})
		return
	}

	c.JSON(http.StatusOK, employees)
}
