package profile_usecase

import "2023_2_Holi/domain"

type profileUseCase struct {
	profileRepo domain.ProfileRepository
}

func NewProfileUsecase(pr domain.ProfileRepository) domain.ProfileUsecase {
	return &profileUseCase{profileRepo: pr}
}

func (u *profileUseCase) GetProfile(userID int) (domain.User, error) {
	return domain.User{}, nil
}

func (u *profileUseCase) UpdateProfile(userID int, newUser domain.User) (domain.User, error) {
	return domain.User{}, nil
}

func (u *profileUseCase) UploadImage(userID int, image []byte) error {
	return nil
}
