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

func TestAddExpenseConvertTextToCommand(t *testing.T) {
	t.Parallel()

	userID := int64(101)
	date := time.Date(2022, 9, 20, 0, 0, 0, 0, time.UTC)

	var handler texthandler.AddExpense

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
			textInput:   "расход",
			cmdExpected: textrouter.Command{
				SetDefaultCurrencyReqDTO: nil,
				AddExpenseReqDTO:         nil,
				GetReportReqDTO:          nil,
				GetReportRespDTO:         nil,
			},
		},
		{
			description: "category",
			textInput:   "расход категория1",
			cmdExpected: textrouter.Command{
				SetDefaultCurrencyReqDTO: nil,
				AddExpenseReqDTO:         nil,
				GetReportReqDTO:          nil,
				GetReportRespDTO:         nil,
			},
		},
		{
			description: "category + price",
			textInput:   "расход категория1 1234.45678",
			cmdExpected: textrouter.Command{
				SetDefaultCurrencyReqDTO: nil,
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   userID,
					Category: "категория1",
					Price:    1234.45678,
					Date:     date,
					Currency: "",
				},
				GetReportReqDTO:  nil,
				GetReportRespDTO: nil,
			},
		},
		{
			description: "category + price + currency",
			textInput:   "расход категория1 1234.45678 EUR",
			cmdExpected: textrouter.Command{
				SetDefaultCurrencyReqDTO: nil,
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   userID,
					Category: "категория1",
					Price:    1234.45678,
					Date:     date,
					Currency: "EUR",
				},
				GetReportReqDTO:  nil,
				GetReportRespDTO: nil,
			},
		},
		{
			description: "invalid request",
			textInput:   "расход категория1 1234.45678 EUR tmp",
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

func TestAddExpenseConvertExecuteCommand(t *testing.T) {
	t.Parallel()

	userID := int64(101)
	date := time.Date(2022, 9, 20, 0, 0, 0, 0, time.UTC)

	type testCase struct {
		description string
		cmd         textrouter.Command

		mockAddExpenseReqDTO usecase.AddExpenseReqDTO
		mockErr              error
		mockTimes            int

		cmdExpected textrouter.Command
		errExpected string
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

			mockAddExpenseReqDTO: usecase.AddExpenseReqDTO{
				UserID:   userID,
				Category: "",
				Price:    0,
				Date:     date,
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
			errExpected: "AddExpense.ExecuteCommand: internal error",
		},
		{
			description: "with error",
			cmd: textrouter.Command{
				SetDefaultCurrencyReqDTO: nil,
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   userID,
					Category: "Category2",
					Price:    43.5678,
					Date:     date,
					Currency: "ABC",
				},
				GetReportReqDTO:  nil,
				GetReportRespDTO: nil,
			},

			mockAddExpenseReqDTO: usecase.AddExpenseReqDTO{
				UserID:   userID,
				Category: "Category2",
				Price:    43.5678,
				Date:     date,
				Currency: "ABC",
			},
			mockErr:   errors.New("unknown currency"),
			mockTimes: 1,

			cmdExpected: textrouter.Command{
				SetDefaultCurrencyReqDTO: nil,
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   userID,
					Category: "Category2",
					Price:    43.5678,
					Date:     date,
					Currency: "ABC",
				},
				GetReportReqDTO:  nil,
				GetReportRespDTO: nil,
			},
			errExpected: "AddExpense.ExecuteCommand: unknown currency",
		},
		{
			description: "without error",
			cmd: textrouter.Command{
				SetDefaultCurrencyReqDTO: nil,
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   userID,
					Category: "Category2",
					Price:    43.5678,
					Date:     date,
					Currency: "USD",
				},
				GetReportReqDTO:  nil,
				GetReportRespDTO: nil,
			},

			mockAddExpenseReqDTO: usecase.AddExpenseReqDTO{
				UserID:   userID,
				Category: "Category2",
				Price:    43.5678,
				Date:     date,
				Currency: "USD",
			},
			mockErr:   nil,
			mockTimes: 1,

			cmdExpected: textrouter.Command{
				SetDefaultCurrencyReqDTO: nil,
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   userID,
					Category: "Category2",
					Price:    43.5678,
					Date:     date,
					Currency: "USD",
				},
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
			expenseUsecase := mock_texthandler.NewMockIExpenseUsecaseAE(ctrl)

			handler := texthandler.NewAddExpense(expenseUsecase)

			expenseUsecase.EXPECT().AddExpense(gomock.Any(), gomock.Eq(scenario.mockAddExpenseReqDTO)).
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

func TestAddExpenseConvertCommandToText(t *testing.T) {
	t.Parallel()

	userID := int64(101)
	date := time.Date(2022, 9, 20, 0, 0, 0, 0, time.UTC)

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
			errExpected:  "AddExpense.ExecuteCommand: internal error",
		},
		{
			description: "category + price",
			cmd: textrouter.Command{
				SetDefaultCurrencyReqDTO: nil,
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   userID,
					Category: "Category2",
					Price:    43.5678,
					Date:     date,
					Currency: "",
				},
				GetReportReqDTO:  nil,
				GetReportRespDTO: nil,
			},
			textExpected: "Добавил Category2 - 43.57  Tue, 20 Sep 2022 00:00:00 UTC",
			errExpected:  "",
		},
		{
			description: "category + price + currency",
			cmd: textrouter.Command{
				SetDefaultCurrencyReqDTO: nil,
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   userID,
					Category: "Category2",
					Price:    43.5678,
					Date:     date,
					Currency: "USD",
				},
				GetReportReqDTO:  nil,
				GetReportRespDTO: nil,
			},
			textExpected: "Добавил Category2 - 43.57 USD Tue, 20 Sep 2022 00:00:00 UTC",
			errExpected:  "",
		},
	}

	for _, scenario := range testCases {
		scenario := scenario
		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

			var handler texthandler.AddExpense

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
