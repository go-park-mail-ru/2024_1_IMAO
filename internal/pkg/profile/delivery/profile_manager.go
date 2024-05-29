package delivery

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	protobuf "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery/protobuf"
	profileusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/usecases"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProfileManager struct {
	protobuf.UnimplementedProfileServer

	profileStorage profileusecases.ProfileStorageInterface
}

func NewProfileManager(storage profileusecases.ProfileStorageInterface) *ProfileManager {
	return &ProfileManager{
		profileStorage: storage,
	}
}

func CleanProfileData(profile *protobuf.ProfileData) *models.Profile {
	return &models.Profile{
		ID:      uint(profile.ID),
		UserID:  uint(profile.UserID),
		Name:    profile.Name,
		Surname: profile.Surname,
		City: models.City{
			ID:          uint(profile.CityID),
			CityName:    profile.CityName,
			Translation: profile.Translation,
		},
		Phone:           profile.Phone,
		Avatar:          profile.Avatar,
		RegisterTime:    profile.RegisterTime.AsTime(),
		Rating:          profile.Rating,
		ReactionsCount:  profile.ReactionsCount,
		Approved:        profile.Approved,
		MerchantsName:   profile.MerchantsName,
		SubersCount:     int(profile.SubersCount),
		SubonsCount:     int(profile.SubonsCount),
		AvatarIMG:       profile.AvatarIMG,
		ActiveAddsCount: int(profile.ActiveAddsCount),
		SoldAddsCount:   int(profile.SoldAddsCount),
		CartNum:         uint(profile.CartNum),
		FavNum:          uint(profile.FavNum),
	}
}

func newProtobufProfile(profile *models.Profile) *protobuf.ProfileData {
	return &protobuf.ProfileData{
		ID:              uint64(profile.ID),
		UserID:          uint64(profile.UserID),
		Name:            profile.Name,
		Surname:         profile.Surname,
		CityID:          uint64(profile.City.ID),
		CityName:        profile.City.CityName,
		Translation:     profile.City.Translation,
		Phone:           profile.Phone,
		Avatar:          profile.Avatar,
		RegisterTime:    timestamppb.New(profile.RegisterTime),
		Rating:          profile.Rating,
		ReactionsCount:  profile.ReactionsCount,
		Approved:        profile.Approved,
		MerchantsName:   profile.MerchantsName,
		SubersCount:     int64(profile.SubersCount),
		SubonsCount:     int64(profile.SubonsCount),
		AvatarIMG:       profile.AvatarIMG,
		ActiveAddsCount: int64(profile.ActiveAddsCount),
		SoldAddsCount:   int64(profile.SoldAddsCount),
		CartNum:         int64(profile.CartNum),
		FavNum:          int64(profile.FavNum),
	}
}

func (manager *ProfileManager) GetProfile(ctx context.Context,
	in *protobuf.ProfileIDRequest) (*protobuf.ProfileData, error) {
	id := in.GetID()
	storage := manager.profileStorage

	profile, err := storage.GetProfileByUserID(ctx, uint(id))
	if err != nil {
		return nil, err
	}

	return newProtobufProfile(profile), nil
}

func (manager *ProfileManager) CreateProfile(ctx context.Context,
	in *protobuf.ProfileIDRequest) (*protobuf.ProfileData, error) {
	id := in.GetID()
	storage := manager.profileStorage

	profile := storage.CreateProfile(ctx, uint(id))

	return newProtobufProfile(profile), nil
}

func (manager *ProfileManager) SetProfileCity(ctx context.Context,
	in *protobuf.SetCityRequest) (*protobuf.ProfileData, error) {
	id := in.GetID()
	storage := manager.profileStorage

	profile, err := storage.SetProfileCity(ctx, uint(id), models.City{
		ID:          uint(in.GetCityID()),
		CityName:    in.GetCityName(),
		Translation: in.GetTranslation(),
	})
	if err != nil {
		return nil, err
	}

	return newProtobufProfile(profile), nil
}

func (manager *ProfileManager) SetProfilePhone(ctx context.Context,
	in *protobuf.SetPhoneRequest) (*protobuf.ProfileData, error) {
	id := in.GetID()
	storage := manager.profileStorage

	profile, err := storage.SetProfilePhone(ctx, uint(id), models.SetProfilePhoneNec{Phone: in.GetPhone()})
	if err != nil {
		return nil, err
	}

	return newProtobufProfile(profile), nil
}

func (manager *ProfileManager) EditProfile(ctx context.Context,
	in *protobuf.EditProfileRequest) (*protobuf.ProfileData, error) {
	storage := manager.profileStorage

	profile, err := storage.SetProfileInfo(ctx, uint(in.GetID()), models.EditProfileNec{
		Name:    in.GetName(),
		Surname: in.GetSurname(),
		Avatar:  in.GetAvatar(),
	})
	if err != nil {
		return nil, err
	}

	return newProtobufProfile(profile), nil
}

func (manager *ProfileManager) AppendSubByIDs(ctx context.Context,
	in *protobuf.UserIdMerchantIdRequest) (*protobuf.AppendSubResponse, error) {
	userID := in.GetUserId()
	merchantID := in.GetMerchantId()
	storage := manager.profileStorage

	isAppended := storage.AppendSubByIDs(ctx, uint(userID), uint(merchantID))

	return &protobuf.AppendSubResponse{IsAppended: isAppended}, nil
}
