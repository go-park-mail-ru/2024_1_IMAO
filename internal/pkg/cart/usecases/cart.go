package usecases

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	advuc "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/usecases"
	useruc "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"
)

type CartInfo interface {
	GetCartByUserID(userID uint, userList useruc.UsersInfo, advertsList advuc.AdvertsInfo) ([]*models.Advert, error)
	DeleteAdvByIDs(userID uint, advertID uint, userList useruc.UsersInfo, advertsList advuc.AdvertsInfo) error
	AppendAdvByIDs(userID uint, advertID uint, userList useruc.UsersInfo, advertsList advuc.AdvertsInfo) bool
}
