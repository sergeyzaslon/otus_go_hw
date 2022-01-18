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
		Phones []string `validate:"len:12"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	ResponseUser struct {
		Code int  `validate:"in:200,404,500"`
		User User `validate:"nested"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			"Hello", // Простые типы не валидируем и не ругаемся на них
			nil,
		},
		{
			User{
				ID:     "2265e743-32ba-4264-8b33-5908afd2978b",
				Name:   "Test Name",
				Age:    25,
				Email:  "sergeyzaslon@gmail.com",
				Role:   "admin",
				Phones: []string{"380699911337"},
				meta:   []byte("{\"is_test\": true}"),
			},
			nil,
		},
		{
			User{
				ID:     "5908afd2978b", // len != 36
				Name:   "Test",
				Age:    16,                                        // < 18
				Email:  "some-invalid-email@",                     // invalid email
				Role:   "invalid_role",                            // not in enum
				Phones: []string{"380699911337", "+380699911337"}, // 2nd phone is invalid: len > 11
				meta:   []byte("{\"is_test\": true}"),
			},
			ValidationErrors{
				ValidationError{"ID", ErrStrLen},
				ValidationError{"Age", ErrNumRange},
				ValidationError{"Email", ErrRegexp},
				ValidationError{"Role", ErrStrEnum},
				ValidationError{"Phones.1", ErrStrLen},
			},
		},
		{
			App{"1.0.2"},
			nil,
		},
		{
			App{"1.0.2.6"},
			ValidationErrors{
				ValidationError{"Version", ErrStrLen},
			},
		},
		{
			Token{[]byte{0x01}, []byte{0x02}, []byte{0x03}},
			nil,
		},
		{
			Response{200, "OK"},
			nil,
		},
		{
			Response{404, "Not Found"},
			nil,
		},
		{
			Response{500, "Internal Server Error"},
			nil,
		},
		{
			Response{503, "Service Temparary Unavailable"},
			ValidationErrors{
				ValidationError{"Code", ErrStrEnum},
			},
		},
		{
			ResponseUser{503, User{ // 503 = invalid
				ID:     "2265e743", // invalid!
				Name:   "Test Name",
				Age:    25,
				Email:  "sergeyzaslon@gmail.com",
				Role:   "admin",
				Phones: []string{"380699911337"},
				meta:   []byte("{\"is_test\": true}"),
			}},
			ValidationErrors{
				ValidationError{"Code", ErrStrEnum},
				ValidationError{"ID", ErrStrLen},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr != nil {
				require.Equal(t, tt.expectedErr, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}
