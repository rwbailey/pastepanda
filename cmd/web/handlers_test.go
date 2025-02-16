package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	validCSRFToken := extractCSRFToken(t, body)

	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = "<form action='/user/signup' method='POST' novalidate>"
	)

	tests := map[string]struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}{
		"Valid submission": {
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusSeeOther,
		},
		"Invalid CSRF Token": {
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "wrongToken",
			wantCode:     http.StatusBadRequest,
		},
		"Empty name": {
			userName:     "",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		"Empty email": {
			userName:     validName,
			userEmail:    "",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		"Empty password": {
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		"Invalid email": {
			userName:     validName,
			userEmail:    "bob@example.",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		"Short password": {
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "pa$$",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		"Duplicate email": {
			userName:     validName,
			userEmail:    "dupe@example.com",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", test.userName)
			form.Add("email", test.userEmail)
			form.Add("password", test.userPassword)
			form.Add("csrf_token", test.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			assert.Equal(t, code, test.wantCode)

			if test.wantFormTag != "" {
				assert.Contains(t, body, test.wantFormTag)
			}
		})
	}
}
