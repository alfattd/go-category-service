package validator_test

import (
	"strings"
	"testing"

	"github.com/alfattd/category-service/internal/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── CategoryNameValidator ────────────────────────────────────────────────────

func TestValidateCategoryName_Valid(t *testing.T) {
	validNames := []string{
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

// ─── CategoryIDValidator ──────────────────────────────────────────────────────

func TestValidateCategoryID_Valid(t *testing.T) {
	validIDs := []string{
		"abc-123",
		"550e8400-e29b-41d4-a716-446655440000",
		"1",
	}

	for _, id := range validIDs {
		t.Run(id, func(t *testing.T) {
			errs := validator.CategoryIDValidator(id)
			assert.Nil(t, errs, "expected no errors for id: %q", id)
		})
	}
}

func TestValidateCategoryID_EmptyID(t *testing.T) {
	errs := validator.CategoryIDValidator("")
	require.NotNil(t, errs)
	assert.Contains(t, errs.Messages, "id is required")
	assert.Len(t, errs.Messages, 1)
}

func TestValidateCategoryID_WhitespaceOnly(t *testing.T) {
	cases := []string{"   ", "\t", "\n"}
	for _, id := range cases {
		t.Run("whitespace", func(t *testing.T) {
			errs := validator.CategoryIDValidator(id)
			require.NotNil(t, errs)
			assert.Contains(t, errs.Messages, "id must not be blank")
		})
	}
}

// ─── ErrorsValidator ──────────────────────────────────────────────────────────

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
