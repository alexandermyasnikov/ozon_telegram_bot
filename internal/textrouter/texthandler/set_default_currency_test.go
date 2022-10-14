package texthandler_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter/texthandler"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter/texthandler/mock_texthandler"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

func TestSetDefaultCurrencyConvertTextToCommand(t *testing.T) {
	t.Parallel()

	userID := int64(101)
	date := time.Date(2022, 9, 20, 0, 0, 0, 0, time.UTC)

	var handler texthandler.SetDefaultCurrency

	type testCase struct {
		description string
		textInput   string
		cmdExpected textrouter.Command
	}

	testCases := [...]testCase{
		{
			description: "empty input",
			textInput:   "",
			cmdExpected: textrouter.Command{
				SetDefaultCurrencyReqDTO: nil,
				AddExpenseReqDTO:         nil,
				GetReportReqDTO:          nil,
				GetReportRespDTO:         nil,
			},
		},
		{
			description: "command only",
			textInput:   "валюта",
			cmdExpected: textrouter.Command{
				SetDefaultCurrencyReqDTO: nil,
				AddExpenseReqDTO:         nil,
				GetReportReqDTO:          nil,
				GetReportRespDTO:         nil,
			},
		},
		{
			description: "valid command",
			textInput:   "валюта RUB",
			cmdExpected: textrouter.Command{
				SetDefaultCurrencyReqDTO: &usecase.SetDefaultCurrencyReqDTO{
					UserID:   userID,
					Currency: "RUB",
				},
				AddExpenseReqDTO: nil,
				GetReportReqDTO:  nil,
				GetReportRespDTO: nil,
			},
		},
		{
			description: "invalid request",
			textInput:   "валюта",
			cmdExpected: textrouter.Command{
				SetDefaultCurrencyReqDTO: nil,
				AddExpenseReqDTO:         nil,
				GetReportReqDTO:          nil,
				GetReportRespDTO:         nil,
			},
		},
		{
			description: "invalid request",
			textInput:   "валюта RUB RUB",
			cmdExpected: textrouter.Command{
				SetDefaultCurrencyReqDTO: nil,
				AddExpenseReqDTO:         nil,
				GetReportReqDTO:          nil,
				GetReportRespDTO:         nil,
			},
		},
	}

	for _, scenario := range testCases {
		scenario := scenario
		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

			var cmd textrouter.Command

			handler.ConvertTextToCommand(userID, scenario.textInput, date, &cmd)
			assert.EqualValues(t, scenario.cmdExpected, cmd)
		})
	}
}

