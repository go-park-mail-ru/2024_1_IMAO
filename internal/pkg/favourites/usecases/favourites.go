package usecases

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	advuc "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/usecases"
	useruc "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"
)

type FavouritesStorageInterface interface {
	GetFavouritesByUserID(ctx context.Context, userID uint) ([]*models.ReturningAdInList, error)
	DeleteAdvByIDs(ctx context.Context, userID uint, advertID uint, userList useruc.UsersStorageInterface, advertsList advuc.AdvertsStorageInterface) error
	AppendAdvByIDs(ctx context.Context, userID uint, advertID uint, userList useruc.UsersStorageInterface, advertsList advuc.AdvertsStorageInterface) bool
}
