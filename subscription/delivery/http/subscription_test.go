package http

import (
	"2023_2_Holi/domain/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestUnSubscribe(t *testing.T) {
	tests := []struct {
		name       string
		subID      int
		expectErr  error
		statusCode int
	}{
		{
			name:       "Success",
			subID:      1,
			expectErr:  nil,
			statusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUsecase := new(mocks.SubsUsecase)
			handler := &SubsHandler{
				SubsUsecase: mockUsecase,
			}

			mockUsecase.On("UnSubscribe", test.subID).Return(test.expectErr)

			req, err := http.NewRequest("POST", "/v1/subs/unsub/"+strconv.Itoa(test.subID), nil)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/v1/subs/unsub/{id}", handler.UnSubscribe).Methods(http.MethodPost, http.MethodOptions)
			router.ServeHTTP(rec, req)

			assert.Equal(t, test.statusCode, rec.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

//func TestCheckSub(t *testing.T) {
//	tests := []struct {
//		name         string
//		subID        int
//		expectUpTo   time.Time
//		expectStatus string
//		expectErr    error
//		statusCode   int
//	}{
//		{
//			name:         "Success",
//			subID:        1,
//			expectUpTo:   time.Now(),
//			expectStatus: "active",
//			expectErr:    nil,
//			statusCode:   http.StatusOK,
//		},
//		// Добавьте другие тестовые случаи по мере необходимости
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			mockUsecase := new(mocks.SubsUsecase)
//			handler := &SubsHandler{
//				SubsUsecase: mockUsecase,
//			}
//
//			// Преобразование времени в Unix-формат
//			expectedUpToUnix := test.expectUpTo.Unix()
//
//			mockUsecase.On("CheckSub", test.subID).Return(expectedUpToUnix, test.expectStatus, test.expectErr)
//
//			req, err := http.NewRequest("GET", "/v1/subs/check/"+strconv.Itoa(test.subID), nil)
//			assert.NoError(t, err)
//
//			rec := httptest.NewRecorder()
//
//			router := mux.NewRouter()
//			router.HandleFunc("/v1/subs/check/{id}", handler.CheckSub).Methods(http.MethodGet, http.MethodOptions)
//			router.ServeHTTP(rec, req)
//
//			assert.Equal(t, test.statusCode, rec.Code)
//
//			// Преобразование возвращаемого значения мока к ожидаемому типу
//			var expectedUpTo time.Time
//			switch v := test.expectUpTo.(type) {
//			case int64:
//				expectedUpTo = time.Unix(v, 0)
//			default:
//				t.Errorf("unexpected type %T for expectUpTo", v)
//			}
//
//			// Проверка, что ожидаемое значение совпадает с тем, что фактически возвращается
//			assert.Equal(t, expectedUpTo.Format(time.RFC3339), rec.Body.String())
//
//			mockUsecase.AssertExpectations(t)
//		})
//	}
//}

func TestSubscribe(t *testing.T) {
	tests := []struct {
		name       string
		subID      int
		expectErr  error
		statusCode int
	}{
		{
			name:       "Success",
			subID:      1,
			expectErr:  nil,
			statusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUsecase := new(mocks.SubsUsecase)
			handler := &SubsHandlerNotAu{
				SubsUsecase: mockUsecase,
			}

			mockUsecase.On("Subscribe", test.subID).Return(test.expectErr)

			req, err := http.NewRequest("POST", "/api/v1/subs/take_request", nil)
			assert.NoError(t, err)

			req.Form = make(map[string][]string)
			req.Form.Add("label", strconv.Itoa(test.subID))

			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/api/v1/subs/take_request", handler.TakeRequest).Methods(http.MethodPost, http.MethodOptions)
			router.ServeHTTP(rec, req)

			assert.Equal(t, test.statusCode, rec.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}
