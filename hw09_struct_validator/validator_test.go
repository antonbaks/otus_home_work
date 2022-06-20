package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "err",
				Name:   "Test",
				Age:    10,
				meta:   nil,
				Email:  "err",
				Role:   "root",
				Phones: []string{"1234567890", "12345678"},
			},
			expectedErr: ValidationErrors{
				ValidationError{"ID", ErrStrLen},
				ValidationError{"Age", ErrIntMin},
				ValidationError{"Email", ErrStrRegexp},
				ValidationError{"Role", ErrStrIn},
				ValidationError{"Phones", ErrStrLen},
				ValidationError{"Phones", ErrStrLen},
			},
		},
		{App{"1234"}, ValidationErrors{ValidationError{"Version", ErrStrLen}}},
		{Response{304, ""}, ValidationErrors{ValidationError{"Code", ErrIntIn}}},
		{
			in: User{
				ID:     "11f47ad5-7b73-42c0-abae-878b1e16adee",
				Name:   "Test",
				Age:    20,
				meta:   nil,
				Email:  "test@test.com",
				Role:   "admin",
				Phones: []string{"89519111511", "89519111511"},
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			require.Equal(t, tt.expectedErr, Validate(tt.in))
		})
	}
}
