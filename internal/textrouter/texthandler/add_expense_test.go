package texthandler_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter/texthandler"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter/texthandler/mock_texthandler"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/utils"
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
			cmdExpected: textrouter.Command{},
		},
		{
			description: "command only",
			textInput:   "расход",
			cmdExpected: textrouter.Command{},
		},
		{
			description: "category",
			textInput:   "расход категория1",
			cmdExpected: textrouter.Command{},
		},
		{
			description: "category + price",
			textInput:   "расход категория1 1234.45678",
			cmdExpected: textrouter.Command{
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   userID,
					Category: "категория1",
					Price:    decimal.New(123445678, -5),
					Date:     date,
				},
			},
		},
		{
			description: "invalid request",
			textInput:   "расход категория1 1234.45678 EUR tmp",
			cmdExpected: textrouter.Command{},
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

		mockAddExpenseReqDTO  usecase.AddExpenseReqDTO
		mockAddExpenseRespDTO usecase.AddExpenseRespDTO
		mockErr               error
		mockTimes             int

		cmdExpected textrouter.Command
		errExpected string
	}

	testCases := [...]testCase{
		{
			description: "empty req",
			cmd:         textrouter.Command{},

			mockAddExpenseReqDTO: usecase.AddExpenseReqDTO{
				UserID:   userID,
				Category: "",
				Price:    decimal.Zero,
				Date:     date,
			},
			mockErr:   nil,
			mockTimes: 0,

			cmdExpected: textrouter.Command{},
			errExpected: "AddExpense.ExecuteCommand: internal error",
		},
		{
			description: "with error",
			cmd: textrouter.Command{
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   userID,
					Category: "Category2",
					Price:    decimal.New(435678, -4),
					Date:     date,
				},
			},

			mockAddExpenseReqDTO: usecase.AddExpenseReqDTO{
				UserID:   userID,
				Category: "Category2",
				Price:    decimal.New(435678, -4),
				Date:     date,
			},
			mockErr:   errors.New("unknown error"),
			mockTimes: 1,

			cmdExpected: textrouter.Command{
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   userID,
					Category: "Category2",
					Price:    decimal.New(435678, -4),
					Date:     date,
				},
			},
			errExpected: "AddExpense.ExecuteCommand: unknown error",
		},
		{
			description: "without error",
			cmd: textrouter.Command{
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   userID,
					Category: "Category2",
					Price:    decimal.New(435678, -4),
					Date:     date,
				},
			},

			mockAddExpenseReqDTO: usecase.AddExpenseReqDTO{
				UserID:   userID,
				Category: "Category2",
				Price:    decimal.New(435678, -4),
				Date:     date,
			},
			mockAddExpenseRespDTO: usecase.AddExpenseRespDTO{
				Currency: "USD",
				Limits: map[int]decimal.Decimal{
					utils.WeekInterval: decimal.New(123, -1),
				},
			},
			mockErr:   nil,
			mockTimes: 1,

			cmdExpected: textrouter.Command{
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   userID,
					Category: "Category2",
					Price:    decimal.New(435678, -4),
					Date:     date,
				},
				AddExpenseRespDTO: &usecase.AddExpenseRespDTO{
					Currency: "USD",
					Limits: map[int]decimal.Decimal{
						utils.WeekInterval: decimal.New(123, -1),
					},
				},
			},
			errExpected: "",
		},
	}

	for _, scenario := range testCases {
		scenario := scenario
		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			expenseUsecase := mock_texthandler.NewMockExpenseUsecaseAE(ctrl)

			handler := texthandler.NewAddExpense(expenseUsecase)

			expenseUsecase.EXPECT().AddExpense(gomock.Any(), gomock.Eq(scenario.mockAddExpenseReqDTO)).
				Return(scenario.mockAddExpenseRespDTO, scenario.mockErr).Times(scenario.mockTimes)

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
			description:  "empty req",
			cmd:          textrouter.Command{},
			textExpected: "",
			errExpected:  "AddExpense.ExecuteCommand: internal error",
		},
		{
			description: "category + price",
			cmd: textrouter.Command{
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   userID,
					Category: "Category2",
					Price:    decimal.New(435678, -4),
					Date:     date,
				},
				AddExpenseRespDTO: &usecase.AddExpenseRespDTO{
					Currency: "EUR",
					Limits:   nil,
				},
			},
			textExpected: "Добавил Category2 - 43.57 EUR Tue, 20 Sep 2022 00:00:00 UTC",
			errExpected:  "",
		},
		{
			description: "category + limits",
			cmd: textrouter.Command{
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   userID,
					Category: "Category2",
					Price:    decimal.New(435678, -4),
					Date:     date,
				},
				AddExpenseRespDTO: &usecase.AddExpenseRespDTO{
					Currency: "USD",
					Limits: map[int]decimal.Decimal{
						utils.DayInterval:   decimal.New(0, 0),
						utils.WeekInterval:  decimal.New(-12345678, -3),
						utils.MonthInterval: decimal.New(345678, -4),
					},
				},
			},
			textExpected: `Добавил Category2 - 43.57 USD Tue, 20 Sep 2022 00:00:00 UTC
Внимание! Превышен лимит: неделя - 12345.68`,
			errExpected: "",
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
