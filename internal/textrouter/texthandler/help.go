package texthandler

import (
	"context"
	"strings"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

type Help struct{}

func NewHelp() *Help {
	return &Help{}
}

func (h *Help) Name() string {
	return usecase.HelpCmdName
}

func (h *Help) ConvertTextToCommand(ctx context.Context, text string, cmd *usecase.Command) bool {
	return strings.HasPrefix(text, "/help")
}

func (h *Help) ConvertCommandToText(ctx context.Context, cmd *usecase.Command) (string, error) {
	return `Я понимаю следующие команды:
/start                               - приветственное сообщение
/help                                - стравочная информация
/about                               - информация о проекте
валюта <валюта>                      - выбрать валюту по умолчанию
расход <категория> <суммa> <валюта>  - добавление расходов
отчет <период>                       - отчет за интервал
лимит <период> <сумма>               - установить бюджет`, nil
}
