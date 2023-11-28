package http

import (
	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSearchData(t *testing.T) {
	tests := []struct {
		name                 string
		searchStr            string
		setUCaseExpectations func(usecase *mocks.SearchUsecase, searchData *domain.SearchData)
		status               int
	}{
		{
			name:      "GoodCase/Common",
			searchStr: "Leonardo",
			setUCaseExpectations: func(usecase *mocks.SearchUsecase, searchData *domain.SearchData) {
				faker.FakeData(searchData)
				usecase.On("GetSearchData", mock.Anything).Return(*searchData, nil)
			},
			status: http.StatusOK,
		},
		{
			name:      "BadCase/Common",
			searchStr: "Laus",
			setUCaseExpectations: func(usecase *mocks.SearchUsecase, searchData *domain.SearchData) {
				usecase.On("GetSearchData", mock.Anything).Return(*searchData, domain.ErrNotFound)
			},
			status: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			router := mux.NewRouter()
			mockUsecase := new(mocks.SearchUsecase)
			var data domain.SearchData
			test.setUCaseExpectations(mockUsecase, &data)

			req, err := http.NewRequest("GET", "/api/v1/search/"+test.searchStr, nil)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()

			NewSearchHandler(router, mockUsecase)

			handler := &SearchHandler{
				SearchUsecase: mockUsecase,
			}

			router.HandleFunc("/api/v1/search/{searchStr}", handler.GetSearchData).Methods("GET")
			router.ServeHTTP(rec, req)

			assert.Equal(t, test.status, rec.Code)
		})
	}
}
