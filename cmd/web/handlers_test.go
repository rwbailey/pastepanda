package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	ping(rr, r)

	rs := rr.Result()

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	require.NoError(t, err)

	body = bytes.TrimSpace(body)

	assert.Equal(t, "OK", string(body))
}

func TestPingEndToEnd(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "OK", body)
}

func TestPasteView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := map[string]struct {
		urlPath  string
		wantCode int
		wantBody string
	}{
		"Valid ID": {
			urlPath:  "/paste/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		"Non-existent ID": {
			urlPath:  "/paste/view/2",
			wantCode: http.StatusNotFound,
		},
		"Negative ID": {
			urlPath:  "/paste/view/-1",
			wantCode: http.StatusNotFound,
		},
		"Decimal ID": {
			urlPath:  "/paste/view/1.23",
			wantCode: http.StatusNotFound,
		},
		"String ID": {
			urlPath:  "/paste/view/foo",
			wantCode: http.StatusNotFound,
		},
		"Empty ID": {
			urlPath:  "/paste/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			code, _, body := ts.get(t, test.urlPath)

			assert.Equal(t, test.wantCode, code)

			if test.wantBody != "" {
				assert.Contains(t, body, test.wantBody)
			}
		})
	}
}
