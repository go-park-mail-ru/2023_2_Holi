package usecase

import (
	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSeriesByGenre(t *testing.T) {
	tests := []struct {
		name                      string
		genre                     int
		setSeriesRepoExpectations func(seriesRepo *mocks.SeriesRepository, series []domain.Video)
		good                      bool
	}{
		{
			name:  "GoodCase/Common",
			genre: 1,
			setSeriesRepoExpectations: func(seriesRepo *mocks.SeriesRepository, series []domain.Video) {
				seriesRepo.On("GetSeriesByGenre", 1).Return(series, nil)
			},
			good: true,
		},
		{
			name:  "ErrorCase/UsecaseError",
			genre: 2,
			setSeriesRepoExpectations: func(seriesRepo *mocks.SeriesRepository, series []domain.Video) {
				seriesRepo.On("GetSeriesByGenre", 2).Return(nil, errors.New("Some error"))
			},
			good: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sr := new(mocks.SeriesRepository)
			var series []domain.Video
			test.setSeriesRepoExpectations(sr, series)

			seriesCase := NewSeriesUsecase(sr)
			seriesCaseVideos, err := seriesCase.GetSeriesByGenre(test.genre)

			if test.good {
				assert.Nil(t, err)
				assert.EqualValues(t, seriesCaseVideos, series)
			} else {
				assert.NotNil(t, err)
				assert.Empty(t, seriesCaseVideos)
			}

			sr.AssertExpectations(t)
		})
	}
}

func TestGetSeriesData(t *testing.T) {
	tests := []struct {
		name                      string
		seriesID                  int
		setSeriesRepoExpectations func(seriesRepo *mocks.SeriesRepository, series *domain.Video, artists []domain.Cast, episodes []domain.Episode)
		good                      bool
	}{
		{
			name:     "GoodCase/Common",
			seriesID: 1,
			setSeriesRepoExpectations: func(seriesRepo *mocks.SeriesRepository, series *domain.Video, artists []domain.Cast, episodes []domain.Episode) {
				seriesRepo.On("GetSeriesData", mock.Anything).Return(*series, nil)
				seriesRepo.On("GetSeriesCast", mock.Anything).Return(artists, nil)
				seriesRepo.On("GetSeriesEpisodes", mock.Anything).Return(episodes, nil)
			},
			good: true,
		},
		{
			name:     "BadCase/GetSeriesDataError",
			seriesID: 2,
			setSeriesRepoExpectations: func(seriesRepo *mocks.SeriesRepository, series *domain.Video, artists []domain.Cast, episodes []domain.Episode) {
				seriesRepo.On("GetSeriesData", mock.Anything).Return(domain.Video{}, errors.New("some error"))
			},
			good: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sr := new(mocks.SeriesRepository)
			var series domain.Video
			var artists []domain.Cast
			var episodes []domain.Episode
			test.setSeriesRepoExpectations(sr, &series, artists, episodes)

			seriesCase := NewSeriesUsecase(sr)
			seriesCaseVideo, seriesCaseArtists, seriesCaseEpisodes, err := seriesCase.GetSeriesData(test.seriesID)

			if test.good {
				assert.Nil(t, err)
				assert.EqualValues(t, seriesCaseVideo, series)
				assert.EqualValues(t, seriesCaseArtists, artists)
				assert.EqualValues(t, seriesCaseEpisodes, episodes)
			} else {
				assert.NotNil(t, err)
				assert.Empty(t, seriesCaseVideo)
				assert.Empty(t, seriesCaseArtists)
				assert.Empty(t, seriesCaseEpisodes)
			}

			sr.AssertExpectations(t)
		})
	}
}

func TestGetCastPageSeries(t *testing.T) {
	tests := []struct {
		name                      string
		seriesID                  int
		setSeriesRepoExpectations func(seriesRepo *mocks.SeriesRepository, series []domain.Video, artist domain.Cast)
		good                      bool
	}{
		{
			name:     "GoodCase/Common",
			seriesID: 1,
			setSeriesRepoExpectations: func(seriesRepo *mocks.SeriesRepository, series []domain.Video, artist domain.Cast) {
				seriesRepo.On("GetCastPageSeries", mock.Anything).Return(series, nil)
				seriesRepo.On("GetCastNameSeries", mock.Anything).Return(artist, nil)
			},
			good: true,
		},
		{
			name:     "BadCase/GetCastPageSeriesError",
			seriesID: 2,
			setSeriesRepoExpectations: func(seriesRepo *mocks.SeriesRepository, series []domain.Video, artist domain.Cast) {
				seriesRepo.On("GetCastPageSeries", mock.Anything).Return([]domain.Video{}, errors.New("some error"))
			},
			good: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sr := new(mocks.SeriesRepository)
			var series []domain.Video
			var artist domain.Cast
			test.setSeriesRepoExpectations(sr, series, artist)

			seriesCase := NewSeriesUsecase(sr)
			seriesCaseVideos, seriesCaseArtist, err := seriesCase.GetCastPageSeries(test.seriesID)

			if test.good {
				assert.Nil(t, err)
				assert.EqualValues(t, seriesCaseVideos, series)
				assert.EqualValues(t, seriesCaseArtist, artist)
			} else {
				assert.NotNil(t, err)
				assert.Empty(t, seriesCaseVideos)
				assert.Empty(t, seriesCaseArtist)
			}

			sr.AssertExpectations(t)
		})
	}
}
