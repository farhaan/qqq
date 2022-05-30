package discount

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/farhaan/kuncie/entity"
)

type repository struct {
	discounts []entity.Discount
}
type Repository interface {
	GetActiveDiscountsForCheckout(ctx context.Context) (discounts []entity.Discount)
}

func NewRepository() Repository {
	return &repository{discounts: initiate()}
}

func (r *repository) GetActiveDiscountsForCheckout(ctx context.Context) (discounts []entity.Discount) {
	var activeDiscount []entity.Discount
	now := time.Now()
	for _, discount := range r.discounts {
		if discount.Rule.StartDate.Before(now) && discount.Rule.EndDate.After(now) {
			activeDiscount = append(activeDiscount, discount)
		}
	}

	return activeDiscount
}

// initiate initiate discount data using fake json
func initiate() []entity.Discount {
	jsonFile, err := os.Open("database/discounts.json")
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	var discounts []entity.Discount

	jsonDecoder := json.NewDecoder(jsonFile)
	jsonDecoder.UseNumber()

	if err := jsonDecoder.Decode(&discounts); err != nil {
		panic(err)
	}

	return discounts
}
