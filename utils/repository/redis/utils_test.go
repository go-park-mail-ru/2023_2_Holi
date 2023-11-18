package redis

import (
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestSessionExists(t *testing.T) {
	tests := []struct {
		name  string
		token string
		id    string
		good  bool
	}{
		{
			name:  "GoodCase/Common",
			id:    "1",
			token: "fo4380cnu3inciou4",
			good:  true,
		},
		{
			name:  "BadCase/EmptyToken",
			token: "",
		},
		{
			name:  "BadCase/InappropriateToken",
			token: "123",
		},
	}
	db, rm := redismock.NewClientMock()
	r := NewUtilsRedisRepository(db)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.good {
				rm.ExpectGet(test.token).SetVal(test.id)
			}
			id, err := r.GetIdFromStorage(test.token)

			if test.good {
				expectedID, _ := strconv.Atoi(test.id)
				assert.Equal(t, expectedID, id)
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}

			err = rm.ExpectationsWereMet()
			assert.Nil(t, err)
		})
	}
}
