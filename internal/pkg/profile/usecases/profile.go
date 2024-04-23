package usecases

import (
	"context"
	"mime/multipart"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

type ProfileInfo interface {
	CreateProfile(ctx context.Context, userID uint) *models.Profile
	GetProfileByUserID(ctx context.Context, userID uint) (*models.Profile, error)

	SetProfileCity(ctx context.Context, userID uint, data models.SetProfileCityNec)
	SetProfilePhone(ctx context.Context, userID uint, data models.SetProfilePhoneNec)
	//SetProfileRating(userID uint, data models.SetProfileRatingNec)
	SetProfile(userID uint, data models.SetProfileNec)
	SetProfileInfo(ctx context.Context, userID uint, file *multipart.FileHeader,
		data models.EditProfileNec) (*models.Profile, error)
	EditProfile(userID uint, data models.EditProfileNec)
	SetProfileApproved(userID uint)
}
