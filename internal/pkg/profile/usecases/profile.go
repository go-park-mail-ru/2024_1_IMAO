package usecases

import (
	"context"
	"mime/multipart"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

type ProfileStorageInterface interface {
	CreateProfile(ctx context.Context, userID uint) *models.Profile
	GetProfileByUserID(ctx context.Context, userID uint) (*models.Profile, error)

	SetProfileCity(ctx context.Context, userID uint, data models.City) (*models.Profile, error)
	SetProfilePhone(ctx context.Context, userID uint, data models.SetProfilePhoneNec) (*models.Profile, error)
	SetProfileInfo(ctx context.Context, userID uint, file *multipart.FileHeader,
		data models.EditProfileNec) (*models.Profile, error)
	//SetProfileApproved(userID uint)
}
