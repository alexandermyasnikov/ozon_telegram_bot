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

func TestGetReportConvertTextToCommand(t *testing.T) {
	t.Parallel()

	date := time.Now()

	var handler texthandler.GetReport

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
			cmdBefore:   usecase.Command{},
			cmdAfter:    usecase.Command{},
		},
		{
			description: "command only",
			textInput:   "отчет",
			matched:     false,
			cmdBefore:   usecase.Command{},
			cmdAfter:    usecase.Command{},
		},
		{
			description: "day",
			textInput:   "отчет день",
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
				GetReportReqDTO: &usecase.GetReportReqDTO{
					UserID:       101,
					Date:         date,
					IntervalType: utils.DayInterval,
				},
			},
		},
		{
			description: "week",
			textInput:   "отчет неделя",
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
				GetReportReqDTO: &usecase.GetReportReqDTO{
					UserID:       101,
					Date:         date,
					IntervalType: utils.WeekInterval,
				},
			},
		},
		{
			description: "month with spaces",
			textInput:   "  отчет             месяц       ",
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
				GetReportReqDTO: &usecase.GetReportReqDTO{
					UserID:       101,
					Date:         date,
					IntervalType: utils.MonthInterval,
				},
			},
		},
		{
			description: "invalid request",
			textInput:   "отчет нед",
			matched:     false,
			cmdBefore:   usecase.Command{},
			cmdAfter:    usecase.Command{},
		},
		{
			description: "invalid request",
			textInput:   "статистика",
			matched:     false,
			cmdBefore:   usecase.Command{},
			cmdAfter:    usecase.Command{},
		},
		{
			description: "invalid request",
			textInput:   "отчет нед RUB RUB",
			matched:     false,
			cmdBefore:   usecase.Command{},
			cmdAfter:    usecase.Command{},
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

func TestGetReportConvertCommandToText(t *testing.T) {
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
			description: "empty req",
			cmd: usecase.Command{
				GetReportRespDTO: &usecase.GetReportRespDTO{
					Currency: "RUB",
					Expenses: []usecase.ExpenseReportDTO{
						{
							Category: "Catergory1",
							Sum:      decimal.New(12, 0),
						},
						{
							Category: "Catergory2",
							Sum:      decimal.New(34567, -3),
						},
					},
				},
			},
			textExpected: "",
			errExpected:  "GetReport.ExecuteCommand: internal error",
		},
		{
			description: "empty resp",
			cmd: usecase.Command{
				GetReportReqDTO: &usecase.GetReportReqDTO{
					UserID:       userID,
					Date:         date,
					IntervalType: utils.MonthInterval,
				},
			},
			textExpected: "",
			errExpected:  "GetReport.ExecuteCommand: internal error",
		},
		{
			description: "categories",
			cmd: usecase.Command{
				GetReportReqDTO: &usecase.GetReportReqDTO{
					UserID:       userID,
					Date:         date,
					IntervalType: utils.MonthInterval,
				},
				GetReportRespDTO: &usecase.GetReportRespDTO{
					Currency: "RUB",
					Expenses: []usecase.ExpenseReportDTO{
						{
							Category: "Catergory1",
							Sum:      decimal.New(12, 0),
						},
						{
							Category: "Catergory2",
							Sum:      decimal.New(34567, -3),
						},
					},
				},
			},
			textExpected: "Расходы по категориям за месяц:\n" +
				"Catergory1 - 12.00\n" +
				"Catergory2 - 34.57",
			errExpected: "",
		},
	}

	for _, scenario := range testCases {
		scenario := scenario
		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

			var handler texthandler.GetReport

			ctx := context.Background()

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
