package router

import "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"

const (
	cmdUnknown = iota + 1
	cmdStart
	cmdHelp
	cmdAbout
	cmdAddProduct
	cmdGetStats

	cmdUnknownText = "" +
		"Не могу понять.\n" +
		"Введи /help для более подробной информации"
	cmdStartText = "" +
		"Привет. Напиши освои расходы и я запомпю их.\n" +
		"Введи /help для более подробной информации"
	cmdHelpText = "" +
		"Я понимаю следующие команды:\n" +
		"/start                          - приветственное сообщение\n" +
		"/help                           - стравочная информация\n" +
		"/about                          - информация о проекте\n" +
		"сохранить <категория> <суммa>   - добавление расходов\n" +
		"отчет <день|неделя|месяц|год>   - отчет за интервал\n"
	cmdAboutText = "" +
		"Я бот для учета расходов." +
		"Автор @amyasnikov"
	cmdAddProductText = "Записал"
)

type command struct {
	id                   int
	addProductReqDTO     *usecase.AddProductReqDTO
	getStatisticsReqDTO  *usecase.GetStatisticsReqDTO
	getStatisticsRespDTO *usecase.GetStatisticsRespDTO
}
