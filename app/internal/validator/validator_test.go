package validator_test

import (
	"strings"
	"testing"

	"github.com/alfattd/category-service/internal/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateCategoryName_Valid(t *testing.T) {
	validNames := []string{
		"Electronics",
		"Books & Music",
		"Home Appliances",
		"Sports",
		"A",
		strings.Repeat("a", 20),
	}

	validNames = []string{
		"Electronics",
		"Home Appliances",
		"Sports",
		"A",
		"Category 123",
		strings.Repeat("a", 20),
	}

	for _, name := range validNames {
		t.Run(name, func(t *testing.T) {
			errs := validator.CategoryNameValidator(name)
			assert.Nil(t, errs, "expected no errors for name: %q", name)
		})
	}
}

func TestValidateCategoryName_EmptyName(t *testing.T) {
	errs := validator.CategoryNameValidator("")
	require.NotNil(t, errs)
	assert.Contains(t, errs.Messages, "name is required")
	assert.Len(t, errs.Messages, 1)
}

func TestValidateCategoryName_WhitespaceOnly(t *testing.T) {
	cases := []string{"   ", "\t", "\n", "  \t  "}
	for _, name := range cases {
		t.Run("whitespace:"+name, func(t *testing.T) {
			errs := validator.CategoryNameValidator(name)
			require.NotNil(t, errs)
			assert.True(t, errs.HasErrors())
		})
	}
}

func TestValidateCategoryName_TooLong(t *testing.T) {
	name := strings.Repeat("a", 21)
	errs := validator.CategoryNameValidator(name)
	require.NotNil(t, errs)
	assert.Contains(t, errs.Messages, "name must not exceed 20 characters")
}

func TestValidateCategoryName_ExactlyMaxLength(t *testing.T) {
	name := strings.Repeat("a", 20)
	errs := validator.CategoryNameValidator(name)
	assert.Nil(t, errs)
}

func TestValidateCategoryName_ForbiddenCharacters(t *testing.T) {
	forbidden := []string{
		"<script>", "name>value", "a;b", "x&y",
		`a\b`, "a/b", "{name}", "(name)", "[name]",
		`name"value`, "name'value",
	}

	for _, name := range forbidden {
		t.Run(name, func(t *testing.T) {
			errs := validator.CategoryNameValidator(name)
			require.NotNil(t, errs)
			assert.True(t, errs.HasErrors())
		})
	}
}

func TestValidateCategoryName_MultipleErrors(t *testing.T) {
	name := strings.Repeat("a<", 51)
	errs := validator.CategoryNameValidator(name)
	require.NotNil(t, errs)
	assert.GreaterOrEqual(t, len(errs.Messages), 2, "expected multiple validation errors")
}

func TestValidationErrors_Error(t *testing.T) {
	errs := &validator.ErrorsValidator{}
	errs.Add("name is required")
	errs.Add("name must not exceed 20 characters")

	assert.Equal(t, "name is required; name must not exceed 20 characters", errs.Error())
}

func TestValidationErrors_HasErrors(t *testing.T) {
	errs := &validator.ErrorsValidator{}
	assert.False(t, errs.HasErrors())

	errs.Add("some error")
	assert.True(t, errs.HasErrors())
}
