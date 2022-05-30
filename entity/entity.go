package entity

import "time"

type Product struct {
	SKU      SKU     `json:"sku"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	StockQty int     `json:"stockQty"`
	BuyQty   int     `json:"buyQty"`
}

type SKU string

type DiscountType int
type DiscountAmountType int

const (
	// TypeBuyXGetY means buy item X to get free Y item
	DiscountTypeBuyXGetY DiscountType = iota + 1
	// Type BuyN means buy N amount of item to get price reduction
	DiscountTypeBuyN
)

const (
	DiscountAmountTypePercentage DiscountAmountType = iota + 1
	DiscountAmountTypeAmount
	DiscountAmountTypePriceReductionIfExist
)

type DiscountRule struct {
	StartDate    time.Time                           `json:"startDate"`
	EndDate      time.Time                           `json:"endDate"`
	Accumulative bool                                `json:"accumulative"`
	ItemValue    map[string]DiscountRuleItemCriteria `json:"itemValue"`
}

type DiscountRuleItemCriteria struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type Discount struct {
	Type       DiscountType       `json:"type"`
	Keywords   []string           `json:"keywords"`
	Rule       DiscountRule       `json:"rule"`
	AmountType DiscountAmountType `json:"amountType"`
	Amount     float64            `json:"amount"`
}

type Checkout struct {
	Products        map[SKU]Product `json:"products"`
	FinalProducts   map[SKU]Product `json:"finalProducts"`
	Discounts       []Discount      `json:"discounts"`
	TotalBasePrice  float64         `json:"totalBasePrice"`
	TotalFinalPrice float64         `json:"finalPrice"`
}
