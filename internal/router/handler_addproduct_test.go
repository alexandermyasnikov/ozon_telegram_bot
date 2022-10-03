package router

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/mocks/usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

func TestGetID(t *testing.T) {
	t.Parallel()

	var handler HandlerTextAddProduct

	assert.Equal(t, cmdAddProduct, handler.GetID())
}

func TestConvertTextToCommand(t *testing.T) {
	t.Parallel()

	userID := int64(1)
	date := time.Date(2022, 9, 20, 10, 0, 0, 0, time.UTC)

	var handler HandlerTextAddProduct

	type testCase struct {
		description string
		textInput   string
		cmdExpected command
	}

	testCases := [...]testCase{
		{
			description: "empty input",
			textInput:   "",
			cmdExpected: command{
				id:                   cmdUnknown,
				addProductReqDTO:     nil,
				getStatisticsReqDTO:  nil,
				getStatisticsRespDTO: nil,
			},
		},
		{
			description: "addProduct command",
			textInput:   "записать продукты 1000",
			cmdExpected: command{
				id: cmdAddProduct,
				addProductReqDTO: &usecase.AddProductReqDTO{
					UserID: userID,
					Product: usecase.ProductDTO{
						Category: "продукты",
						Price:    1000,
						Date:     date,
					},
				},
				getStatisticsReqDTO:  nil,
				getStatisticsRespDTO: nil,
			},
		},
		{
			description: "addProduct command rub",
			textInput:   "записать продукты 1000руб",
			cmdExpected: command{
				id: cmdAddProduct,
				addProductReqDTO: &usecase.AddProductReqDTO{
					UserID: userID,
					Product: usecase.ProductDTO{
						Category: "продукты",
						Price:    1000,
						Date:     date,
					},
				},
				getStatisticsReqDTO:  nil,
				getStatisticsRespDTO: nil,
			},
		},
		/*{
			description: "addProduct command",
			textInput:   "записать продукты 1000 вчера",
			cmdExpected: command{
				id: cmdAddProduct,
				addProductReqDTO: &usecase.AddProductReqDTO{
					UserID: userID,
					Product: usecase.ProductDTO{
						Category: "продукты",
						Price:    1000,
						Date:     date.AddDate(0, 0, -1),
					},
				},
				getStatisticsReqDTO:  nil,
				getStatisticsRespDTO: nil,
			},
		},*/
	}

	for _, scenario := range testCases {
		scenario := scenario
		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

			cmd := command{
				id:                   cmdUnknown,
				addProductReqDTO:     nil,
				getStatisticsReqDTO:  nil,
				getStatisticsRespDTO: nil,
			}

			handler.ConvertTextToCommand(userID, scenario.textInput, date, &cmd)
			assert.Equal(t, scenario.cmdExpected.id, cmd.id)
			assert.Equal(t, scenario.cmdExpected.addProductReqDTO, cmd.addProductReqDTO)
			assert.Equal(t, scenario.cmdExpected.getStatisticsReqDTO, cmd.getStatisticsReqDTO)
			assert.Equal(t, scenario.cmdExpected.getStatisticsRespDTO, cmd.getStatisticsRespDTO)
		})
	}
}

func TestConvertExecuteCommandEmptyDTO(t *testing.T) {
	t.Parallel()

	cmdEmptyDTO := command{
		id:                   cmdAddProduct,
		addProductReqDTO:     nil,
		getStatisticsReqDTO:  nil,
		getStatisticsRespDTO: nil,
	}

	var handler HandlerTextAddProduct

	ctrl := gomock.NewController(t)
	sender := mocks.NewMockProductUsecaseInterface(ctrl)

	err := handler.ExecuteCommand(&cmdEmptyDTO, sender)
	assert.Error(t, err)
}

func TestConvertExecuteCommand(t *testing.T) {
	t.Parallel()

	userID := int64(1)
	date := time.Date(2022, 9, 20, 10, 0, 0, 0, time.UTC)

	var handler HandlerTextAddProduct

	cmd := command{
		id: cmdAddProduct,
		addProductReqDTO: &usecase.AddProductReqDTO{
			UserID: userID,
			Product: usecase.ProductDTO{
				Category: "продукты",
				Price:    1000,
				Date:     date,
			},
		},
		getStatisticsReqDTO:  nil,
		getStatisticsRespDTO: nil,
	}

	ctrl := gomock.NewController(t)
	sender := mocks.NewMockProductUsecaseInterface(ctrl)

	sender.EXPECT().AddProduct(*cmd.addProductReqDTO)

	err := handler.ExecuteCommand(&cmd, sender)
	assert.NoError(t, err)
}

func TestConvertCommandToText(t *testing.T) {
	t.Parallel()

	userID := int64(1)
	date := time.Date(2022, 9, 20, 10, 0, 0, 0, time.UTC)

	var handler HandlerTextAddProduct

	cmd := command{
		id: cmdAddProduct,
		addProductReqDTO: &usecase.AddProductReqDTO{
			UserID: userID,
			Product: usecase.ProductDTO{
				Category: "продукты",
				Price:    1000,
				Date:     date,
			},
		},
		getStatisticsReqDTO:  nil,
		getStatisticsRespDTO: nil,
	}

	textActual := handler.ConvertCommandToText(cmd)
	assert.Equal(t, cmdAddProductText, textActual)
}
