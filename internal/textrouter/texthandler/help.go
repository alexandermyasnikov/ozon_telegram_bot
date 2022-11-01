package texthandler

import (
	"context"
	"strings"
	"time"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
)

type Help struct{}

func NewHelp() *Help {
	return &Help{}
}

func (h *Help) Name() string {
	return "help"
}

func (h *Help) ConvertTextToCommand(userID int64, text string, date time.Time, cmd *textrouter.Command) bool {
	return strings.HasPrefix(text, "/help")
}

func (h *Help) ExecuteCommand(ctx context.Context, cmd *textrouter.Command) error {
	return nil
}

func (h *Help) ConvertCommandToText(cmd *textrouter.Command) (string, error) {
	return `Я понимаю следующие команды:
/start                               - приветственное сообщение
/help                                - стравочная информация
/about                               - информация о проекте
валюта <валюта>                      - выбрать валюту по умолчанию
расход <категория> <суммa> <валюта>  - добавление расходов
отчет <период>                       - отчет за интервал
лимит <период> <сумма>               - установить бюджет`, nil
}
