package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid user",
			user: User{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			user: User{
				Name:  "",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "empty email",
			user: User{
				Name:  "John Doe",
				Email: "",
				Age:   25,
			},
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "negative age",
			user: User{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   -1,
			},
			wantErr: true,
			errMsg:  "age must be between 0 and 150",
		},
		{
			name: "age too high",
			user: User{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   151,
			},
			wantErr: true,
			errMsg:  "age must be between 0 and 150",
		},
		{
			name: "minimum age",
			user: User{
				Name:  "Baby Doe",
				Email: "baby@example.com",
				Age:   0,
			},
			wantErr: false,
		},
		{
			name: "maximum age",
			user: User{
				Name:  "Elder Doe",
				Email: "elder@example.com",
				Age:   150,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
