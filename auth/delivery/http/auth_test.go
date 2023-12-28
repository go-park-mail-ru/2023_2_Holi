package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"
)

const userID = "1"

func TestLogin(t *testing.T) {
	tests := []struct {
		name                 string
		getBody              func() []byte
		setUCaseExpectations func(session *domain.Session, uCase *mocks.AuthUsecase)
		status               int
		wantCookie           bool
		auth                 string
		setAuth              func(r *http.Request, uCase *mocks.AuthUsecase, session *domain.Session)
	}{
		//{
		//	name: "GoodCase/Common",
		//	getBody: func() []byte {
		//		var creds domain.Credentials
		//		faker.FakeData(&creds.Password)
		//		creds.Email = "ferfg@fsf.ru"
		//		creds.Password = []byte{61,
		//			73, 76, 31}
		//		jsonBody, _ := json.Marshal(creds)
		//		return jsonBody
		//	},
		//	setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
		//		err := faker.FakeData(session)
		//		assert.NoError(t, err)
		//		session.ExpiresAt = time.Now().Add(24 * time.Hour)
		//
		//		uCase.On("Login", mock.Anything).Return(*session, 1, nil)
		//	},
		//	status:     http.StatusOK,
		//	wantCookie: true,
		//},
		{
			name: "BadCase/EmptyCredentials",
			getBody: func() []byte {
				jsonBody, _ := json.Marshal(domain.Credentials{})
				return jsonBody
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				*session = domain.Session{}
				uCase.On("Login", mock.Anything).Return(*session, 0, domain.ErrWrongCredentials).Maybe()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/InvalidJson",
			getBody: func() []byte {
				return []byte(`{ "password":"3490rjuv", email: rszdxtfcyguhj@sgf.ru }`)
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				*session = domain.Session{}
				uCase.On("Login", mock.Anything).Return(*session, 0, domain.ErrWrongCredentials).Maybe()
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
				uCase.On("Login", mock.Anything).Return(*session, 0, domain.ErrWrongCredentials).Maybe()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/NoBody",
			getBody: func() []byte {
				return []byte{}
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				*session = domain.Session{}
				uCase.On("Login", mock.Anything).Return(*session, 0, domain.ErrWrongCredentials).Maybe()
			},
			status: http.StatusBadRequest,
		},
		//{
		//	name: "BadCase/AlreadyAuthorized",
		//	getBody: func() []byte {
		//		var creds domain.Credentials
		//		faker.FakeData(&creds.Password)
		//		creds.Email = "ferfg@fsf.ru"
		//		jsonBody, _ := json.Marshal(creds)
		//		return jsonBody
		//	},
		//	setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
		//		err := faker.FakeData(session)
		//		assert.NoError(t, err)
		//		session.ExpiresAt = time.Now().Add(24 * time.Hour)
		//
		//		uCase.On("Login", mock.Anything).Return(*session, 0, nil).Maybe()
		//	},
		//	status: http.StatusConflict,
		//	auth:   userID,
		//	setAuth: func(r *http.Request, uCase *mocks.AuthUsecase, session *domain.Session) {
		//		r.AddCookie(&http.Cookie{
		//			Name:     "session_token",
		//			Value:    session.Token,
		//			Expires:  session.ExpiresAt,
		//			Path:     "/",
		//			HttpOnly: true,
		//		})
		//		uCase.On("IsAuth", mock.Anything).Return(userID, nil)
		//	},
		//},
		//{
		//	name: "GoodCase/AlreadyAuthorizedExpiredCookie",
		//	getBody: func() []byte {
		//		var creds domain.Credentials
		//		faker.FakeData(&creds.Password)
		//		creds.Email = "ferfg@fsf.ru"
		//		jsonBody, _ := json.Marshal(creds)
		//		return jsonBody
		//	},
		//	setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
		//		err := faker.FakeData(session)
		//		assert.NoError(t, err)
		//		session.ExpiresAt = time.Now()
		//
		//		uCase.On("Login", mock.Anything).Return(*session, 1, nil).Maybe()
		//	},
		//	status:     http.StatusOK,
		//	wantCookie: true,
		//	auth:       userID,
		//	setAuth: func(r *http.Request, uCase *mocks.AuthUsecase, session *domain.Session) {
		//		r.AddCookie(&http.Cookie{
		//			Name:     "session_token",
		//			Value:    session.Token,
		//			Expires:  session.ExpiresAt,
		//			Path:     "/",
		//			HttpOnly: true,
		//		})
		//		uCase.On("IsAuth", mock.Anything).Return("", nil).Maybe()
		//	},
		//},
		//{
		//	name: "GoodCase/AlreadyAuthorizedWrongCookie",
		//	getBody: func() []byte {
		//		var creds domain.Credentials
		//		faker.FakeData(&creds.Password)
		//		creds.Email = "ferfg@fsf.ru"
		//		jsonBody, _ := json.Marshal(creds)
		//		return jsonBody
		//	},
		//	setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
		//		err := faker.FakeData(session)
		//		assert.NoError(t, err)
		//		session.ExpiresAt = time.Now().Add(24 * time.Hour)
		//
		//		uCase.On("Login", mock.Anything).Return(*session, 1, nil).Maybe()
		//	},
		//	status:     http.StatusOK,
		//	wantCookie: true,
		//	auth:       userID,
		//	setAuth: func(r *http.Request, uCase *mocks.AuthUsecase, session *domain.Session) {
		//		r.AddCookie(&http.Cookie{
		//			Name:     "fevk",
		//			Value:    session.Token,
		//			Expires:  session.ExpiresAt,
		//			Path:     "/",
		//			HttpOnly: true,
		//		})
		//		uCase.On("IsAuth", mock.Anything).Return("", nil).Maybe()
		//	},
		//},
		//{
		//	name: "BadCase/UserNotFound",
		//	getBody: func() []byte {
		//		var creds domain.Credentials
		//		faker.FakeData(&creds.Password)
		//		creds.Email = "ferfg@fsf.ru"
		//		jsonBody, _ := json.Marshal(creds)
		//		return jsonBody
		//	},
		//	setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
		//		*session = domain.Session{}
		//
		//		uCase.On("IsAuth", mock.Anything).Return("", nil).Maybe()
		//		uCase.On("Login", mock.Anything).Return(*session, 0, domain.ErrNotFound)
		//	},
		//	status: http.StatusNotFound,
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(test.getBody()))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			var mockSession domain.Session
			mockUCase := new(mocks.AuthUsecase)
			test.setUCaseExpectations(&mockSession, mockUCase)
			if test.auth != "" {
				test.setAuth(req, mockUCase, &mockSession)
			}
			rec := httptest.NewRecorder()
			NewAuthHandler(mux.NewRouter(), mockUCase)
			handler := &AuthHandler{
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
				uCase.On("IsAuth", mock.Anything).Return(userID, nil)
				uCase.On("Logout", mock.Anything).Return(nil)
			},
			status:     http.StatusNoContent,
			wantCookie: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			req, err := http.NewRequest("POST", "/api/v1/auth/logout", strings.NewReader(""))
			assert.NoError(t, err)
			req.AddCookie(test.getCookie())

			mockUCase := new(mocks.AuthUsecase)
			test.setUCaseExpectations(mockUCase)

			rec := httptest.NewRecorder()
			handler := &AuthHandler{
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
		setUCaseExpectations func(uCase *mocks.AuthUsecase, session *domain.Session)
		status               int
		auth                 string
		setAuth              func(r *http.Request, uCase *mocks.AuthUsecase, session *domain.Session)
	}{
		//{
		//	name: "GoodCase/Common",
		//	getBody: func() []byte {
		//		var user domain.User
		//		faker.FakeData(&user)
		//		user.Email = "chgvj@mail.ru"
		//		jsonBody, _ := json.Marshal(user)
		//		return jsonBody
		//	},
		//	setUCaseExpectations: func(uCase *mocks.AuthUsecase, session *domain.Session) {
		//		uCase.On("Register", mock.Anything).Return(1, nil)
		//
		//		err := faker.FakeData(session)
		//		assert.NoError(t, err)
		//		uCase.On("Login", mock.Anything).Return(*session, 1, nil)
		//	},
		//	status: http.StatusOK,
		//},
		//{
		//	name: "BadCase/AlreadyRegistered",
		//	getBody: func() []byte {
		//		var user domain.User
		//		faker.FakeData(&user)
		//		user.Email = "chgvj@mail.ru"
		//		jsonBody, _ := json.Marshal(user)
		//		return jsonBody
		//	},
		//	setUCaseExpectations: func(uCase *mocks.AuthUsecase, session *domain.Session) {
		//		uCase.On("Register", mock.Anything).Return(0, domain.ErrAlreadyExists)
		//
		//		err := faker.FakeData(session)
		//		assert.NoError(t, err)
		//		uCase.On("Login", mock.Anything).Return(*session, 1, nil).Maybe()
		//	},
		//	status: http.StatusConflict,
		//},
		{
			name: "BadCase/EmptyJson",
			getBody: func() []byte {
				return []byte("{}")
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase, session *domain.Session) {
				uCase.On("Register", mock.Anything).Return(0, nil).Maybe()
				uCase.On("Login", mock.Anything).Return(*session, 1, domain.ErrWrongCredentials).Maybe()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/EmptyBody",
			getBody: func() []byte {
				return []byte("")
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase, session *domain.Session) {
				uCase.On("Register", mock.Anything).Return(0, nil).Maybe()
				uCase.On("Login", mock.Anything).Return(*session, 1, domain.ErrBadRequest).Maybe()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/InvalidJson",
			getBody: func() []byte {
				return []byte("{043895uith,redfsvdf;vfdv4er")
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase, session *domain.Session) {
				uCase.On("Register", mock.Anything).Return(0, nil).Maybe()
				uCase.On("Login", mock.Anything).Return(*session, 1, domain.ErrBadRequest).Maybe()
			},
			status: http.StatusBadRequest,
		},
		//{
		//	name: "BadCase/AlreadyAuthorized",
		//	getBody: func() []byte {
		//		var user domain.User
		//		faker.FakeData(&user)
		//		user.Email = "chgvj@mail.ru"
		//		jsonBody, _ := json.Marshal(user)
		//		return jsonBody
		//	},
		//	setUCaseExpectations: func(uCase *mocks.AuthUsecase, session *domain.Session) {
		//		uCase.On("Register", mock.Anything).Return(0, errors.New("some")).Maybe()
		//
		//		err := faker.FakeData(session)
		//		session.ExpiresAt = time.Now().Add(24 * time.Hour)
		//		assert.NoError(t, err)
		//		uCase.On("Login", mock.Anything).Return(*session, 0, nil).Maybe()
		//	},
		//	status: http.StatusConflict,
		//	auth:   userID,
		//	setAuth: func(r *http.Request, uCase *mocks.AuthUsecase, session *domain.Session) {
		//		r.AddCookie(&http.Cookie{
		//			Name:     "session_token",
		//			Value:    session.Token,
		//			Expires:  session.ExpiresAt,
		//			Path:     "/",
		//			HttpOnly: true,
		//		})
		//		uCase.On("IsAuth", mock.Anything).Return(userID, nil)
		//	},
		//},
		//{
		//	name: "GoodCase/AlreadyAuthorizedExpiredCookie",
		//	getBody: func() []byte {
		//		var user domain.User
		//		faker.FakeData(&user)
		//		user.Email = "chgvj@mail.ru"
		//		jsonBody, _ := json.Marshal(user)
		//		return jsonBody
		//	},
		//	setUCaseExpectations: func(uCase *mocks.AuthUsecase, session *domain.Session) {
		//		uCase.On("Register", mock.Anything).Return(1, nil)
		//
		//		err := faker.FakeData(session)
		//		session.ExpiresAt = time.Now().Add(24 * time.Hour)
		//		assert.NoError(t, err)
		//		uCase.On("Login", mock.Anything).Return(*session, 1, nil)
		//	},
		//	status: http.StatusOK,
		//	auth:   userID,
		//	setAuth: func(r *http.Request, uCase *mocks.AuthUsecase, session *domain.Session) {
		//		r.AddCookie(&http.Cookie{
		//			Name:     "session_token",
		//			Value:    session.Token,
		//			Expires:  time.Now(),
		//			Path:     "/",
		//			HttpOnly: true,
		//		})
		//		uCase.On("IsAuth", mock.Anything).Return("", nil).Maybe()
		//	},
		//},
		//{
		//	name: "GoodCase/AlreadyAuthorizedWrongCookie",
		//	getBody: func() []byte {
		//		var user domain.User
		//		faker.FakeData(&user)
		//		user.Email = "chgvj@mail.ru"
		//		jsonBody, _ := json.Marshal(user)
		//		return jsonBody
		//	},
		//	setUCaseExpectations: func(uCase *mocks.AuthUsecase, session *domain.Session) {
		//		uCase.On("Register", mock.Anything).Return(1, nil)
		//
		//		err := faker.FakeData(session)
		//		session.ExpiresAt = time.Now().Add(24 * time.Hour)
		//		assert.NoError(t, err)
		//		uCase.On("Login", mock.Anything).Return(*session, 1, nil)
		//	},
		//	status: http.StatusOK,
		//	auth:   userID,
		//	setAuth: func(r *http.Request, uCase *mocks.AuthUsecase, session *domain.Session) {
		//		r.AddCookie(&http.Cookie{
		//			Name:     "fevk",
		//			Value:    session.Token,
		//			Expires:  session.ExpiresAt,
		//			Path:     "/",
		//			HttpOnly: true,
		//		})
		//		uCase.On("IsAuth", mock.Anything).Return("", nil).Maybe()
		//	},
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			req, err := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(test.getBody()))
			assert.NoError(t, err)

			mockUCase := new(mocks.AuthUsecase)
			var mockSession domain.Session
			test.setUCaseExpectations(mockUCase, &mockSession)
			if test.auth != "" {
				test.setAuth(req, mockUCase, &mockSession)
			}
			rec := httptest.NewRecorder()
			handler := &AuthHandler{
				AuthUsecase: mockUCase,
			}

			handler.Register(rec, req)

			assert.Equal(t, test.status, rec.Code)
			mockUCase.AssertExpectations(t)

			var result *domain.Response
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
