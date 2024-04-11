package jwt

import (
	"github.com/stretchr/testify/assert"
	"src/internal/models"
	"testing"
	"time"
)

func TestTokenProvider(t *testing.T) {
	testTable := []struct {
		name        string
		exp         time.Duration
		user        *models.User
		role        string
		expectedErr error
		isValid     bool
	}{
		{
			name: "Test admin",
			exp:  time.Hour,
			user: &models.User{
				Id:       1,
				Name:     "test_name",
				Password: "test_pass",
				Role:     "admin",
				Email:    "test_email",
			},
			role:        "admin",
			expectedErr: nil,
			isValid:     true,
		},
		{
			name: "Test user",
			exp:  time.Hour,
			user: &models.User{
				Id:       1,
				Name:     "test_name",
				Password: "test_pass",
				Role:     "user",
				Email:    "test_email",
			},
			role:        "user",
			expectedErr: nil,
			isValid:     true,
		},
		{
			name: "Test expired",
			exp:  time.Duration(-10),
			user: &models.User{
				Id:       1,
				Name:     "test_name",
				Password: "test_pass",
				Role:     "user",
				Email:    "test_email",
			},
			role:        "user",
			expectedErr: nil,
			isValid:     false,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tp := NewTokenProvider("test", tc.exp)

			tok, err := tp.GenerateToken(tc.user)
			val, _ := tp.IsTokenValid(tok)
			assert.Equal(t, tc.isValid, val)
			if tc.expectedErr == nil {
				assert.Nil(t, err)
				assert.NotNil(t, tok)

				role, err2 := tp.GetRole(tok)
				if tc.isValid {
					assert.Equal(t, role, tc.role)
					assert.Nil(t, err2)
				} else {
					assert.Equal(t, role, "")
					assert.Error(t, err2)
				}
			} else {
				assert.Nil(t, tok)
				assert.Error(t, err)
			}

		})
	}
}
