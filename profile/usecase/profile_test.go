package usecase

import (
	"2023_2_Holi/domain"
	"2023_2_Holi/domain/mocks"
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserData(t *testing.T) {
	tests := []struct {
		name                       string
		id                         int
		setProfileRepoExpectations func(userID int, prRepo *mocks.ProfileRepository, user *domain.User)
		good                       bool
	}{
		{
			name: "GoodCase/Common",
			id:   1,
			setProfileRepoExpectations: func(userID int, prRepo *mocks.ProfileRepository, user *domain.User) {
				faker.FakeData(user)
				user.ID = userID
				prRepo.On("GetUser", mock.Anything).Return(*user, nil)
			},
			good: true,
		},
		{
			name: "BadCase/UserNotFound",
			id:   1,
			setProfileRepoExpectations: func(userID int, prRepo *mocks.ProfileRepository, user *domain.User) {
				faker.FakeData(user)
				user.ID = userID
				prRepo.On("GetUser", mock.Anything).Return(*user, errors.New("user not found"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			pr := new(mocks.ProfileRepository)
			var user domain.User
			test.setProfileRepoExpectations(test.id, pr, &user)

			profileUcase := NewProfileUsecase(pr, &mocks.MockS3Client{})
			userData, err := profileUcase.GetUserData(test.id)

			if test.good {
				assert.Nil(t, err)
				assert.NotEmpty(t, user)
				assert.Equal(t, userData.ID, user.ID)
			} else {
				assert.NotNil(t, err)
				assert.Empty(t, userData)
			}

			pr.AssertExpectations(t)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name                       string
		oldUser                    domain.User
		newUser                    domain.User
		setProfileRepoExpectations func(user *domain.User, prRepo *mocks.ProfileRepository, newUser *domain.User, oldUser *domain.User)
		good                       bool
	}{
		{
			name: "GoodCase/Common",
			oldUser: domain.User{
				ID:    1,
				Name:  "Alex",
				Email: "alex@mail.ru",
			},
			newUser: domain.User{
				ID:    1,
				Name:  "Max",
				Email: "max@mail.ru",
			},
			setProfileRepoExpectations: func(user *domain.User, prRepo *mocks.ProfileRepository, newUser *domain.User, oldUser *domain.User) {
				prRepo.On("GetUser", mock.Anything).Return(*oldUser, nil)
				prRepo.On("UpdateUser", mock.Anything).Return(*newUser, nil)
			},
			good: true,
		},
		{
			name: "BadCase/UserNotFound",
			oldUser: domain.User{
				ID:    1,
				Name:  "Alex",
				Email: "alex@mail.ru",
			},
			newUser: domain.User{
				ID:    2,
				Name:  "Max",
				Email: "max@mail.ru",
			},
			setProfileRepoExpectations: func(user *domain.User, prRepo *mocks.ProfileRepository, newUser *domain.User, oldUser *domain.User) {
				prRepo.On("GetUser", mock.Anything).Return(*oldUser, errors.New("user not found"))
			},
		},
		{
			name: "BadCase/FailedToUpdate",
			oldUser: domain.User{
				ID:    1,
				Name:  "Alex",
				Email: "alex@mail.ru",
			},
			newUser: domain.User{
				ID:    2,
				Name:  "Max",
				Email: "max@mail.ru",
			},
			setProfileRepoExpectations: func(user *domain.User, prRepo *mocks.ProfileRepository, newUser *domain.User, oldUser *domain.User) {
				prRepo.On("GetUser", mock.Anything).Return(*oldUser, nil)
				prRepo.On("UpdateUser", mock.Anything).Return(*newUser, errors.New("failed to update"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			pr := new(mocks.ProfileRepository)
			newUser := test.newUser
			oldUser := test.oldUser
			test.setProfileRepoExpectations(&test.newUser, pr, &newUser, &oldUser)

			mockSvc := &mocks.MockS3Client{}
			profileUcase := NewProfileUsecase(pr, mockSvc)
			updatedUser, err := profileUcase.UpdateUser(test.newUser)

			if test.good {
				assert.Nil(t, err)
				assert.NotEmpty(t, newUser)
				assert.Equal(t, newUser, updatedUser)
			} else {
				assert.NotNil(t, err)
				assert.Empty(t, updatedUser)
			}

			pr.AssertExpectations(t)
		})
	}
}

func TestUploadImage(t *testing.T) {

	pr := new(mocks.ProfileRepository)
	mockSvc := &mocks.MockS3Client{}
	usecase := NewProfileUsecase(pr, mockSvc)

	userID := 123
	imageData := getImageData()

	expectedUploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Body:        bytes.NewReader(imageData),
		ACL:         aws.String("public-read"),
		ContentType: aws.String("image/jpeg"),
		Key:         aws.String(directory + "/" + strconv.Itoa(userID)),
	}

	_, err := mockSvc.PutObject(expectedUploadInput)
	expectedImagePath := "https://" + bucketName + ".hb." + defaultRegion + ".vkcs.cloud/" + directory + "/" + strconv.Itoa(userID)

	imagePath, err := usecase.UploadImage(userID, imageData)
	assert.NoError(t, err)
	assert.Equal(t, expectedImagePath, imagePath)

	assert.Equal(t, 2, mockSvc.PutObjectCallCount)
	assert.Equal(t, expectedUploadInput, mockSvc.PutObjectInput)
}

func getImageData() []byte {
	file, err := os.Open("ava.jpg")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return nil
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println("Ошибка при декодировании изображения:", err)
		return nil
	}
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		return nil
	}
	return buf.Bytes()

}
