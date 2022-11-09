package texthandler_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter/texthandler"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

func TestSetDefaultCurrencyConvertTextToCommand(t *testing.T) {
	t.Parallel()

	var handler texthandler.SetDefaultCurrency

	type testCase struct {
		description string
		textInput   string
		matched     bool
		cmdExpected usecase.Command
	}

	testCases := [...]testCase{
		{
			description: "empty input",
			textInput:   "",
			matched:     false,
			cmdExpected: usecase.Command{},
		},
		{
			description: "command only",
			textInput:   "валюта",
			matched:     false,
			cmdExpected: usecase.Command{},
		},
		{
			description: "valid command",
			textInput:   "валюта RUB",
			matched:     true,
			cmdExpected: usecase.Command{
				SetDefaultCurrencyReqDTO: &usecase.SetDefaultCurrencyReqDTO{
					Currency: "RUB",
					UserID:   0,
				},
			},
		},
		{
			description: "invalid request",
			textInput:   "валюта RUB RUB",
			matched:     false,
			cmdExpected: usecase.Command{},
		},
	}

	for _, scenario := range testCases {
		scenario := scenario
		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			var cmd usecase.Command

			matched := handler.ConvertTextToCommand(ctx, scenario.textInput, &cmd)
			assert.EqualValues(t, scenario.matched, matched)
			assert.EqualValues(t, scenario.cmdExpected, cmd)
		})
	}
}

func TestSetDefaultCurrencyConvertCommandToText(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description  string
		cmd          usecase.Command
		textExpected string
		errExpected  string
	}

	date := time.Now()

	testCases := [...]testCase{
		{
			description:  "empty req",
			cmd:          usecase.Command{},
			textExpected: "",
			errExpected:  "SetDefaultCurrency.ExecuteCommand: internal error",
		},
		{
			description: "USD",
			cmd: usecase.Command{
				MessageInfo: usecase.MessageInfo{
					UserID: 101,
					Date:   date,
				},
				SetDefaultCurrencyReqDTO: &usecase.SetDefaultCurrencyReqDTO{
					UserID:   101,
					Currency: "USD",
				},
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
