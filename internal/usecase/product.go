package usecase

import (
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/storage"
)

type ProductUsecaseInterface interface {
	AddProduct(AddProductReqDTO) error
	GetStatistics(GetStatisticsReqDTO) (GetStatisticsRespDTO, error)
}

type ProductDTO struct {
	Category string
	Price    int
	Date     time.Time
}

type AddProductReqDTO struct {
	UserID  int64
	Product ProductDTO
}

type GetStatisticsReqDTO struct {
	UserID int64
	Date   time.Time
	Days   int
}

type GetStatisticsRespDTO struct {
	Products map[string]int
}

type ProductUsecase struct {
	productStorage storage.ProductStorageInterface
}

var _ ProductUsecaseInterface = (*ProductUsecase)(nil)

func NewProductUsecase(productStorage storage.ProductStorageInterface) *ProductUsecase {
	return &ProductUsecase{
		productStorage: productStorage,
	}
}

func (uc *ProductUsecase) AddProduct(req AddProductReqDTO) error {
	product := entity.NewProduct("", req.Product.Category, req.Product.Price, req.Product.Date)

	err := uc.productStorage.Add(req.UserID, product)

	return errors.Wrap(err, "productUsecase AddProduct")
}

func (uc *ProductUsecase) GetStatistics(req GetStatisticsReqDTO) (GetStatisticsRespDTO, error) {
	resp := GetStatisticsRespDTO{
		Products: make(map[string]int),
	}

	products, err := uc.productStorage.GetAll(req.UserID, req.Date, req.Days)
	if err != nil {
		return GetStatisticsRespDTO{}, errors.Wrap(err, "productUsecase GetAll")
	}

	for _, product := range products {
		resp.Products[product.GetCategory()] += product.GetPrice()
	}

	return resp, nil
}
