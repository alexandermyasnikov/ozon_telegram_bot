package usecase

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
)

type ICurrencyStorage interface {
	Get(string) (entity.Rate, error)
	GetAll() ([]entity.Rate, error)
	Update(entity.Rate) error
}

type IUserStorage interface {
	GetDefaultCurrency(entity.UserID) (string, error)
	UpdateDefaultCurrency(entity.UserID, string) error
}

type IExpenseStorage interface {
	Create(entity.UserID, entity.Expense) error
	Get(entity.UserID, entity.Date, int) ([]entity.Expense, error)
}

type IRatesUpdaterService interface {
	Get(ctx context.Context, base string, codes []string) ([]entity.Rate, error)
}

type IConfig interface {
	GetBaseCurrencyCode() string
	GetCurrencyCodes() []string
	GetFrequencyRateUpdateSec() int
}

type ExpenseUsecase struct {
	currencyStorage     ICurrencyStorage
	userStorage         IUserStorage
	expenseStorage      IExpenseStorage
	ratesUpdaterService IRatesUpdaterService
	config              IConfig
}

func NewExpenseUsecase(currencyStorage ICurrencyStorage, userStorage IUserStorage, expenseStorage IExpenseStorage,
	ratesUpdaterService IRatesUpdaterService, config IConfig,
) *ExpenseUsecase {
	return &ExpenseUsecase{
		currencyStorage:     currencyStorage,
		userStorage:         userStorage,
		expenseStorage:      expenseStorage,
		ratesUpdaterService: ratesUpdaterService,
		config:              config,
	}
}

func (uc *ExpenseUsecase) SetDefaultCurrency(ctx context.Context, req SetDefaultCurrencyReqDTO) error {
	userID := entity.NewUserID(req.UserID)

	ok := uc.isSupportedCurrencyCode(req.Currency)
	if !ok {
		return errors.New("currency is unsupported")
	}

	err := uc.userStorage.UpdateDefaultCurrency(userID, req.Currency)

	return errors.Wrap(err, "ExpenseUsecase.SetDefaultCurrency")
}

func (uc *ExpenseUsecase) AddExpense(ctx context.Context, req AddExpenseReqDTO) error {
	userID := entity.NewUserID(req.UserID)

	err := uc.tryUpdateRates(ctx, false)
	if err != nil {
		return errors.Wrap(err, "ExpenseUsecase.AddExpense")
	}

	currency, err := uc.getCurrencyForUser(userID, req.Currency)
	if err != nil {
		return errors.Wrap(err, "ExpenseUsecase.AddExpense")
	}

	rate, err := uc.currencyStorage.Get(currency)
	if err != nil {
		return errors.Wrap(err, "ExpenseUsecase.AddExpense")
	}

	price := entity.NewDecimalFromFloat(req.Price)
	price.Div(rate.GetRatio())

	date := entity.NewDateFromTime(req.Date)
	expense := entity.NewExpense(req.Category, price, date)

	err = uc.expenseStorage.Create(userID, expense)

	return errors.Wrap(err, "ExpenseUsecase.AddExpense")
}

func (uc *ExpenseUsecase) GetReport(ctx context.Context, req GetReportReqDTO) (GetReportRespDTO, error) {
	userID := entity.NewUserID(req.UserID)
	date := entity.NewDateFromTime(req.Date)

	err := uc.tryUpdateRates(ctx, false)
	if err != nil {
		return GetReportRespDTO{}, errors.Wrap(err, "ExpenseUsecase.AddExpense")
	}

	expenses, err := uc.expenseStorage.Get(userID, date, req.Days)
	if err != nil {
		return GetReportRespDTO{}, errors.Wrap(err, "ExpenseUsecase.GetReport")
	}

	currency, err := uc.getCurrencyForUser(userID, req.Currency)
	if err != nil {
		return GetReportRespDTO{}, errors.Wrap(err, "ExpenseUsecase.AddExpense")
	}

	rate, err := uc.currencyStorage.Get(currency)
	if err != nil {
		return GetReportRespDTO{}, errors.Wrap(err, "ExpenseUsecase.GetReport")
	}

	categories := make(map[string]entity.Decimal)

	for _, expense := range expenses {
		price := expense.GetPrice()
		price.Mult(rate.GetRatio())
		price.Add(categories[expense.GetCategory()])

		categories[expense.GetCategory()] = price
	}

	resp := GetReportRespDTO{
		Currency:   currency,
		Categories: categories, // TODO sorted slice
	}

	return resp, nil
}

// ----

// TODO PresentorUsecase ?
func (uc *ExpenseUsecase) GetAllCurrencyNames(ctx context.Context) (GetAllCurrencyNamesRespDTO, error) {
	currencyRates, err := uc.currencyStorage.GetAll()

	names := make([]string, 0, len(currencyRates))
	for _, currencyRate := range currencyRates {
		names = append(names, currencyRate.GetCode())
	}

	resp := GetAllCurrencyNamesRespDTO{
		Currencies: names,
	}

	return resp, errors.Wrap(err, "ExpenseUsecase.GetAllCurrencyNames")
}

// TODO WorkerUsecase ?
func (uc *ExpenseUsecase) UpdateCurrency(ctx context.Context) error {
	err := uc.tryUpdateRates(ctx, true)

	return errors.Wrap(err, "ExpenseUsecase.UpdateCurrency")
}

// ---- helpers

func (uc *ExpenseUsecase) tryUpdateRates(ctx context.Context, force bool) error {
	if !force && !uc.needUpdateRates() {
		return nil
	}

	rates, err := uc.ratesUpdaterService.Get(ctx, uc.config.GetBaseCurrencyCode(), uc.config.GetCurrencyCodes())
	if err != nil {
		return errors.Wrap(err, "ExpenseUsecase.tryUpdateRates")
	}

	for _, rate := range rates {
		err := uc.currencyStorage.Update(rate)
		if err != nil {
			return errors.Wrap(err, "ExpenseUsecase.tryUpdateRates")
		}
	}

	return nil
}

func (uc *ExpenseUsecase) needUpdateRates() bool {
	rate, err := uc.currencyStorage.Get(uc.config.GetBaseCurrencyCode())
	if err != nil {
		return true
	}

	duration := time.Since(rate.GetTime().ToTime()).Seconds()

	return duration > float64(uc.config.GetFrequencyRateUpdateSec())
}

func (uc *ExpenseUsecase) getCurrencyForUser(userID entity.UserID, currency string) (string, error) {
	if len(currency) != 0 {
		// Используется валюта указанная в операции
		if uc.isSupportedCurrencyCode(currency) {
			return currency, nil
		}

		return "", errors.New("currency is unsupported")
	}

	currency, err := uc.userStorage.GetDefaultCurrency(userID)
	if err == nil {
		// Используется валюта дефолтная для пользователя
		return currency, nil
	}

	// Используется валюта дефолтная для всех
	return uc.config.GetBaseCurrencyCode(), nil
}

func (uc *ExpenseUsecase) isSupportedCurrencyCode(currency string) bool {
	if currency == uc.config.GetBaseCurrencyCode() {
		return true
	}

	for _, supportedCode := range uc.config.GetCurrencyCodes() {
		if currency == supportedCode {
			return true
		}
	}

	return false
}
