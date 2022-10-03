package router

import (
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

var _ HandlerTextInterface = (*HandlerTextAddProduct)(nil)

var (
	reNamePrice = regexp.MustCompile(genRE(
		[]string{"cmd", `сохр`, `запис`},
		[]string{"name", `\S+`},
		[]string{"price", `\d+`},
	))

	reNamePriceDate = regexp.MustCompile(genRE(
		[]string{"cmd", `сохр`, `запис`},
		[]string{"name", `\S+`},
		[]string{"price", `\d+`},
		[]string{"date", `ден`, `нед`, `мес`, `год`},
	))
)

type HandlerTextAddProduct struct{}

func (h *HandlerTextAddProduct) GetID() int {
	return cmdAddProduct
}

func (h *HandlerTextAddProduct) ConvertTextToCommand(userID int64, text string, date time.Time, cmd *command) {
	if matches := reNamePrice.FindStringSubmatch(text); len(matches) > 0 {
		category := matches[reNamePrice.SubexpIndex("name")]
		price, _ := strconv.Atoi(matches[reNamePrice.SubexpIndex("price")])

		cmd.id = cmdAddProduct
		cmd.addProductReqDTO = &usecase.AddProductReqDTO{
			UserID: userID,
			Product: usecase.ProductDTO{
				Category: category,
				Price:    price,
				Date:     date,
			},
		}

		return
	}

	if matches := reNamePriceDate.FindStringSubmatch(text); len(matches) > 0 {
		category := matches[reNamePriceDate.SubexpIndex("name")]
		price, _ := strconv.Atoi(matches[reNamePriceDate.SubexpIndex("price")])
		days := dateToDays(matches[reNamePriceDate.SubexpIndex("date")])
		date := date.AddDate(0, 0, days)

		cmd.id = cmdAddProduct
		cmd.addProductReqDTO = &usecase.AddProductReqDTO{
			UserID: userID,
			Product: usecase.ProductDTO{
				Category: category,
				Price:    price,
				Date:     date,
			},
		}

		return
	}
}

func (h *HandlerTextAddProduct) ExecuteCommand(cmd *command, productUsecase usecase.ProductUsecaseInterface) error {
	if cmd.addProductReqDTO == nil {
		return errors.New("empty addProductReqDTO")
	}

	err := productUsecase.AddProduct(*cmd.addProductReqDTO)

	return errors.Wrap(err, "productUsecase.AddProduct")
}

func (h *HandlerTextAddProduct) ConvertCommandToText(cmd command) string {
	return cmdAddProductText
}
