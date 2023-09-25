package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	_http "2023_2_Holi/auth/delivery/http"
	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name                 string
		getBody              func() []byte
		setUCaseExpectations func(session *domain.Session, uCase *mocks.AuthUsecase)
		status               int
		wantCookie           bool
	}{
		{
			name: "GoodCase/Common",
			getBody: func() []byte {
				var creds domain.Credentials
				faker.FakeData(&creds)
				jsonBody, _ := json.Marshal(creds)
				return jsonBody
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				err := faker.FakeData(session)
				assert.NoError(t, err)

				uCase.On("Login", mock.Anything).Return(*session, nil)
			},
			status:     http.StatusOK,
			wantCookie: true,
		},
		{
			name: "BadCase/EmptyCredentials",
			getBody: func() []byte {
				jsonBody, _ := json.Marshal(domain.Credentials{})
				return jsonBody
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				*session = domain.Session{}
				uCase.On("Login", mock.Anything).Return(*session, domain.ErrWrongCredentials).Maybe()
			},
			status: http.StatusForbidden,
		},
		{
			name: "BadCase/InvalidJson",
			getBody: func() []byte {
				return []byte(`{ "password":"3490rjuv", name: rszdxtfcyguhj }`)
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				*session = domain.Session{}
				uCase.On("Login", mock.Anything).Return(*session, domain.ErrWrongCredentials).Maybe()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/EmptyJson",
			getBody: func() []byte {
				return []byte(`{}`)
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				*session = domain.Session{}
				uCase.On("Login", mock.Anything).Return(*session, domain.ErrWrongCredentials).Maybe()
			},
			status: http.StatusForbidden,
		},
		{
			name: "BadCase/NoBody",
			getBody: func() []byte {
				return []byte{}
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				*session = domain.Session{}
				uCase.On("Login", mock.Anything).Return(*session, domain.ErrWrongCredentials).Maybe()
			},
			status: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(test.getBody()))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			var mockSession domain.Session
			mockUCase := new(mocks.AuthUsecase)
			test.setUCaseExpectations(&mockSession, mockUCase)

			rec := httptest.NewRecorder()
			handler := &_http.AuthHandler{
				AuthUsecase: mockUCase,
			}

			handler.Login(rec, req)

			assert.Equal(t, test.status, rec.Code)
			mockUCase.AssertExpectations(t)

			cookies := rec.Result().Cookies()
			assert.NotEqual(t, len(cookies) == 0, test.wantCookie)

			if test.wantCookie {
				var sessionCookie *http.Cookie
				for _, cookie := range cookies {
					if cookie.Name == "session_token" {
						sessionCookie = cookie
						break
					}
				}
				assert.NotNil(t, sessionCookie)
				assert.Equal(t, mockSession.Token, sessionCookie.Value)

				expectedExpires := mockSession.ExpiresAt
				assert.WithinDuration(t, expectedExpires, sessionCookie.Expires, time.Second)
			}
		})
	}
}
