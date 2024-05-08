package delivery

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	protobuf "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/delivery/protobuf"
	cartusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/usecases"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CartManager struct {
	protobuf.UnimplementedCartServer

	cartStorage cartusecases.CartStorageInterface
}

func NewCartManager(storage cartusecases.CartStorageInterface) *CartManager {
	return &CartManager{
		cartStorage: storage,
	}
}

func ReturningAdvertItem(returningAdvertList *protobuf.ReturningAdvertList) []*models.ReturningAdvert {

	var returningAdvertsList []*models.ReturningAdvert

	for i := 0; i < len(returningAdvertList.Adverts); i++ {
		returningAdvert := models.ReturningAdvert{
			Advert: models.Advert{
				ID:          uint(returningAdvertList.Adverts[i].Advert.Id),
				UserID:      uint(returningAdvertList.Adverts[i].Advert.UserId),
				CityID:      uint(returningAdvertList.Adverts[i].Advert.CityId),
				CategoryID:  uint(returningAdvertList.Adverts[i].Advert.CategoryId),
				Title:       returningAdvertList.Adverts[i].Advert.Title,
				Description: returningAdvertList.Adverts[i].Advert.Description,
				Price:       uint(returningAdvertList.Adverts[i].Advert.Price),
				CreatedTime: returningAdvertList.Adverts[i].Advert.CreateTime.AsTime(),
				ClosedTime:  returningAdvertList.Adverts[i].Advert.CloseTime.AsTime(),
				Active:      returningAdvertList.Adverts[i].Advert.Active,
				IsUsed:      returningAdvertList.Adverts[i].Advert.IsUsed,
				Deleted:     returningAdvertList.Adverts[i].Advert.IsUsed,
			},
			City: models.City{
				ID:          uint(returningAdvertList.Adverts[i].City.Id),
				CityName:    returningAdvertList.Adverts[i].City.CityName,
				Translation: returningAdvertList.Adverts[i].City.Translation,
			},
			Category: models.Category{
				ID:          uint(returningAdvertList.Adverts[i].Category.Id),
				Name:        returningAdvertList.Adverts[i].Category.Name,
				Translation: returningAdvertList.Adverts[i].City.Translation,
			},
			Photos:    returningAdvertList.Adverts[i].Photos,
			PhotosIMG: returningAdvertList.Adverts[i].PhotosIMG,
		}

		returningAdvertsList = append(returningAdvertsList, &returningAdvert)
	}

	if returningAdvertsList == nil {
		returningAdvertsList = []*models.ReturningAdvert{}
	}

	return returningAdvertsList
}

func newProtobufAdvertList(advertList []*models.ReturningAdvert) *protobuf.ReturningAdvertList {

	var protobufAdvertsList protobuf.ReturningAdvertList

	for i := 0; i < len(advertList); i++ {
		protobufAdvert := protobuf.Advert{
			Id:          uint32(advertList[i].Advert.ID),
			UserId:      uint32(advertList[i].Advert.UserID),
			CityId:      uint32(advertList[i].Advert.CityID),
			CategoryId:  uint32(advertList[i].Advert.CategoryID),
			Title:       advertList[i].Advert.Title,
			Description: advertList[i].Advert.Description,
			Price:       uint32(advertList[i].Advert.Price),
			CreateTime:  timestamppb.New(advertList[i].Advert.CreatedTime),
			CloseTime:   timestamppb.New(advertList[i].Advert.ClosedTime),
			Active:      advertList[i].Advert.Active,
			IsUsed:      advertList[i].Advert.IsUsed,
		}

		protobufCity := protobuf.City{
			Id:          uint32(advertList[i].City.ID),
			CityName:    advertList[i].Category.Name,
			Translation: advertList[i].City.Translation,
		}

		protobufCategory := protobuf.Category{
			Id:          uint32(advertList[i].Category.ID),
			Name:        advertList[i].Category.Name,
			Translation: advertList[i].Category.Translation,
		}

		protobufAdvertItem := protobuf.ReturningAdvert{
			Advert:    &protobufAdvert,
			City:      &protobufCity,
			Category:  &protobufCategory,
			Photos:    advertList[i].Photos,
			PhotosIMG: advertList[i].PhotosIMG,
		}

		protobufAdvertsList.Adverts = append(protobufAdvertsList.Adverts, &protobufAdvertItem)

	}

	return &protobufAdvertsList
}

func (manager *CartManager) GetCartByUserID(ctx context.Context,
	in *protobuf.UserIdRequest) (*protobuf.ReturningAdvertList, error) {
	id := in.GetUserId()
	storage := manager.cartStorage

	cart, err := storage.GetCartByUserID(ctx, uint(id))
	if err != nil {
		return nil, err
	}

	return newProtobufAdvertList(cart), nil
}

func (manager *CartManager) DeleteAdvByIDs(ctx context.Context,
	in *protobuf.UserIdAdvertIdRequest) (*protobuf.DeleteAdvResponse, error) {
	userId := in.GetUserId()
	cartId := in.GetAdvertId()
	storage := manager.cartStorage

	err := storage.DeleteAdvByIDs(ctx, uint(userId), uint(cartId))
	if err != nil {
		return nil, err
	}

	return &protobuf.DeleteAdvResponse{IsAppended: false}, nil
}

func (manager *CartManager) AppendAdvByIDs(ctx context.Context,
	in *protobuf.UserIdAdvertIdRequest) (*protobuf.AppendAdvResponse, error) {
	userId := in.GetUserId()
	cartId := in.GetAdvertId()
	storage := manager.cartStorage

	isAppended := storage.AppendAdvByIDs(ctx, uint(userId), uint(cartId))

	return &protobuf.AppendAdvResponse{IsAppended: isAppended}, nil
}
