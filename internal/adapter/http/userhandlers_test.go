package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/golang/mock/gomock"
	"github.com/k1nky/gophkeeper/internal/entity/user"
)

func (suite *adapterTestSuite) TestLogin() {
	type want struct {
		statusCode          int
		authorizationHeader string
	}
	tests := []struct {
		name        string
		payload     string
		want        want
		expectLogin []interface{}
	}{
		{
			name:        "Valid",
			payload:     `{"login": "user", "password": "pass"}`,
			want:        want{statusCode: http.StatusOK, authorizationHeader: "sometoken"},
			expectLogin: []interface{}{"sometoken", nil},
		},
		{
			name:        "Invalid json",
			payload:     `{"login": "user", `,
			want:        want{statusCode: http.StatusBadRequest, authorizationHeader: ""},
			expectLogin: []interface{}{},
		},
		{
			name:        "Invalid body format",
			payload:     `{"login": "user", "pass": ""} `,
			want:        want{statusCode: http.StatusBadRequest, authorizationHeader: ""},
			expectLogin: []interface{}{},
		},
		{
			name:        "Invalid credentials",
			payload:     `{"login": "user", "password": "invalid_password"}`,
			want:        want{statusCode: http.StatusUnauthorized, authorizationHeader: ""},
			expectLogin: []interface{}{"", user.ErrInvalidCredentials},
		},
		{
			name:        "Unexpected error",
			payload:     `{"login": "user", "password": "somepassword"}`,
			want:        want{statusCode: http.StatusInternalServerError, authorizationHeader: ""},
			expectLogin: []interface{}{"", errors.New("unexpected error")},
		},
	}
	a := &Adapter{
		auth: suite.authService,
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.payload))
		if len(tt.expectLogin) > 0 {
			suite.authService.EXPECT().Login(gomock.Any(), gomock.Any()).Return(tt.expectLogin...)
		}
		a.Login(w, r)
		suite.Equal(tt.want.statusCode, w.Code)
	}
}
