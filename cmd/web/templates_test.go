package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHumanDate(t *testing.T) {
	tests := map[string]struct {
		tm   time.Time
		want string
	}{
		"UTC": {
			tm:   time.Date(2025, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2025 at 10:15",
		},
		"Empty": {
			tm:   time.Time{},
			want: "",
		},
		"CET": {
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2024 at 09:15",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			hd := humanDate(test.tm)

			assert.Equal(t, test.want, hd)
		})
	}
}
