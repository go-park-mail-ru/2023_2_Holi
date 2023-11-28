package usecase

import (
	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSearchData(t *testing.T) {
	tests := []struct {
		name                      string
		searchStr                 string
		films                     []domain.Video
		cast                      []domain.Cast
		setSearchRepoExpectations func(searchRepo *mocks.SearchRepository, films []domain.Video, cast []domain.Cast)
		good                      bool
	}{
		{
			name:      "GoodCase/Common",
			searchStr: "Leonardo",
			films: []domain.Video{
				{
					ID:          1,
					Name:        "Leonardo",
					PreviewPath: "path",
				},
			},
			cast: []domain.Cast{
				{
					ID:   1,
					Name: "Leonardo Di Caprio",
				},
				{
					ID:   2,
					Name: "Leonardo Turtle",
				},
			},
			setSearchRepoExpectations: func(searchRepo *mocks.SearchRepository, films []domain.Video, cast []domain.Cast) {
				searchRepo.On("GetSuitableFilms", mock.Anything).Return(films, nil)
				searchRepo.On("GetSuitableCast", mock.Anything).Return(cast, nil)
			},
			good: true,
		},
		{
			name:      "BadCase/NoData",
			searchStr: "Le",
			setSearchRepoExpectations: func(searchRepo *mocks.SearchRepository, films []domain.Video, cast []domain.Cast) {
				searchRepo.On("GetSuitableFilms", mock.Anything).Return(films, domain.ErrNotFound)
				searchRepo.On("GetSuitableCast", mock.Anything).Return(cast, domain.ErrNotFound)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			sr := new(mocks.SearchRepository)
			test.setSearchRepoExpectations(sr, test.films, test.cast)

			searchUcase := NewSearchUsecase(sr)
			searchData, err := searchUcase.GetSearchData(test.searchStr)

			if test.good {
				assert.Nil(t, err)
				assert.NotEmpty(t, test.films)
				assert.NotEmpty(t, test.cast)
				assert.Equal(t, searchData.Films, test.films)
				assert.Equal(t, searchData.Cast, test.cast)
			} else {
				assert.NotNil(t, err)
				assert.Empty(t, searchData)
			}

			sr.AssertExpectations(t)
		})
	}
}
