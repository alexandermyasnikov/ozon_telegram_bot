package router

import (
	"fmt"
	"regexp"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

var _ HandlerTextInterface = (*HandlerTextGetStats)(nil)

var reDate = regexp.MustCompile(genRE(
	[]string{"cmd", `отчет`, `стат`},
	[]string{"date", `ден`, `нед`, `мес`, `год`},
))

type HandlerTextGetStats struct{}

func (h *HandlerTextGetStats) GetID() int {
	return cmdGetStats
}

func (h *HandlerTextGetStats) ConvertTextToCommand(userID int64, text string, date time.Time, cmd *command) {
	if matches := reDate.FindStringSubmatch(text); len(matches) > 0 {
		days := dateToDays(matches[reDate.SubexpIndex("date")])

		cmd.id = cmdGetStats
		cmd.getStatisticsReqDTO = &usecase.GetStatisticsReqDTO{
			UserID: userID,
			Date:   date,
			Days:   days,
		}
	}
}

func (h *HandlerTextGetStats) ExecuteCommand(cmd *command, productUsecase usecase.ProductUsecaseInterface) error {
	if cmd.getStatisticsReqDTO == nil {
		return errors.New("empty getStatisticsReqDTO")
	}

	resp, err := productUsecase.GetStatistics(*cmd.getStatisticsReqDTO)
	*cmd.getStatisticsRespDTO = resp

	return errors.Wrap(err, "productUsecase.GetStatistics")
}

func (h *HandlerTextGetStats) ConvertCommandToText(cmd command) string {
	if cmd.getStatisticsRespDTO == nil {
		return "empty getStatisticsRespDTO"
	}

	textOut := "Расходы:\n"
	for k, v := range cmd.getStatisticsRespDTO.Products {
		textOut += fmt.Sprintf("%s - %d\n", k, v)
	}

	return textOut
}
