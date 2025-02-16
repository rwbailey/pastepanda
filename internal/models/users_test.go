package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserModelExists(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := map[string]struct {
		name   string
		userID int
		want   bool
	}{
		"Valid ID": {
			userID: 1,
			want:   true,
		},
		"Zero ID": {
			userID: 0,
			want:   false,
		},
		"Non-existent ID": {
			userID: 2,
			want:   false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db := newTestDB(t)

			m := UserModel{db}

			exists, err := m.Exists(test.userID)

			assert.Equal(t, test.want, exists)
			assert.NoError(t, err)
		})
	}
}
