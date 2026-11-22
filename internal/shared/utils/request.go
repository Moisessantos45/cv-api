package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type QueryInt interface {
	~int64 | ~uint64 | ~int | ~uint
}

var URLRegex = regexp.MustCompile(`^https?:\/\/.+\..+`)

var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

var PhoneRegex = regexp.MustCompile(`^\+?[0-9\s\-\(\)]{7,20}$`)

func ValidateURL(link string, fieldName string) error {
	if strings.TrimSpace(link) == "" {
		return nil
	}

	if !URLRegex.MatchString(link) {
		return fmt.Errorf("el campo %s tiene un formato de URL inválido", fieldName)
	}

	_, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("el campo %s tiene un formato de URL inválido: %s", fieldName, err.Error())
	}

	return nil
}

func ValidateJSONFormat(data string, fieldName string) error {
	if strings.TrimSpace(data) == "" {
		return fmt.Errorf("el campo %s es obligatorio y debe tener un formato JSON válido", fieldName)
	}

	var js json.RawMessage
	if err := json.Unmarshal([]byte(data), &js); err != nil {
		return fmt.Errorf("el campo %s tiene un formato JSON inválido: %s", fieldName, err.Error())
	}

	return nil
}

func ValidateJSONArray(data string, fieldName string) error {
	if strings.TrimSpace(data) == "" {
		return fmt.Errorf("el campo %s es obligatorio y debe ser un array JSON válido", fieldName)
	}

	var arr []any
	if err := json.Unmarshal([]byte(data), &arr); err != nil {
		return fmt.Errorf("el campo %s debe ser un array JSON válido: %s", fieldName, err.Error())
	}

	return nil
}

func ValidateJSONObject(data string, fieldName string) error {
	if strings.TrimSpace(data) == "" {
		return fmt.Errorf("el campo %s es obligatorio y debe ser un objeto JSON válido", fieldName)
	}

	var obj map[string]any
	if err := json.Unmarshal([]byte(data), &obj); err != nil {
		return fmt.Errorf("el campo %s debe ser un objeto JSON válido: %s", fieldName, err.Error())
	}

	return nil
}

func ValidateParamsId(c *gin.Context, params string) (uint64, error) {
	newParams := params
	if len(strings.TrimSpace(params)) == 0 {
		newParams = "id"
	}

	idStr := c.Param(newParams)
	if len(strings.TrimSpace(idStr)) == 0 {
		return 0, errors.New("ID de producto no proporcionado")
	}

	id, newErr := strconv.ParseUint(idStr, 10, 64)
	if newErr != nil {
		return 0, errors.New("ID de producto inválido: " + newErr.Error())
	}

	return id, nil
}

func ValidateParamsQuery[T QueryInt](c *gin.Context, paramName string) (T, error) {
	paramStr := c.DefaultQuery(paramName, "1")
	if len(strings.TrimSpace(paramStr)) == 0 {
		var zero T
		return zero, fmt.Errorf("parámetro de consulta %q no proporcionado", paramName)
	}

	var zero T

	switch any(zero).(type) {
	case int64:
		v, err := strconv.ParseInt(paramStr, 10, 64)
		if err != nil {
			return zero, fmt.Errorf("parámetro de consulta %q inválido: %w", paramName, err)
		}
		return T(v), nil

	case uint64:
		v, err := strconv.ParseUint(paramStr, 10, 64)
		if err != nil {
			return zero, fmt.Errorf("parámetro de consulta %q inválido: %w", paramName, err)
		}
		return T(v), nil
	case int:
		v, err := strconv.Atoi(paramStr)
		if err != nil {
			return zero, fmt.Errorf("parámetro de consulta %q inválido: %w", paramName, err)
		}
		return T(v), nil
	case uint:
		v, err := strconv.ParseUint(paramStr, 10, 64)
		if err != nil {
			return zero, fmt.Errorf("parámetro de consulta %q inválido: %w", paramName, err)
		}
		return T(v), nil
	default:
		return zero, fmt.Errorf("tipo no soportado")
	}
}

func ExtractedParamsJwt(c *gin.Context) (string, uint64, error) {
	userIDIface, exists := c.Get("userID")
	if !exists {
		return "", 0, fmt.Errorf("userID no existe en el contexto")
	}

	userIDStr, ok := userIDIface.(string)
	if !ok {
		return "", 0, fmt.Errorf("userID no es string")
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		log.Printf("Error al convertir userID '%s' a uint64: %v", userIDStr, err)
		return "", 0, fmt.Errorf("userID inválido")
	}

	tokenIface, exists := c.Get("token")
	if !exists {
		return "", 0, fmt.Errorf("token no existe en el contexto")
	}

	tokenStr, ok := tokenIface.(string)
	if !ok {
		return "", 0, fmt.Errorf("token no es string")
	}

	return tokenStr, userID, nil
}

func ValidateQueryPagination(c *gin.Context) (int, int, string, error) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	searchQuery := c.DefaultQuery("search", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 0, 0, "", errors.New("Parámetro 'page' inválido: debe ser un número entero positivo")
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		return 0, 0, "", errors.New("Parámetro 'page_size' inválido: debe ser un número entero positivo")
	}

	return page, pageSize, searchQuery, nil
}
