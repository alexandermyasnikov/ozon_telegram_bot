package texthandler

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/utils"
)

type SetLimit struct{}

func NewSetLimit() *SetLimit {
	return &SetLimit{}
}

func (h *SetLimit) Name() string {
	return usecase.SetLimitCmdName
}

func (h *SetLimit) ConvertTextToCommand(ctx context.Context, text string, cmd *usecase.Command) bool {
	intervalIndex := 1
	limitIndex := 2
	argsCountMin := 3
	argsCountMax := 3

	fields := strings.Fields(text)
	if len(fields) < argsCountMin || len(fields) > argsCountMax || fields[0] != "лимит" {
		return false
	}

	intervalType, ok := utils.IntervalFromStr(fields[intervalIndex])
	if !ok {
		return false
	}

	limit, err := decimal.NewFromString(fields[limitIndex])
	if err != nil {
		return false
	}

	cmd.SetLimitReqDTO = &usecase.SetLimitReqDTO{
		UserID:       cmd.UserID,
		Limit:        limit,
		IntervalType: intervalType,
	}

	return true
}

func (h *SetLimit) ConvertCommandToText(ctx context.Context, cmd *usecase.Command) (string, error) {
	if cmd.SetLimitReqDTO == nil || cmd.SetLimitRespDTO == nil {
		return "", errors.Wrap(textrouter.ErrInvalidCommand, "SetLimit.ExecuteCommand")
	}

	precision := 2

	intervalType, _ := utils.IntervalToStr(cmd.SetLimitReqDTO.IntervalType)

	textOut := fmt.Sprintf("Установил лимит: %s - %s - %s",
		intervalType,
		cmd.SetLimitReqDTO.Limit.StringFixed(int32(precision)),
		cmd.SetLimitRespDTO.Currency)

	return textOut, nil
}
