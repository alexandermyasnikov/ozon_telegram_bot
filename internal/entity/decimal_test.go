package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
)

func TestDecimalToFloat(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description string
		number      int64
		scale       uint64
		floatExp    float64
	}

	testCases := [...]testCase{
		{
			description: "decimal",
			number:      1234567,
			scale:       5,
			floatExp:    12.34567,
		},
		{
			description: "zero decimal",
			number:      0,
			scale:       0,
			floatExp:    0,
		},
		{
			description: "only fractional part",
			number:      12,
			scale:       5,
			floatExp:    0.00012,
		},
		{
			description: "only integer part",
			number:      10000,
			scale:       3,
			floatExp:    10,
		},
		{
			description: "negitive decimal",
			number:      -1234567,
			scale:       5,
			floatExp:    -12.34567,
		},
	}

	for _, scenario := range testCases {
		scenario := scenario
		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

			dec := entity.NewDecimal(scenario.number, scenario.scale)
			assert.Equal(t, scenario.floatExp, dec.ToFloat())
		})
	}
}

func TestDecimalFromFloat(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description string
		float       float64
	}

	testCases := [...]testCase{
		{
			description: "float",
			float:       123456789,
		},
		{
			description: "float",
			float:       10,
		},
		{
			description: "float",
			float:       123.456789,
		},
		{
			description: "float",
			float:       0.012345,
		},
		{
			description: "negative float",
			float:       -0.012345,
		},
	}

	for _, scenario := range testCases {
		scenario := scenario
		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

			dec := entity.NewDecimalFromFloat(scenario.float)
			assert.InDelta(t, scenario.float, dec.ToFloat(), 10e-5)
		})
	}
}

func TestDecimalAdd(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description string
		number1     entity.Decimal
		number2     entity.Decimal
		numberExp   entity.Decimal
	}

	testCases := [...]testCase{
		{
			description: "add",
			number1:     entity.NewDecimal(123456, 3),
			number2:     entity.NewDecimal(55, 1),
			numberExp:   entity.NewDecimal(128956, 3),
		},
		{
			description: "add",
			number1:     entity.NewDecimal(10000, 0),
			number2:     entity.NewDecimal(12, 6),
			numberExp:   entity.NewDecimal(10000000012, 6),
		},
	}

	for _, scenario := range testCases {
		scenario := scenario
		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

			scenario.number1.Add(scenario.number2)
			assert.Equal(t, scenario.numberExp, scenario.number1)
		})
	}
}

func TestDecimalMult(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description string
		number1     entity.Decimal
		number2     entity.Decimal
		numberExp   entity.Decimal
	}

	testCases := [...]testCase{
		{
			description: "mult",
			number1:     entity.NewDecimal(55, 1),
			number2:     entity.NewDecimal(1234567, 3),
			numberExp:   entity.NewDecimal(67901185, 4),
		},
	}

	for _, scenario := range testCases {
		scenario := scenario
		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

			scenario.number1.Mult(scenario.number2)
			assert.Equal(t, scenario.numberExp, scenario.number1)
		})
	}
}

func TestDecimalDiv(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description string
		number1     entity.Decimal
		number2     entity.Decimal
		numberExp   entity.Decimal
	}

	testCases := [...]testCase{
		{
			description: "div",
			number1:     entity.NewDecimal(67901185, 4),
			number2:     entity.NewDecimal(1234567, 3),
			numberExp:   entity.NewDecimal(55, 1),
		},
	}

	for _, scenario := range testCases {
		scenario := scenario
		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

			scenario.number1.Div(scenario.number2)
			assert.Equal(t, scenario.numberExp, scenario.number1)
		})
	}
}
