package productmemorystorage_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	productmemorystorage "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
)

func Test_GetAll(t *testing.T) {
	t.Parallel()

	type testInParam struct {
		userID  int64
		product entity.Product
	}

	type testOutParam struct {
		description string
		userID      int64
		days        int
		date        time.Time
		products    []entity.Product
	}

	type testCase struct {
		description   string
		testInParams  []testInParam
		testOutParams []testOutParam
	}

	testCases := [...]testCase{
		{
			description:  "Empty storage",
			testInParams: []testInParam{},
			testOutParams: []testOutParam{
				{
					description: "0 days",
					userID:      1,
					days:        0,
					date:        time.Date(2022, 9, 20, 12, 0, 0, 0, time.UTC),
					products:    []entity.Product{},
				},
				{
					description: "365 days",
					userID:      1,
					days:        365,
					date:        time.Date(2022, 9, 20, 12, 0, 0, 0, time.UTC),
					products:    []entity.Product{},
				},
			},
		},
		{
			description: "Single user, same-day products",
			testInParams: []testInParam{
				{
					userID:  1,
					product: entity.NewProduct("name1", "cat1", 1000, time.Date(2022, 9, 20, 12, 0, 0, 0, time.UTC)),
				},
				{
					userID:  1,
					product: entity.NewProduct("name2", "cat2", 2000, time.Date(2022, 9, 20, 12, 10, 0, 0, time.UTC)),
				},
				{
					userID:  1,
					product: entity.NewProduct("name3", "cat3", 3000, time.Date(2022, 9, 20, 10, 0, 0, 0, time.UTC)),
				},
				{
					userID:  1,
					product: entity.NewProduct("name2", "cat2", 2000, time.Date(2022, 9, 20, 14, 30, 0, 0, time.UTC)),
				},
			},
			testOutParams: []testOutParam{
				{
					description: "0 days, 20 Sep",
					userID:      1,
					days:        0,
					date:        time.Date(2022, 9, 20, 12, 0, 0, 0, time.UTC),
					products:    []entity.Product{},
				},
				{
					description: "1 days, 19 Sep",
					userID:      1,
					days:        1,
					date:        time.Date(2022, 9, 19, 12, 0, 0, 0, time.UTC),
					products:    []entity.Product{},
				},
				{
					description: "1 days, 21 Sep",
					userID:      1,
					days:        1,
					date:        time.Date(2022, 9, 21, 12, 0, 0, 0, time.UTC),
					products:    []entity.Product{},
				},
				{
					description: "1 days, 20 Sep",
					userID:      1,
					days:        1,
					date:        time.Date(2022, 9, 20, 12, 0, 0, 0, time.UTC),
					products: []entity.Product{
						entity.NewProduct("name1", "cat1", 1000, time.Date(2022, 9, 20, 12, 0, 0, 0, time.UTC)),
						entity.NewProduct("name2", "cat2", 2000, time.Date(2022, 9, 20, 12, 10, 0, 0, time.UTC)),
						entity.NewProduct("name3", "cat3", 3000, time.Date(2022, 9, 20, 10, 0, 0, 0, time.UTC)),
						entity.NewProduct("name2", "cat2", 2000, time.Date(2022, 9, 20, 14, 30, 0, 0, time.UTC)),
					},
				},
			},
		},
		{
			description: "Single user, different days products, unordered",
			testInParams: []testInParam{
				{
					userID:  1,
					product: entity.NewProduct("name1", "cat1", 2000, time.Date(2022, 9, 20, 12, 0, 0, 0, time.UTC)),
				},
				{
					userID:  1,
					product: entity.NewProduct("name2", "cat2", 3000, time.Date(2022, 9, 17, 14, 0, 0, 0, time.UTC)),
				},
				{
					userID:  1,
					product: entity.NewProduct("name3", "cat3", 3000, time.Date(2022, 9, 17, 12, 0, 0, 0, time.UTC)),
				},
				{
					userID:  1,
					product: entity.NewProduct("name4", "cat4", 3000, time.Date(2022, 9, 17, 13, 0, 0, 0, time.UTC)),
				},
				{
					userID:  1,
					product: entity.NewProduct("name5", "cat5", 2000, time.Date(2022, 9, 22, 12, 0, 0, 0, time.UTC)),
				},
			},
			testOutParams: []testOutParam{
				{
					description: "1 days, 17 Sep",
					userID:      1,
					days:        1,
					date:        time.Date(2022, 9, 17, 0, 0, 0, 0, time.UTC),
					products: []entity.Product{
						entity.NewProduct("name2", "cat2", 3000, time.Date(2022, 9, 17, 14, 0, 0, 0, time.UTC)),
						entity.NewProduct("name3", "cat3", 3000, time.Date(2022, 9, 17, 12, 0, 0, 0, time.UTC)),
						entity.NewProduct("name4", "cat4", 3000, time.Date(2022, 9, 17, 13, 0, 0, 0, time.UTC)),
					},
				},
				{
					description: "1 days, 20 Sep",
					userID:      1,
					days:        1,
					date:        time.Date(2022, 9, 20, 0, 0, 0, 0, time.UTC),
					products: []entity.Product{
						entity.NewProduct("name1", "cat1", 2000, time.Date(2022, 9, 20, 12, 0, 0, 0, time.UTC)),
					},
				},
				{
					description: "1 days, 22 Sep",
					userID:      1,
					days:        1,
					date:        time.Date(2022, 9, 22, 0, 0, 0, 0, time.UTC),
					products: []entity.Product{
						entity.NewProduct("name5", "cat5", 2000, time.Date(2022, 9, 22, 12, 0, 0, 0, time.UTC)),
					},
				},
				{
					description: "2 days, 22 Sep",
					userID:      1,
					days:        2,
					date:        time.Date(2022, 9, 22, 0, 0, 0, 0, time.UTC),
					products: []entity.Product{
						entity.NewProduct("name5", "cat5", 2000, time.Date(2022, 9, 22, 12, 0, 0, 0, time.UTC)),
					},
				},
				{
					description: "3 days, 22 Sep",
					userID:      1,
					days:        3,
					date:        time.Date(2022, 9, 22, 0, 0, 0, 0, time.UTC),
					products: []entity.Product{
						entity.NewProduct("name5", "cat5", 2000, time.Date(2022, 9, 22, 12, 0, 0, 0, time.UTC)),
						entity.NewProduct("name1", "cat1", 2000, time.Date(2022, 9, 20, 12, 0, 0, 0, time.UTC)),
					},
				},
				{
					description: "7 days, 22 Sep",
					userID:      1,
					days:        7,
					date:        time.Date(2022, 9, 22, 0, 0, 0, 0, time.UTC),
					products: []entity.Product{
						entity.NewProduct("name5", "cat5", 2000, time.Date(2022, 9, 22, 12, 0, 0, 0, time.UTC)),
						entity.NewProduct("name1", "cat1", 2000, time.Date(2022, 9, 20, 12, 0, 0, 0, time.UTC)),
						entity.NewProduct("name2", "cat2", 3000, time.Date(2022, 9, 17, 14, 0, 0, 0, time.UTC)),
						entity.NewProduct("name3", "cat3", 3000, time.Date(2022, 9, 17, 12, 0, 0, 0, time.UTC)),
						entity.NewProduct("name4", "cat4", 3000, time.Date(2022, 9, 17, 13, 0, 0, 0, time.UTC)),
					},
				},
				{
					description: "7 days, 21 Sep",
					userID:      1,
					days:        7,
					date:        time.Date(2022, 9, 21, 0, 0, 0, 0, time.UTC),
					products: []entity.Product{
						entity.NewProduct("name1", "cat1", 2000, time.Date(2022, 9, 20, 12, 0, 0, 0, time.UTC)),
						entity.NewProduct("name2", "cat2", 3000, time.Date(2022, 9, 17, 14, 0, 0, 0, time.UTC)),
						entity.NewProduct("name3", "cat3", 3000, time.Date(2022, 9, 17, 12, 0, 0, 0, time.UTC)),
						entity.NewProduct("name4", "cat4", 3000, time.Date(2022, 9, 17, 13, 0, 0, 0, time.UTC)),
					},
				},
				{
					description: "7 days, 17 Sep",
					userID:      1,
					days:        7,
					date:        time.Date(2022, 9, 17, 0, 0, 0, 0, time.UTC),
					products: []entity.Product{
						entity.NewProduct("name2", "cat2", 3000, time.Date(2022, 9, 17, 14, 0, 0, 0, time.UTC)),
						entity.NewProduct("name3", "cat3", 3000, time.Date(2022, 9, 17, 12, 0, 0, 0, time.UTC)),
						entity.NewProduct("name4", "cat4", 3000, time.Date(2022, 9, 17, 13, 0, 0, 0, time.UTC)),
					},
				},
			},
		},
		{
			description: "Many users",
			testInParams: []testInParam{
				{
					userID:  1,
					product: entity.NewProduct("name1", "cat1", 2000, time.Date(2022, 9, 20, 12, 0, 0, 0, time.UTC)),
				},
				{
					userID:  2,
					product: entity.NewProduct("name2", "cat2", 2000, time.Date(2022, 9, 20, 12, 0, 0, 0, time.UTC)),
				},
				{
					userID:  3,
					product: entity.NewProduct("name3", "cat3", 2000, time.Date(2022, 9, 20, 12, 0, 0, 0, time.UTC)),
				},
				{
					userID:  1,
					product: entity.NewProduct("name4", "cat4", 2000, time.Date(2022, 9, 21, 14, 0, 0, 0, time.UTC)),
				},
				{
					userID:  2,
					product: entity.NewProduct("name5", "cat5", 2000, time.Date(2022, 9, 22, 14, 0, 0, 0, time.UTC)),
				},
				{
					userID:  3,
					product: entity.NewProduct("name6", "cat6", 2000, time.Date(2022, 9, 23, 14, 0, 0, 0, time.UTC)),
				},
				{
					userID:  2,
					product: entity.NewProduct("name7", "cat7", 2000, time.Date(2022, 9, 25, 12, 0, 0, 0, time.UTC)),
				},
			},
			testOutParams: []testOutParam{
				{
					description: "1 uesrId, 7 days, 25 Sep",
					userID:      1,
					days:        7,
					date:        time.Date(2022, 9, 25, 0, 0, 0, 0, time.UTC),
					products: []entity.Product{
						entity.NewProduct("name4", "cat4", 2000, time.Date(2022, 9, 21, 14, 0, 0, 0, time.UTC)),
						entity.NewProduct("name1", "cat1", 2000, time.Date(2022, 9, 20, 12, 0, 0, 0, time.UTC)),
					},
				},
				{
					description: "3 uesrId, 7 days, 25 Sep",
					userID:      3,
					days:        7,
					date:        time.Date(2022, 9, 25, 0, 0, 0, 0, time.UTC),
					products: []entity.Product{
						entity.NewProduct("name6", "cat6", 2000, time.Date(2022, 9, 23, 14, 0, 0, 0, time.UTC)),
						entity.NewProduct("name3", "cat3", 2000, time.Date(2022, 9, 20, 12, 0, 0, 0, time.UTC)),
					},
				},
				{
					description: "2 uesrId, 5 days, 25 Sep",
					userID:      2,
					days:        5,
					date:        time.Date(2022, 9, 25, 0, 0, 0, 0, time.UTC),
					products: []entity.Product{
						entity.NewProduct("name7", "cat7", 2000, time.Date(2022, 9, 25, 12, 0, 0, 0, time.UTC)),
						entity.NewProduct("name5", "cat5", 2000, time.Date(2022, 9, 22, 14, 0, 0, 0, time.UTC)),
					},
				},
				{
					description: "2 uesrId, 2 days, 25 Sep",
					userID:      2,
					days:        2,
					date:        time.Date(2022, 9, 25, 0, 0, 0, 0, time.UTC),
					products: []entity.Product{
						entity.NewProduct("name7", "cat7", 2000, time.Date(2022, 9, 25, 12, 0, 0, 0, time.UTC)),
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()
			storage := productmemorystorage.NewProductStorage()
			for _, param := range testCase.testInParams {
				err := storage.Add(param.userID, param.product)
				assert.NoError(t, err)
			}

			for _, param := range testCase.testOutParams {
				param := param
				t.Run(param.description, func(t *testing.T) {
					t.Parallel()
					actualProducts, err := storage.GetAll(param.userID, param.date, param.days)
					assert.NoError(t, err)

					assert.Equal(t, len(param.products), len(actualProducts))
					assert.Equal(t, param.products, actualProducts)
				})
			}
		})
	}
}
