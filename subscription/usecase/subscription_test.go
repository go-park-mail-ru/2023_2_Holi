package usecase

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"2023_2_Holi/domain/mocks"
)

func TestSubscribe(t *testing.T) {
	tests := []struct {
		name                string
		setRepoExpectations func(repo *mocks.SubsRepository)
		subID               int
		expectedError       error
	}{
		{
			name: "Subscribe Success",
			setRepoExpectations: func(repo *mocks.SubsRepository) {
				repo.On("Subscribe", 1).Return(nil)
			},
			subID:         1,
			expectedError: nil,
		},
		{
			name: "Subscribe Error",
			setRepoExpectations: func(repo *mocks.SubsRepository) {
				repo.On("Subscribe", 1).Return(errors.New("subscription error"))
			},
			subID:         1,
			expectedError: errors.New("subscription error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := new(mocks.SubsRepository)
			test.setRepoExpectations(repo)

			usecase := NewSubsUsecase(repo)

			err := usecase.Subscribe(test.subID)

			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestUnSubscribe(t *testing.T) {
	tests := []struct {
		name                string
		setRepoExpectations func(repo *mocks.SubsRepository)
		subID               int
		expectedError       error
	}{
		{
			name: "UnSubscribe Success",
			setRepoExpectations: func(repo *mocks.SubsRepository) {
				repo.On("UnSubscribe", 1).Return(nil)
			},
			subID:         1,
			expectedError: nil,
		},
		{
			name: "UnSubscribe Error",
			setRepoExpectations: func(repo *mocks.SubsRepository) {
				repo.On("UnSubscribe", 1).Return(errors.New("unsubscription error"))
			},
			subID:         1,
			expectedError: errors.New("unsubscription error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := new(mocks.SubsRepository)
			test.setRepoExpectations(repo)

			usecase := NewSubsUsecase(repo)

			err := usecase.UnSubscribe(test.subID)

			assert.Equal(t, test.expectedError, err)
		})
	}
}

//func TestCheckSub(t *testing.T) {
//	tests := []struct {
//		name                string
//		setRepoExpectations func(repo *mocks.SubsRepository)
//		subID               int
//		expectedSub         time.Time
//		expectedStatus      bool
//		expectedError       error
//	}{
//		{
//			name: "CheckSub Active Subscription",
//			setRepoExpectations: func(repo *mocks.SubsRepository) {
//				subUpTo := time.Now().Add(24 * time.Hour)
//				repo.On("CheckSub", 1).Return(subUpTo, nil)
//			},
//			subID:          1,
//			expectedSub:    time.Now().Add(24 * time.Hour),
//			expectedStatus: true,
//			expectedError:  nil,
//		},
//		{
//			name: "CheckSub Expired Subscription",
//			setRepoExpectations: func(repo *mocks.SubsRepository) {
//				subUpTo := time.Now().Add(-24 * time.Hour)
//				repo.On("CheckSub", 1).Return(subUpTo, nil)
//			},
//			subID:          1,
//			expectedSub:    time.Now(),
//			expectedStatus: false,
//			expectedError:  nil,
//		},
//		{
//			name: "CheckSub Error",
//			setRepoExpectations: func(repo *mocks.SubsRepository) {
//				repo.On("CheckSub", 1).Return(time.Now(), errors.New("subscription check error"))
//			},
//			subID:          1,
//			expectedSub:    time.Now(),
//			expectedStatus: false,
//			expectedError:  errors.New("subscription check error"),
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			repo := new(mocks.SubsRepository)
//			test.setRepoExpectations(repo)
//
//			usecase := NewSubsUsecase(repo)
//
//			sub, status, err := usecase.CheckSub(test.subID)
//
//			assert.Equal(t, test.expectedError, err)
//
//			assertion.Equal(t, test.expectedSub, sub)
//			assert.Equal(t, test.expectedStatus, status)
//		})
//	}
//}
