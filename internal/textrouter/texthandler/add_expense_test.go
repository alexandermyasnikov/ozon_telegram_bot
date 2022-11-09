package texthandler_test

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter/texthandler"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/utils"
)

func TestAddExpenseConvertTextToCommand(t *testing.T) {
	t.Parallel()

	date := time.Now()

	var handler texthandler.AddExpense

	type testCase struct {
		description string
		textInput   string
		matched     bool
		cmdBefore   usecase.Command
		cmdAfter    usecase.Command
	}

	testCases := [...]testCase{
		{
			description: "empty input",
			textInput:   "",
			matched:     false,
		},
		{
			description: "command only",
			textInput:   "расход",
			matched:     false,
		},
		{
			description: "category",
			textInput:   "расход категория1",
			matched:     false,
		},
		{
			description: "category + price",
			textInput:   "расход категория1 1234.45678",
			matched:     true,
			cmdBefore: usecase.Command{
				MessageInfo: usecase.MessageInfo{
					UserID: 101,
					Date:   date,
				},
			},
			cmdAfter: usecase.Command{
				MessageInfo: usecase.MessageInfo{
					UserID: 101,
					Date:   date,
				},
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   101,
					Category: "категория1",
					Price:    decimal.RequireFromString("1234.45678"),
					Date:     date,
				},
			},
		},
		{
			description: "invalid request",
			textInput:   "расход категория1 1234.45678 EUR tmp",
			matched:     false,
		},
	}

	for _, scenario := range testCases {
		scenario := scenario
		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			cmd := scenario.cmdBefore

			matched := handler.ConvertTextToCommand(ctx, scenario.textInput, &cmd)
			assert.EqualValues(t, scenario.matched, matched)
			assert.EqualValues(t, scenario.cmdAfter, cmd)
		})
	}
}

func TestAddExpenseConvertCommandToText(t *testing.T) {
	t.Parallel()

	userID := int64(101)
	date := time.Date(2022, 9, 20, 0, 0, 0, 0, time.UTC)

	type testCase struct {
		description  string
		cmd          usecase.Command
		textExpected string
		errExpected  string
	}

	testCases := [...]testCase{
		{
			description:  "empty req",
			cmd:          usecase.Command{},
			textExpected: "",
			errExpected:  "AddExpense.ExecuteCommand: internal error",
		},
		{
			description: "category + price",
			cmd: usecase.Command{
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   userID,
					Category: "Category2",
					Price:    decimal.RequireFromString("43.5678"),
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
			cmd: usecase.Command{
				AddExpenseReqDTO: &usecase.AddExpenseReqDTO{
					UserID:   userID,
					Category: "Category2",
					Price:    decimal.RequireFromString("43.5678"),
					Date:     date,
				},
				AddExpenseRespDTO: &usecase.AddExpenseRespDTO{
					Currency: "USD",
					Limits: map[int]decimal.Decimal{
						utils.DayInterval:   decimal.New(0, 0),
						utils.WeekInterval:  decimal.RequireFromString("-12345.678"),
						utils.MonthInterval: decimal.RequireFromString("34.5678"),
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

			ctx := context.Background()

			var handler texthandler.AddExpense

			textOutput, err := handler.ConvertCommandToText(ctx, &scenario.cmd)
			assert.Equal(t, scenario.textExpected, textOutput)
			if len(scenario.errExpected) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, scenario.errExpected)
			}
		})
	}
}
