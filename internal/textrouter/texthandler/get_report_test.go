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

func TestGetReportConvertTextToCommand(t *testing.T) {
	t.Parallel()

	userID := int64(101)
	date := time.Date(2022, 9, 20, 0, 0, 0, 0, time.UTC)

	var handler texthandler.GetReport

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
			textInput:   "отчет",
			cmdExpected: textrouter.Command{},
		},
		{
			description: "day",
			textInput:   "отчет день",
			cmdExpected: textrouter.Command{
				GetReportReqDTO: &usecase.GetReportReqDTO{
					UserID:       userID,
					Date:         date,
					IntervalType: utils.DayInterval,
				},
			},
		},
		{
			description: "week",
			textInput:   "отчет неделя",
			cmdExpected: textrouter.Command{
				GetReportReqDTO: &usecase.GetReportReqDTO{
					UserID:       userID,
					Date:         date,
					IntervalType: utils.WeekInterval,
				},
			},
		},
		{
			description: "month",
			textInput:   "отчет месяц",
			cmdExpected: textrouter.Command{
				GetReportReqDTO: &usecase.GetReportReqDTO{
					UserID:       userID,
					Date:         date,
					IntervalType: utils.MonthInterval,
				},
			},
		},
		{
			description: "month with spaces",
			textInput:   "  отчет             месяц       ",
			cmdExpected: textrouter.Command{
				GetReportReqDTO: &usecase.GetReportReqDTO{
					UserID:       userID,
					Date:         date,
					IntervalType: utils.MonthInterval,
				},
			},
		},
		{
			description: "invalid request",
			textInput:   "отчет нед",
			cmdExpected: textrouter.Command{},
		},
		{
			description: "invalid request",
			textInput:   "статистика",
			cmdExpected: textrouter.Command{},
		},
		{
			description: "invalid request",
			textInput:   "отчет нед RUB RUB",
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

func TestGetReportConvertExecuteCommand(t *testing.T) {
	t.Parallel()

	userID := int64(101)
	date := time.Date(2022, 9, 20, 0, 0, 0, 0, time.UTC)

	type testCase struct {
		description string
		cmd         textrouter.Command

		mockGetReportReqDTO  usecase.GetReportReqDTO
		mockGetReportRespDTO usecase.GetReportRespDTO
		mockErr              error
		mockTimes            int

		cmdExpected textrouter.Command
		errExpected string
	}

	testCases := [...]testCase{
		{
			description: "empty req",
			cmd:         textrouter.Command{},

			mockGetReportReqDTO: usecase.GetReportReqDTO{
				UserID:       userID,
				Date:         date,
				IntervalType: utils.DayInterval,
			},
			mockGetReportRespDTO: usecase.GetReportRespDTO{
				Currency: "",
				Expenses: nil,
			},
			mockErr:   nil,
			mockTimes: 0,

			cmdExpected: textrouter.Command{},
			errExpected: "GetReport.ExecuteCommand: internal error",
		},
		{
			description: "with error",
			cmd: textrouter.Command{
				GetReportReqDTO: &usecase.GetReportReqDTO{
					UserID:       userID,
					Date:         date,
					IntervalType: utils.MonthInterval,
				},
			},

			mockGetReportReqDTO: usecase.GetReportReqDTO{
				UserID:       userID,
				Date:         date,
				IntervalType: utils.MonthInterval,
			},
			mockGetReportRespDTO: usecase.GetReportRespDTO{
				Currency: "",
				Expenses: nil,
			},
			mockErr:   errors.New("unknown error"),
			mockTimes: 1,

			cmdExpected: textrouter.Command{
				GetReportReqDTO: &usecase.GetReportReqDTO{
					UserID:       userID,
					Date:         date,
					IntervalType: utils.MonthInterval,
				},
			},
			errExpected: "GetReport.ExecuteCommand: unknown error",
		},
		{
			description: "without error",
			cmd: textrouter.Command{
				GetReportReqDTO: &usecase.GetReportReqDTO{
					UserID:       userID,
					Date:         date,
					IntervalType: utils.MonthInterval,
				},
			},

			mockGetReportReqDTO: usecase.GetReportReqDTO{
				UserID:       userID,
				Date:         date,
				IntervalType: utils.MonthInterval,
			},
			mockGetReportRespDTO: usecase.GetReportRespDTO{
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
			mockErr:   nil,
			mockTimes: 1,

			cmdExpected: textrouter.Command{
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
			errExpected: "",
		},
	}

	for _, scenario := range testCases {
		scenario := scenario
		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			expenseUsecase := mock_texthandler.NewMockExpenseUsecaseGR(ctrl)

			handler := texthandler.NewGetReport(expenseUsecase)

			expenseUsecase.EXPECT().GetReport(gomock.Any(), gomock.Eq(scenario.mockGetReportReqDTO)).
				Return(scenario.mockGetReportRespDTO, scenario.mockErr).Times(scenario.mockTimes)

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

func TestGetReportConvertCommandToText(t *testing.T) {
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
			cmd: textrouter.Command{
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
			cmd: textrouter.Command{
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
