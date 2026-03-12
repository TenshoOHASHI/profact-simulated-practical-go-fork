package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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

	if req.Name == "" && req.Rent == 0 && req.Address == "" && req.Layout == nil && req.Status == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    400,
			Message: "更新するフィールドを指定してください",
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

func (h *PropertyHandler) ImportProperties(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    400,
			Message: "ファイルを選択してください",
		})
		return
	}
	defer file.Close()

	const maxFileSize = 10 << 20 // 10MB
	if header.Size > maxFileSize {
		c.JSON(http.StatusRequestEntityTooLarge, response.ErrorResponse{
			Code:    413,
			Message: "ファイルサイズが大きすぎます（最大10MB）",
		})
		return
	}

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".csv") {
		c.JSON(http.StatusUnsupportedMediaType, response.ErrorResponse{
			Code:    415,
			Message: "CSVファイルのみアップロード可能です",
		})
		return
	}

	result, validationErrors, err := h.usecase.ImportProperties(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    500,
			Message: "サーバー内部エラーが発生しました",
		})
		return
	}
	if len(validationErrors) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "インポートに失敗しました",
			"errors":  validationErrors,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "インポートが完了しました",
		"data":    result,
	})
}

func (h *PropertyHandler) ExportProperties(c *gin.Context) {
	csvData, err := h.usecase.ExportProperties()
	if err != nil {
		// デバッグ用: エラー詳細を出力
		fmt.Printf("Export error: %v\n", err)
		if err.Error() == "no data" {
			c.JSON(http.StatusNotFound, response.ErrorResponse{
				Code:    404,
				Message: "エクスポート対象のデータが存在しません",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    500,
			Message: "エクスポートに失敗しました",
		})
		return
	}

	filename := fmt.Sprintf("properties_%s.csv", time.Now().Format("20060102"))
	c.Header("Content-Type", "text/csv; charset=UTF-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "text/csv; charset=UTF-8", csvData)
}
