package postgresql_test

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"2023_2_Holi/collections/repository/collections_postgresql"
	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"
)

type MockDB struct {
	repository *mocks.MockFilmRepository
}

func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	// Вместо выполнения запроса к реальной базе данных,
	// делаем вызов к мок-репозиторию и возвращаем фиктивные данные.
	films, err := m.repository.GetFilmsByGenre(args[0].(string))
	if err != nil {
		return nil, err
	}

	// Создаем фиктивные строки результата для совместимости с интерфейсом *sql.DB.
	// В реальном приложении здесь можно использовать библиотеки для создания фейковых рядов.
	// Это простой пример, и он может потребовать доработки.
	rows := &sql.Rows{}
	for _, film := range films {
		rows.Rows = append(rows.Rows, []driver.Value{film.ID, film.Name, film.PreviewPath, film.Rating})
	}
	return rows, nil
}

func TestGetFilmsByGenre(t *testing.T) {
	tests := []struct {
		name             string
		genre            string
		expectedFilms    []domain.Film
		expectedError    error
		repositoryExpect func(repository *mocks.MockFilmRepository)
	}{
		{
			name:  "Success",
			genre: "Action",
			expectedFilms: []domain.Film{
				{ID: 1, Name: "Фильм 1", PreviewPath: "путь1", Rating: 4.5},
				{ID: 2, Name: "Фильм 2", PreviewPath: "путь2", Rating: 4.0},
			},
			repositoryExpect: func(repository *mocks.MockFilmRepository) {
				repository.On("GetFilmsByGenre", "Action").Return([]domain.Film{
					{ID: 1, Name: "Фильм 1", PreviewPath: "путь1", Rating: 4.5},
					{ID: 2, Name: "Фильм 2", PreviewPath: "путь2", Rating: 4.0},
				}, nil)
			},
		},
		{
			name:          "Empty result",
			genre:         "Comedy",
			expectedFilms: []domain.Film{},
			repositoryExpect: func(repository *mocks.MockFilmRepository) {
				repository.On("GetFilmsByGenre", "Comedy").Return([]domain.Film{}, nil)
			},
		},
		{
			name:          "Error from repository",
			genre:         "ErrorGenre",
			expectedError: errors.New("ошибка репозитория"),
			repositoryExpect: func(repository *mocks.MockFilmRepository) {
				repository.On("GetFilmsByGenre", "ErrorGenre").Return(nil, errors.New("ошибка репозитория"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mockRepository := new(mocks.MockFilmRepository)
			test.repositoryExpect(mockRepository)

			handler := collections_postgresql.NewFilmPostgresqlRepository(mockRepository)

			films, err := handler.GetFilmsByGenre(test.genre)

			assert.Equal(t, test.expectedFilms, films)
			assert.Equal(t, test.expectedError, err)

			mockRepository.AssertExpectations(t)
		})
	}
}
