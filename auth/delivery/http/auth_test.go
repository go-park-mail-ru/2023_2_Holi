package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

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

func TestLogout(t *testing.T) {
	tests := []struct {
		name                 string
		getCookie            func() *http.Cookie
		setUCaseExpectations func(uCase *mocks.AuthUsecase)
		status               int
		wantCookie           bool
	}{
		{
			name: "GoodCase/Common",
			getCookie: func() *http.Cookie {
				return &http.Cookie{
					Name:    "session_token",
					Value:   uuid.NewString(),
					Expires: time.Now().Add(24 * time.Hour),
				}
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase) {
				uCase.On("Logout", mock.Anything).Return(nil)
			},
			status:     http.StatusOK,
			wantCookie: true,
		},
		{
			name: "BadCase/EmptyCookiePayload",
			getCookie: func() *http.Cookie {
				return &http.Cookie{
					Name:    "session_token",
					Value:   "",
					Expires: time.Time{},
				}
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase) {
				uCase.On("Logout", mock.Anything).Return(nil).Maybe()
			},
			status:     http.StatusUnauthorized,
			wantCookie: false,
		},
		{
			name: "BadCase/UnsuitableName",
			getCookie: func() *http.Cookie {
				return &http.Cookie{
					Name:    "rdtcfyvgubhj",
					Value:   uuid.NewString(),
					Expires: time.Now().Add(24 * time.Hour),
				}
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase) {
				uCase.On("Logout", mock.Anything).Return(nil).Maybe()
			},
			status:     http.StatusUnauthorized,
			wantCookie: false,
		},
		{
			name: "BadCase/EmptyCookie",
			getCookie: func() *http.Cookie {
				return &http.Cookie{}
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase) {
				uCase.On("Logout", mock.Anything).Return(nil).Maybe()
			},
			status:     http.StatusUnauthorized,
			wantCookie: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest("POST", "/api/v1/auth/logout", strings.NewReader(""))
			assert.NoError(t, err)
			req.AddCookie(test.getCookie())

			mockUCase := new(mocks.AuthUsecase)
			test.setUCaseExpectations(mockUCase)

			rec := httptest.NewRecorder()
			handler := &_http.AuthHandler{
				AuthUsecase: mockUCase,
			}

			handler.Logout(rec, req)

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
				assert.Empty(t, sessionCookie.Value)

				assert.WithinDuration(t, time.Now(), sessionCookie.Expires, 10*time.Second)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name                 string
		getBody              func() []byte
		setUCaseExpectations func(uCase *mocks.AuthUsecase)
		status               int
	}{
		{
			name: "GoodCase/Common",
			getBody: func() []byte {
				var user domain.User
				faker.FakeData(&user)
				jsonBody, _ := json.Marshal(user)
				return jsonBody
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase) {
				uCase.On("Register", mock.Anything).Return(1, nil)
			},
			status: http.StatusOK,
		},
		{
			name: "BadCase/EmptyJson",
			getBody: func() []byte {
				return []byte("{}")
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase) {
				uCase.On("Register", mock.Anything).Return(0, nil).Maybe()
			},
			status: http.StatusForbidden,
		},
		{
			name: "BadCase/EmptyBody",
			getBody: func() []byte {
				return []byte("")
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase) {
				uCase.On("Register", mock.Anything).Return(0, nil).Maybe()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/InvalidJson",
			getBody: func() []byte {
				return []byte("{043895uith,redfsvdf;vfdv4er")
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase) {
				uCase.On("Register", mock.Anything).Return(0, nil).Maybe()
			},
			status: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			req, err := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(test.getBody()))
			assert.NoError(t, err)

			mockUCase := new(mocks.AuthUsecase)
			test.setUCaseExpectations(mockUCase)

			rec := httptest.NewRecorder()
			handler := &_http.AuthHandler{
				AuthUsecase: mockUCase,
			}

			handler.Register(rec, req)

			assert.Equal(t, test.status, rec.Code)
			mockUCase.AssertExpectations(t)

			var result *_http.Result
			err = json.NewDecoder(rec.Result().Body).Decode(&result)
			assert.NoError(t, err)

			if test.status < 300 {
				assert.NotEmpty(t, result.Body)
				assert.Empty(t, result.Err)
			} else {
				assert.Empty(t, result.Body)
				assert.NotEmpty(t, result.Err)
			}
		})
	}
}
