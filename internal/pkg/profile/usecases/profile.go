package usecases

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

type ProfileInfo interface {
	CreateProfile(userID uint) *models.Profile
	GetProfileByUserID(userID uint) (*models.Profile, error)

	SetProfileCity(userID uint, data models.SetProfileCityNec)
	SetProfilePhone(userID uint, data models.SetProfilePhoneNec)
	SetProfileRating(userID uint, data models.SetProfileRatingNec)
	SetProfile(userID uint, data models.SetProfileNec)
	EditProfile(userID uint, data models.EditProfileNec)
	SetProfileApproved(userID uint)
}
