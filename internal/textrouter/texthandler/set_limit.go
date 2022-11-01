package texthandler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/utils"
)

type ExpenseUsecaseSL interface {
	SetLimit(ctx context.Context, req usecase.SetLimitReqDTO) (usecase.SetLimitRespDTO, error)
}

type SetLimit struct {
	expenseUsecase ExpenseUsecaseSL
}

func NewSetLimit(expenseUsecase ExpenseUsecaseSL) *SetLimit {
	return &SetLimit{
		expenseUsecase: expenseUsecase,
	}
}

func (h *SetLimit) Name() string {
	return "setLimit"
}

func (h *SetLimit) ConvertTextToCommand(userID int64, text string, date time.Time, cmd *textrouter.Command) bool {
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
		UserID:       userID,
		Limit:        limit,
		IntervalType: intervalType,
	}

	return true
}

func (h *SetLimit) ExecuteCommand(ctx context.Context, cmd *textrouter.Command) error {
	if cmd.SetLimitReqDTO == nil {
		return errors.Wrap(textrouter.ErrInvalidCommand, "SetLimit.ExecuteCommand")
	}

	resp, err := h.expenseUsecase.SetLimit(ctx, *cmd.SetLimitReqDTO)
	if err != nil {
		return errors.Wrap(err, "SetLimit.ExecuteCommand")
	}

	cmd.SetLimitRespDTO = &resp

	return nil
}

func (h *SetLimit) ConvertCommandToText(cmd *textrouter.Command) (string, error) {
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
