package http

import (
	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserData(t *testing.T) {
	tests := []struct {
		name                 string
		id                   string
		setUCaseExpectations func(usecase *mocks.ProfileUsecase, user *domain.User)
		status               int
	}{
		{
			name: "GoodCase/Common",
			id:   "1",
			setUCaseExpectations: func(usecase *mocks.ProfileUsecase, user *domain.User) {
				usecase.On("GetUserData", mock.Anything).Return(*user, nil)
			},
			status: http.StatusOK,
		},
		{
			name: "BadCase/EmptyID",
			setUCaseExpectations: func(usecase *mocks.ProfileUsecase, user *domain.User) {
				usecase.On("GetUserData", mock.Anything).Return(*user, errors.New("empty id"))
			},
			status: http.StatusNotFound,
		},
		{
			id:   "hello",
			name: "BadCase/WrongIDFormat",
			setUCaseExpectations: func(usecase *mocks.ProfileUsecase, user *domain.User) {
				usecase.On("GetUserData", mock.Anything).Return(*user, errors.New("wrong id format"))
			},
			status: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			router := mux.NewRouter()
			mockUsecase := new(mocks.ProfileUsecase)
			var user domain.User
			test.setUCaseExpectations(mockUsecase, &user)

			req, err := http.NewRequest("GET", "/api/v1/profile/"+test.id, nil)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()

			NewProfileHandler(router, mockUsecase, bluemonday.UGCPolicy())

			handler := &ProfileHandler{
				ProfileUsecase: mockUsecase,
			}

			router.HandleFunc("/api/v1/profile/{id}", handler.GetUserData).Methods("GET")
			router.ServeHTTP(rec, req)

			assert.Equal(t, test.status, rec.Code)
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	tests := []struct {
		name                 string
		getBody              func() []byte
		setUCaseExpectations func(uCase *mocks.ProfileUsecase, updatedUser *domain.User)
		status               int
	}{
		{
			name: "GoodCase/Common",
			getBody: func() []byte {
				var user domain.User
				faker.FakeData(&user)
				user.Email = "chgvj@mail.ru"
				jsonBody, _ := json.Marshal(user)
				return jsonBody
			},
			setUCaseExpectations: func(uCase *mocks.ProfileUsecase, updatedUser *domain.User) {
				uCase.On("UploadImage", mock.Anything, mock.Anything).Return("path/to/image", nil)
				uCase.On("UpdateUser", mock.Anything).Return(*updatedUser, nil)
			},
			status: http.StatusOK,
		},
		{
			name: "BadCase/EmptyJson",
			getBody: func() []byte {
				return []byte("{}")
			},
			setUCaseExpectations: func(uCase *mocks.ProfileUsecase, updatedUser *domain.User) {
				uCase.On("UploadImage", mock.Anything, mock.Anything).Return("", nil).Maybe()
				uCase.On("UpdateUser", mock.Anything).Return(*updatedUser, domain.ErrWrongCredentials).Maybe()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/EmptyBody",
			getBody: func() []byte {
				return []byte("")
			},

			setUCaseExpectations: func(uCase *mocks.ProfileUsecase, updatedUser *domain.User) {
				uCase.On("UploadImage", mock.Anything, mock.Anything).Return("", nil).Maybe()
				uCase.On("UpdateUser", mock.Anything).Return(*updatedUser, domain.ErrBadRequest).Maybe()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/InvalidJson",
			getBody: func() []byte {
				return []byte("{043895uith,redfsvdf;vfdv4er")
			},
			setUCaseExpectations: func(uCase *mocks.ProfileUsecase, updatedUser *domain.User) {
				uCase.On("UploadImage", mock.Anything, mock.Anything).Return("", nil).Maybe()
				uCase.On("UpdateUser", mock.Anything).Return(*updatedUser, domain.ErrBadRequest).Maybe()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/BadImageData",
			getBody: func() []byte {
				var user domain.User
				faker.FakeData(&user)
				user.Email = "chgvj@mail.ru"
				user.ImageData = []byte("123")
				jsonBody, _ := json.Marshal(user)
				return jsonBody
			},
			setUCaseExpectations: func(uCase *mocks.ProfileUsecase, updatedUser *domain.User) {
				uCase.On("UploadImage", mock.Anything, mock.Anything).Return("", domain.ErrInternalServerError)
			},
			status: http.StatusInternalServerError,
		},
		{
			name: "BadCase/NoSuchUser",
			getBody: func() []byte {
				var user domain.User
				faker.FakeData(&user)
				user.Email = "chgvj@mail.ru"
				jsonBody, _ := json.Marshal(user)
				return jsonBody
			},
			setUCaseExpectations: func(uCase *mocks.ProfileUsecase, updatedUser *domain.User) {
				uCase.On("UploadImage", mock.Anything, mock.Anything).Return("path/to/image", nil)
				uCase.On("UpdateUser", mock.Anything).Return(*updatedUser, domain.ErrNotFound).Maybe()
			},
			status: http.StatusNotFound,
		},
		{
			name: "BadCase/FailedToUpdate",
			getBody: func() []byte {
				var user domain.User
				faker.FakeData(&user)
				user.Email = "chgvj@mail.ru"
				jsonBody, _ := json.Marshal(user)
				return jsonBody
			},
			setUCaseExpectations: func(uCase *mocks.ProfileUsecase, updatedUser *domain.User) {
				uCase.On("UploadImage", mock.Anything, mock.Anything).Return("path/to/image", nil)
				uCase.On("UpdateUser", mock.Anything).Return(*updatedUser, domain.ErrInternalServerError).Maybe()
			},
			status: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			req, err := http.NewRequest("POST", "/api/v1/profile/update", bytes.NewReader(test.getBody()))
			assert.NoError(t, err)

			mockUCase := new(mocks.ProfileUsecase)
			var updatedUser domain.User
			test.setUCaseExpectations(mockUCase, &updatedUser)

			rec := httptest.NewRecorder()
			handler := &ProfileHandler{
				ProfileUsecase: mockUCase,
			}

			handler.UpdateProfile(rec, req)

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
