package usecase

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/yamu-studio/profact-simulated-practical-go/internal/domain"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/handler/request"
	"github.com/yamu-studio/profact-simulated-practical-go/internal/handler/validator"
	"golang.org/x/text/encoding/japanese"
)

// 簡易的なステータスバリデーション例
var validPropertyStatuses = map[string]bool{
	"available":  true,
	"contracted": true,
	"hidden":     true,
}

type PropertyUsecase interface {
	ListProperties() ([]*domain.Property, error)
	GetProperty(id string) (*domain.Property, error)
	CreateProperty(property *domain.Property) error
	UpdateProperty(property *domain.Property) (*domain.Property, error)
	DeleteProperty(id string) error
	ImportProperties(file multipart.File) (*request.ImportResult, []request.ValidationError, error)
	ExportProperties() ([]byte, error)
}

type propertyUsecase struct {
	repo domain.PropertyRepository
}

func NewPropertyUsecase(repo domain.PropertyRepository) PropertyUsecase {
	return &propertyUsecase{repo: repo}
}

func (u *propertyUsecase) ListProperties() ([]*domain.Property, error) {
	return u.repo.FindAll()
}

func (u *propertyUsecase) GetProperty(id string) (*domain.Property, error) {
	return u.repo.FindByID(id)
}

func (u *propertyUsecase) CreateProperty(property *domain.Property) error {
	return u.repo.Create(property)
}

func (u *propertyUsecase) UpdateProperty(property *domain.Property) (*domain.Property, error) {
	existing, err := u.repo.FindByID(property.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("property not found")
	}

	if property.Name != "" {
		existing.Name = property.Name
	}

	if property.Rent != 0 {
		existing.Rent = property.Rent
	}

	if property.Address != "" {
		existing.Address = property.Address
	}

	if property.Layout != nil {
		existing.Layout = property.Layout
	}

	if property.Status != "" && !validPropertyStatuses[property.Status] {
		return nil, errors.New("invalid property status transition")
	}

	if property.Status != "" {
		existing.Status = property.Status
	}
	if err := u.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, err
}

func (u *propertyUsecase) DeleteProperty(id string) error {
	return u.repo.Delete(id)
}

func (u *propertyUsecase) ImportProperties(file multipart.File) (*request.ImportResult, []request.ValidationError, error) {
	decoder := japanese.ShiftJIS.NewDecoder()
	reader := csv.NewReader(decoder.Reader(file))

	existingMap, err := u.repo.GetExistingPropertiesMap()
	if err != nil {
		return nil, nil, err
	}

	var properties []*domain.Property
	var errors []request.ValidationError
	seen := make(map[string]bool)
	result := &request.ImportResult{}

	_, err = reader.Read()
	if err != nil {
		return nil, nil, err
	}
	lineNumber := 2
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		errs := validator.ValidateCSRow(row, lineNumber)
		if len(errs) > 0 {
			errors = append(errors, errs...)
			lineNumber++
		}

		key := row[0] + "|" + row[2]
		if seen[key] {
			result.SkippedCount++
			result.SkippedItems = append(result.SkippedItems, request.SkippedItem{
				Row:    lineNumber,
				Reason: fmt.Sprintf("CSV内で重複（物件名：%s）", row[0]),
			})
			lineNumber++
			continue
		}
		seen[key] = true

		if existingMap[key] {
			result.SkippedCount++
			result.SkippedItems = append(result.SkippedItems, request.SkippedItem{
				Row:    lineNumber,
				Reason: fmt.Sprintf("既に登録済み（物件名：%s）", row[0]),
			})
			lineNumber++
			continue
		}
		layout := row[3]
		properties = append(properties, &domain.Property{
			Name:    row[0],
			Rent:    validator.ParseInt(row[1]),
			Address: row[2],
			Layout:  &layout,
			Status:  row[4],
		})
	}
	if len(errors) > 0 {
		return nil, errors, nil
	}

	if len(properties) > 0 {
		if err := u.repo.BulkCreate(properties); err != nil {
			return nil, nil, err
		}
	}
	result.ImportedCount = len(properties)
	return result, nil, nil

}

func (u *propertyUsecase) ExportProperties() ([]byte, error) {
	properties, err := u.repo.FindAll()
	if err != nil {
		return nil, err
	}
	if len(properties) == 0 {
		return nil, errors.New("no data")
	}
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	header := []string{"物件名", "賃料", "住所", "間取り", "ステータス"}
	writer.Write(header)

	for _, p := range properties {
		layout := ""
		if p.Layout != nil {
			layout = *p.Layout
		}
		writer.Write([]string{p.Name, fmt.Sprintf("%d", p.Rent), p.Address, layout, p.Status})
	}
	writer.Flush()

	return buf.Bytes(), nil
}