func TestSetDefaultCurrencyConvertExecuteCommand(t *testing.T) {
	t.Parallel()

	userID := int64(101)

	type testCase struct {
		description string
		cmd         textrouter.Command

		mockSetDefaultCurrencyReqDTO usecase.SetDefaultCurrencyReqDTO
		mockErr                      error
		mockTimes                    int

		cmdExpected textrouter.Command
		errExpected string
	}

	// TODO выглядит сложно
	testCases := [...]testCase{
		{
			description: "empty req",
			cmd: textrouter.Command{
				SetDefaultCurrencyReqDTO: nil,
				AddExpenseReqDTO:         nil,
				GetReportReqDTO:          nil,
				GetReportRespDTO:         nil,
			},

			mockSetDefaultCurrencyReqDTO: usecase.SetDefaultCurrencyReqDTO{
				UserID:   userID,
				Currency: "",
			},
			mockErr:   nil,
			mockTimes: 0,

			cmdExpected: textrouter.Command{
				SetDefaultCurrencyReqDTO: nil,
				AddExpenseReqDTO:         nil,
				GetReportReqDTO:          nil,
				GetReportRespDTO:         nil,
			},
			errExpected: "SetDefaultCurrency.ExecuteCommand: internal error",
		},
		{
			description: "with error",
			cmd: textrouter.Command{
				SetDefaultCurrencyReqDTO: &usecase.SetDefaultCurrencyReqDTO{
					UserID:   userID,
					Currency: "EUR",
				},
				AddExpenseReqDTO: nil,
				GetReportReqDTO:  nil,
				GetReportRespDTO: nil,
			},

			mockSetDefaultCurrencyReqDTO: usecase.SetDefaultCurrencyReqDTO{
				UserID:   userID,
				Currency: "EUR",
			},
			mockErr:   errors.New("invalid currency"),
			mockTimes: 1,

			cmdExpected: textrouter.Command{
				SetDefaultCurrencyReqDTO: &usecase.SetDefaultCurrencyReqDTO{
					UserID:   userID,
					Currency: "EUR",
				},
				AddExpenseReqDTO: nil,
				GetReportReqDTO:  nil,
				GetReportRespDTO: nil,
			},
			errExpected: "SetDefaultCurrency.ExecuteCommand: invalid currency",
		},
		{
			description: "without error",
			cmd: textrouter.Command{
				SetDefaultCurrencyReqDTO: &usecase.SetDefaultCurrencyReqDTO{
					UserID:   userID,
					Currency: "EUR",
				},
				AddExpenseReqDTO: nil,
				GetReportReqDTO:  nil,
				GetReportRespDTO: nil,
			},

			mockSetDefaultCurrencyReqDTO: usecase.SetDefaultCurrencyReqDTO{
				UserID:   userID,
				Currency: "EUR",
			},
			mockErr:   nil,
			mockTimes: 1,

			cmdExpected: textrouter.Command{
				SetDefaultCurrencyReqDTO: &usecase.SetDefaultCurrencyReqDTO{
					UserID:   userID,
					Currency: "EUR",
				},
				AddExpenseReqDTO: nil,
				GetReportReqDTO:  nil,
				GetReportRespDTO: nil,
			},
			errExpected: "",
		},
	}

	for _, scenario := range testCases {
		scenario := scenario
		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			expenseUsecase := mock_texthandler.NewMockIExpenseUsecaseSDC(ctrl)

			handler := texthandler.NewSetDefaultCurrency(expenseUsecase)

			expenseUsecase.EXPECT().SetDefaultCurrency(gomock.Any(), gomock.Eq(scenario.mockSetDefaultCurrencyReqDTO)).
				Return(scenario.mockErr).Times(scenario.mockTimes)

			ctx := context.Background()
			err := handler.ExecuteCommand(ctx, &scenario.cmd)
			if len(scenario.errExpected) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, scenario.errExpected)
			}

			assert.EqualValues(t, scenario.cmdExpected, scenario.cmd)
		})
	}
}

func TestSetDefaultCurrencyConvertCommandToText(t *testing.T) {
	t.Parallel()

	userID := int64(101)

	type testCase struct {
		description  string
		cmd          textrouter.Command
		textExpected string
		errExpected  string
	}

	testCases := [...]testCase{
		{
			description: "empty req",
			cmd: textrouter.Command{
				SetDefaultCurrencyReqDTO: nil,
				AddExpenseReqDTO:         nil,
				GetReportReqDTO:          nil,
				GetReportRespDTO:         nil,
			},
			textExpected: "",
			errExpected:  "SetDefaultCurrency.ExecuteCommand: internal error",
		},
		{
			description: "USD",
			cmd: textrouter.Command{
				SetDefaultCurrencyReqDTO: &usecase.SetDefaultCurrencyReqDTO{
					UserID:   userID,
					Currency: "USD",
				},
				AddExpenseReqDTO: nil,
				GetReportReqDTO:  nil,
				GetReportRespDTO: nil,
			},
			textExpected: "Задана валюта по умолчанию USD",
			errExpected:  "",
		},
	}

	for _, scenario := range testCases {
		scenario := scenario
		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

			var handler texthandler.SetDefaultCurrency

			textOutput, err := handler.ConvertCommandToText(&scenario.cmd)
			assert.Equal(t, scenario.textExpected, textOutput)
			if len(scenario.errExpected) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, scenario.errExpected)
			}
		})
	}
}
