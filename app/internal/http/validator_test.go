package http

import (
	"net/http"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"gte=0,lte=130"`
}

func TestCustomValidator_Validate_Success(t *testing.T) {
	cv := &CustomValidator{
		validator: validator.New(),
	}

	testData := TestStruct{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	err := cv.Validate(testData)

	assert.NoError(t, err)
}

func TestCustomValidator_Validate_RequiredFieldMissing(t *testing.T) {
	cv := &CustomValidator{
		validator: validator.New(),
	}

	testData := TestStruct{
		Name:  "",
		Email: "john@example.com",
		Age:   30,
	}

	err := cv.Validate(testData)

	assert.Error(t, err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestCustomValidator_Validate_InvalidEmail(t *testing.T) {
	cv := &CustomValidator{
		validator: validator.New(),
	}

	testData := TestStruct{
		Name:  "John Doe",
		Email: "invalid-email",
		Age:   30,
	}

	err := cv.Validate(testData)

	assert.Error(t, err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestCustomValidator_Validate_AgeOutOfRange(t *testing.T) {
	cv := &CustomValidator{
		validator: validator.New(),
	}

	testData := TestStruct{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   150,
	}

	err := cv.Validate(testData)

	assert.Error(t, err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestCustomValidator_Validate_MultipleErrors(t *testing.T) {
	cv := &CustomValidator{
		validator: validator.New(),
	}

	testData := TestStruct{
		Name:  "",
		Email: "invalid-email",
		Age:   -5,
	}

	err := cv.Validate(testData)

	assert.Error(t, err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestCustomValidator_Validate_ValidEdgeCases(t *testing.T) {
	cv := &CustomValidator{
		validator: validator.New(),
	}

	testCases := []TestStruct{
		{Name: "A", Email: "a@b.c", Age: 0},
		{Name: "Very Long Name", Email: "test@example.com", Age: 130},
	}

	for _, tc := range testCases {
		err := cv.Validate(tc)
		assert.NoError(t, err, "Should validate successfully: %+v", tc)
	}
}
