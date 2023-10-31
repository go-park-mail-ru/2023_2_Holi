package auth_redis

import (
	"2023_2_Holi/domain"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		session domain.Session
		good    bool
		err     error
	}{
		{
			name: "GoodCase/Common",
			session: domain.Session{
				Token:     "123",
				ExpiresAt: time.Now().Add(24 * time.Hour),
				UserID:    1,
			},
			good: true,
			err:  nil,
		},
		{
			name: "GoodCase/SameToken",
			session: domain.Session{
				Token:     "123",
				ExpiresAt: time.Now().Add(24 * time.Hour),
				UserID:    1,
			},
			good: true,
		},
		{
			name:    "BadCase/EmptyToken",
			session: domain.Session{},
			err:     errors.New("empty token"),
		},
	}
	db, mock := redismock.NewClientMock()
	r := NewSessionRedisRepository(db)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.good {
				mock.ExpectSet(test.session.Token, test.session.UserID, test.session.ExpiresAt.Sub(time.Now())).SetVal("")
			}
			err := r.Add(test.session)

			if test.good {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			} else {
				if err == nil || err.Error() != test.err.Error() {
					t.Errorf("Expected error: %v, got: %v", test.err, err)
				}
			}

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestDeleteByToken(t *testing.T) {
	tests := []struct {
		name  string
		token string
		good  bool
		err   error
	}{
		{
			name:  "GoodCase/Common",
			token: "12312dcdscsad",
			good:  true,
			err:   nil,
		},
		{
			name:  "BadCase/EmptyToken",
			token: "",
			err:   errors.New("empty token"),
		},
	}
	db, mock := redismock.NewClientMock()
	r := NewSessionRedisRepository(db)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.good {
				mock.ExpectDel(test.token).SetVal(1)
			}
			err := r.DeleteByToken(test.token)

			if test.good {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			} else {
				if err == nil || err.Error() != test.err.Error() {
					t.Errorf("Expected error: %v, got: %v", test.err, err)
				}
			}

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestSessionExists(t *testing.T) {
	tests := []struct {
		name  string
		token string
		good  bool
		err   error
	}{
		{
			name:  "GoodCase/Common",
			token: "12312dcdscsad",
			good:  true,
			err:   nil,
		},
		{
			name:  "BadCase/EmptyToken",
			token: "",
			err:   errors.New("empty token"),
		},
	}
	db, mock := redismock.NewClientMock()
	r := NewSessionRedisRepository(db)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.good {
				mock.ExpectExists(test.token).SetVal(1)
			}
			_, err := r.SessionExists(test.token)

			if test.good {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			} else {
				if err == nil || err.Error() != test.err.Error() {
					t.Errorf("Expected error: %v, got: %v", test.err, err)
				}
			}

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}
