package postgres

//
//import (
//	"frontend/domain"
//	"context"
//	"errors"
//	"github.com/jackc/pgx/v5/pgconn"
//	"github.com/pashagolub/pgxmock/v3"
//	"github.com/stretchr/testify/require"
//	"testing"
//)
//
//const userID = 1
//
//const testAddToFavouritesQuery = `
//	INSERT INTO favourite \(video_id, user_id\)
//	VALUES \(\$1, \$2\)
//`
//
//const testDeleteFromFavouritesQuery = `
//	DELETE FROM favourite
//	WHERE video_id = \$1 AND user_id = \$2
//`
//
//const testSelectAllQuery = `
//	SELECT v.id, v.name, v.description,
//		v.preview_path, v.preview_video_path, v.release_year, v.rating, v.age_restriction
//	FROM video AS v
//		JOIN favourite AS f ON video_id = v.id
//	WHERE f.user_id = \$1
//`
//
//func TestInsertIntoFavourites(t *testing.T) {
//	tests := []struct {
//		name    string
//		videoID int
//		err     error
//		good    bool
//	}{
//		{
//			name:    "GoodCase/Common",
//			videoID: 1,
//			good:    true,
//		},
//		{
//			name:    "BadCase/NegativeVideoID",
//			err:     &pgconn.PgError{Code: "23503"},
//			videoID: -1,
//		},
//		{
//			name:    "BadCase/OutOfRangeVideoID",
//			err:     &pgconn.PgError{Code: "23503"},
//			videoID: 123456789,
//		},
//		{
//			name:    "BadCase/ZeroVideoID",
//			err:     &pgconn.PgError{Code: "23503"},
//			videoID: 0,
//		},
//		{
//			name:    "BadCase/AlreadyExistedVideo",
//			err:     &pgconn.PgError{Code: "23505"},
//			videoID: 4,
//		},
//	}
//
//	mockDB, err := pgxmock.NewPool()
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer mockDB.Close()
//
//	r := NewFavouritesPostgresqlRepository(mockDB, context.Background())
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			eq := mockDB.ExpectExec(testAddToFavouritesQuery).
//				WithArgs(test.videoID, userID)
//			if test.good {
//				eq.WillReturnResult(pgxmock.NewResult("", 1))
//			} else {
//				eq.WillReturnError(test.err)
//			}
//
//			err = r.InsertIntoFavourites(test.videoID, userID)
//			if test.good {
//				require.Nil(t, err)
//			} else {
//				require.NotNil(t, err)
//			}
//
//			err = mockDB.ExpectationsWereMet()
//			require.Nil(t, err)
//		})
//	}
//}
//
//func TestDeleteFromFavourites(t *testing.T) {
//	tests := []struct {
//		name    string
//		videoID int
//		result  pgconn.CommandTag
//		good    bool
//	}{
//		{
//			name:    "GoodCase/Common",
//			videoID: 1,
//			result:  pgxmock.NewResult("", 1),
//			good:    true,
//		},
//		{
//			name:    "BadCase/NegativeVideoID",
//			result:  pgxmock.NewResult("", 0),
//			videoID: -1,
//		},
//		{
//			name:    "BadCase/OutOfRangeVideoID",
//			result:  pgxmock.NewResult("", 0),
//			videoID: 123456789,
//		},
//		{
//			name:    "BadCase/ZeroVideoID",
//			result:  pgxmock.NewResult("", 0),
//			videoID: 0,
//		},
//	}
//
//	mockDB, err := pgxmock.NewPool()
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer mockDB.Close()
//
//	r := NewFavouritesPostgresqlRepository(mockDB, context.Background())
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			mockDB.ExpectExec(testDeleteFromFavouritesQuery).
//				WithArgs(test.videoID, userID).
//				WillReturnResult(test.result)
//
//			err = r.DeleteFromFavourites(test.videoID, userID)
//			if test.good {
//				require.Nil(t, err)
//			} else {
//				require.NotNil(t, err)
//			}
//
//			err = mockDB.ExpectationsWereMet()
//			require.Nil(t, err)
//		})
//	}
//}
//
//func TestSelectAllFavourites(t *testing.T) {
//	tests := []struct {
//		name   string
//		videos []domain.Video
//		good   bool
//	}{
//		{
//			name: "GoodCase/Common",
//			videos: []domain.Video{
//				domain.Video{
//					ID:               1,
//					Name:             "some",
//					Description:      "desc",
//					PreviewPath:      "path",
//					PreviewVideoPath: "video_path",
//					ReleaseYear:      2007,
//					Rating:           9.5,
//					AgeRestriction:   13,
//				},
//				domain.Video{
//					ID:               2,
//					Name:             "some",
//					Description:      "desc",
//					PreviewPath:      "path",
//					PreviewVideoPath: "video_path",
//					ReleaseYear:      2007,
//					Rating:           9.5,
//					AgeRestriction:   13,
//				},
//			},
//			good: true,
//		},
//		{
//			name:   "GoodCase/EmptyFavourites",
//			videos: make([]domain.Video, 2, 2),
//			good:   true,
//		},
//	}
//
//	mockDB, err := pgxmock.NewPool()
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer mockDB.Close()
//
//	r := NewFavouritesPostgresqlRepository(mockDB, context.Background())
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			rows := mockDB.NewRows([]string{"id", "name", "description", "preview_path",
//				"preview_video_path", "release_year", "rating", "age_restriction"}).
//				AddRow(test.videos[0].ID, test.videos[0].Name, test.videos[0].Description,
//					test.videos[0].PreviewPath, test.videos[0].PreviewVideoPath, test.videos[0].ReleaseYear,
//					test.videos[0].Rating, test.videos[0].AgeRestriction).
//				AddRow(test.videos[1].ID, test.videos[1].Name, test.videos[1].Description,
//					test.videos[1].PreviewPath, test.videos[1].PreviewVideoPath, test.videos[1].ReleaseYear,
//					test.videos[1].Rating, test.videos[1].AgeRestriction)
//
//			eq := mockDB.ExpectQuery(testSelectAllQuery).
//				WithArgs(userID)
//			if test.good {
//				eq.WillReturnRows(rows)
//			} else {
//				eq.WillReturnError(errors.New("some"))
//			}
//
//			videos, err := r.SelectAllFavourites(userID)
//			if test.good {
//				require.Equal(t, test.videos, videos)
//				require.Nil(t, err)
//			} else {
//				require.NotNil(t, err)
//			}
//
//			err = mockDB.ExpectationsWereMet()
//			require.Nil(t, err)
//		})
//	}
//}
