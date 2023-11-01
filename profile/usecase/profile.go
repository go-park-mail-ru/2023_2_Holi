package profile_usecase

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"bytes"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type profileUseCase struct {
	profileRepo domain.ProfileRepository
}

func NewProfileUsecase(pr domain.ProfileRepository) domain.ProfileUsecase {
	return &profileUseCase{profileRepo: pr}
}

func (u *profileUseCase) GetUserData(userID int) (domain.User, error) {
	user, err := u.profileRepo.GetUser(userID)
	if err != nil {
		logs.LogError(logs.Logger, "profile_usecase", "GetUserData", err, err.Error())
		return domain.User{}, err
	}
	logs.Logger.Debug("Usecase GetUserData:", user)

	return user, nil
}

func (u *profileUseCase) UpdateUser(newUser domain.User) (domain.User, error) {
	updatedUser, err := u.profileRepo.UpdateUser(newUser)
	if err != nil {
		logs.LogError(logs.Logger, "profile_usecase", "UpdateUser", err, err.Error())
		return domain.User{}, err
	}
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
	sess, _ := session.NewSession()
	svc := s3.New(sess, aws.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))
	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(directory + "/" + strconv.Itoa(userID)),
		Body:        bytes.NewReader(imageData),
		ACL:         aws.String("public-read"),
		ContentType: aws.String("image/jpeg"),
	}

	if _, err := svc.PutObject(uploadInput); err != nil {
		logs.LogError(logs.Logger, "profile_usecase", "UploadImage", err, "Failed to upload image")
		return "", err
	}
	imagePath := "https://" + bucketName + ".hb." + defaultRegion + ".vkcs.cloud/" + directory + "/" + strconv.Itoa(userID)

	return imagePath, nil
}
