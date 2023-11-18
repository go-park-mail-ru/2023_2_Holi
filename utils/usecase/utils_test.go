package usecase

import (
	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGetIdBy(t *testing.T) {
	tests := []struct {
		name            string
		setExpectations func(ur *mocks.UtilsRepository, id int)
		token           string
		id              int
		good            bool
	}{
		{
			name: "GoodCase/Common",
			setExpectations: func(ur *mocks.UtilsRepository, id int) {
				ur.On("GetIdFromStorage", mock.Anything).Return(id, nil)
			},
			id:    1,
			token: "fo4380cnu3inciou4",
			good:  true,
		},
		{
			name: "BadCase/EmptyToken",
			setExpectations: func(ur *mocks.UtilsRepository, id int) {
				ur.On("GetIdFromStorage", mock.Anything).Return(id, domain.ErrInvalidToken)
			},
			token: "",
		},
		{
			name: "BadCase/InappropriateToken",
			setExpectations: func(ur *mocks.UtilsRepository, id int) {
				ur.On("GetIdFromStorage", mock.Anything).Return(id, errors.New("some"))
			},
			id:    0,
			token: "123",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			ur := new(mocks.UtilsRepository)
			test.setExpectations(ur, test.id)

			uu := NewUtilsUsecase(ur)
			id, err := uu.GetIdBy(test.token)

			if test.good {
				assert.Equal(t, test.id, id)
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}

			ur.AssertExpectations(t)
		})
	}
}
