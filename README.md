# telegram-bot myasnikov.alexander.s@gmail.com

### Следано
* Архитектура близка к Clean Architecture
* unit тесты

### TODO
* Увеличить покрытие
* Вынести все ошибки в отдельные файлы
* Добавить логер
* Попробовать добавить интеграционные тесты
* убрать явные зависимости от времени

### Проблемы
* Все хендлеры находятся в одном пакете, из-за этого пересекаются названия интерфейсов
* textrouter.Command имеет слишком много полей, неудобно инифиализировать в тестах
* непонятно со стилем именования пакетов, директорий, интерфейсов
* как и нужно ли тестировать неэкспортируемые функции, из-за этого тесты могут получаться большого размеры

### Структура:

```
├── cmd/bot
│   └── main.go
└── internal
    ├── adapter
    │   ├── service
    │   │   ├── ratesupdaterservicecbr
    │   │   │   ├── rates_updater_service.go
    │   │   │   └── rates_updater_service_test.go
    │   │   └── ratesupdaterserviceexchangerate
    │   │       ├── rates_updater_service.go
    │   │       └── rates_updater_service_test.go
    │   └── storage
    │       ├── currencymemorystorage
    │       │   ├── currency_storage.go
    │       │   └── currency_storage_test.go
    │       ├── expensememorystorage
    │       │   ├── expense_storage.go
    │       │   └── expense_storage_test.go
    │       └── usermemorystorage
    │           ├── user_storage.go
    │           └── user_storage_test.go
    ├── app
    │   └── app.go
    ├── clients/tg
    │   ├── tgclient.go
    │   └── tgclient_test.go
    ├── config
    │   └── config.go
    ├── entity
    │   ├── decimal.go
    │   ├── expense.go
    │   ├── rate.go
    │   ├── time.go
    │   └── user.go
    ├── textrouter
    │   ├── errors.go
    │   ├── router.go
    │   └── texthandler
    │       ├── about.go
    │       ├── add_expense.go
    │       ├── add_expense_test.go
    │       ├── get_report.go
    │       ├── get_report_test.go
    │       ├── help.go
    │       ├── mock_texthandler
    │       │   ├── add_expense.go
    │       │   ├── get_report.go
    │       │   └── set_default_currency.go
    │       ├── set_default_currency.go
    │       ├── set_default_currency_test.go
    │       ├── start.go
    │       └── unknown.go
    ├── usecase
    │   ├── dto.go
    │   ├── expense.go
    │   ├── expense_test.go
    │   └── mock_usecase
    │       └── expense.go
    ├── util
    │   └── date.go
    └── worker/rate_updater_worker
        └── rate_updater_worker.go
```
