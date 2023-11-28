package usecase

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"bytes"
	"crypto/rand"
	"strconv"

	auth_usecase "2023_2_Holi/auth/usecase"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type profileUseCase struct {
	profileRepo domain.ProfileRepository
	svc         s3iface.S3API
}

func NewProfileUsecase(pr domain.ProfileRepository, svc s3iface.S3API) domain.ProfileUsecase {
	return &profileUseCase{profileRepo: pr, svc: svc}
}

func (u *profileUseCase) GetUserData(userID int) (domain.User, error) {
	user, err := u.profileRepo.GetUser(userID)
	if err != nil {
		logs.LogError(logs.Logger, "profile_usecase", "GetUserData", err, err.Error())
		return domain.User{}, err
	}
	user.Password = []byte{}
	logs.Logger.Debug("Usecase GetUserData:", user)

	return user, nil
}

func (u *profileUseCase) UpdateUser(newUser domain.User) (domain.User, error) {
	if len(newUser.Password) != 0 {
		salt := make([]byte, 8)
		rand.Read(salt)
		newUser.Password = auth_usecase.HashPassword(salt, newUser.Password)
	}

	oldUser, err := u.profileRepo.GetUser(newUser.ID)
	if err != nil {
		logs.LogError(logs.Logger, "profile_usecase", "UpdateUser", err, err.Error())
		return domain.User{}, err
	}

	newUser = getOldFields(newUser, oldUser)
	updatedUser, err := u.profileRepo.UpdateUser(newUser)
	if err != nil {
		logs.LogError(logs.Logger, "profile_usecase", "UpdateUser", err, err.Error())
		return domain.User{}, err
	}
	updatedUser.Password = []byte{}
	logs.Logger.Debug("Usecase UpdateUser:", updatedUser)

	return updatedUser, nil
}

const (
	vkCloudHotboxEndpoint = "https://hb.vkcs.cloud"
	defaultRegion         = "ru-msk"
	bucketName            = "static_holi"
	directory             = "User_Images"
)

func (u *profileUseCase) UploadImage(userID int, imageData []byte) (string, error) {
	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Body:        bytes.NewReader(imageData),
		ACL:         aws.String("public-read"),
		ContentType: aws.String("image/jpeg"),
		Key:         aws.String(directory + "/" + strconv.Itoa(userID)),
	}

	if _, err := u.svc.PutObject(uploadInput); err != nil {
		logs.LogError(logs.Logger, "usecase", "UploadImage", err, "Failed to upload image")
		return "", err
	}
	imagePath := "https://" + bucketName + ".hb." + defaultRegion + ".vkcs.cloud/" + directory + "/" + strconv.Itoa(userID)

	return imagePath, nil
}

func getOldFields(newUser domain.User, oldUser domain.User) domain.User {
	if newUser.Name == "" {
		newUser.Name = oldUser.Name
	}
	if newUser.Email == "" {
		newUser.Email = oldUser.Email
	}
	if newUser.ImagePath == "" {
		newUser.ImagePath = oldUser.ImagePath
	}
	if string(newUser.Password) == "" {
		newUser.Password = oldUser.Password
	}
	return newUser
}
