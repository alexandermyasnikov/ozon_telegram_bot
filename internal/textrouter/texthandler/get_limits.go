package texthandler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/utils"
)

type ExpenseUsecaseGL interface {
	GetLimits(ctx context.Context, req usecase.GetLimitsReqDTO) (usecase.GetLimitsRespDTO, error)
}

type GetLimits struct {
	expenseUsecase ExpenseUsecaseGL
}

func NewGetLimits(expenseUsecase ExpenseUsecaseGL) *GetLimits {
	return &GetLimits{
		expenseUsecase: expenseUsecase,
	}
}

func (h *GetLimits) ConvertTextToCommand(userID int64, text string, date time.Time, cmd *textrouter.Command) bool {
	fields := strings.Fields(text)
	if len(fields) != 1 || fields[0] != "лимиты" {
		return false
	}

	cmd.GetLimitsReqDTO = &usecase.GetLimitsReqDTO{
		UserID: userID,
	}

	return true
}

func (h *GetLimits) ExecuteCommand(ctx context.Context, cmd *textrouter.Command) error {
	if cmd.GetLimitsReqDTO == nil {
		return errors.Wrap(textrouter.ErrInvalidCommand, "GetLimits.ExecuteCommand")
	}

	resp, err := h.expenseUsecase.GetLimits(ctx, *cmd.GetLimitsReqDTO)
	if err != nil {
		return errors.Wrap(err, "GetLimits.ExecuteCommand")
	}

	cmd.GetLimitsRespDTO = &resp

	return nil
}

func (h *GetLimits) ConvertCommandToText(cmd *textrouter.Command) (string, error) {
	if cmd.GetLimitsReqDTO == nil || cmd.GetLimitsRespDTO == nil {
		return "", errors.Wrap(textrouter.ErrInvalidCommand, "GetLimits.ExecuteCommand")
	}

	precision := 2

	textOut := fmt.Sprintf(`Текущие лимиты:
Дневной - %s %s
Недельный - %s %s
Месячный - %0s %s`,
		cmd.GetLimitsRespDTO.Limits[utils.DayInterval].StringFixed(int32(precision)), cmd.GetLimitsRespDTO.Currency,
		cmd.GetLimitsRespDTO.Limits[utils.WeekInterval].StringFixed(int32(precision)), cmd.GetLimitsRespDTO.Currency,
		cmd.GetLimitsRespDTO.Limits[utils.MonthInterval].StringFixed(int32(precision)), cmd.GetLimitsRespDTO.Currency)

	return textOut, nil
}
