package texthandler

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/utils"
)

type GetLimits struct{}

func NewGetLimits() *GetLimits {
	return &GetLimits{}
}

func (h *GetLimits) Name() string {
	return usecase.GetLimitsCmdName
}

func (h *GetLimits) ConvertTextToCommand(ctx context.Context, text string, cmd *usecase.Command) bool {
	fields := strings.Fields(text)
	if len(fields) != 1 || fields[0] != "лимиты" {
		return false
	}

	cmd.GetLimitsReqDTO = &usecase.GetLimitsReqDTO{
		UserID: cmd.UserID,
	}

	return true
}

func (h *GetLimits) ConvertCommandToText(ctx context.Context, cmd *usecase.Command) (string, error) {
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
