package discount

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/farhaan/kuncie/entity"
	"github.com/farhaan/kuncie/utility"
)

type service struct {
	repository Repository
}
type Service interface {
	Apply(ctx context.Context, productCartMap map[entity.SKU]entity.Product) (map[entity.SKU]entity.Product, []entity.Discount)
}

func NewService(repository Repository) Service {
	return &service{repository}
}

// Deduct deduct product cart according to discount that match their criteria
func (s *service) Apply(ctx context.Context, productCartMap map[entity.SKU]entity.Product) (map[entity.SKU]entity.Product, []entity.Discount) {
	var ok bool
	discounts := s.SelectDiscountByMatchingKeyword(ctx, productCartMap)
	appliedDiscount := []entity.Discount{}
	discountedProductCartMap := productCartMap

	// we may improve this deduct function later buy using recursive for accumulative discount
	for _, discount := range discounts {
		switch discount.Type {
		case entity.DiscountTypeBuyXGetY:
			discountedProductCartMap, ok = s.DiscountBuyXGetY(ctx, discount, discountedProductCartMap)
			if ok {
				appliedDiscount = append(appliedDiscount, discount)
			}
		case entity.DiscountTypeBuyN:
			discountedProductCartMap, ok = s.DiscountBuyN(ctx, discount, discountedProductCartMap)
			if ok {
				appliedDiscount = append(appliedDiscount, discount)
			}
		}
	}
	return discountedProductCartMap, appliedDiscount
}

// SelectDiscountByMatchingKeyword select correct discount that matched discount keyword and scanned product SKU
// this operation supposed to be executed in database side using approriate query
func (s *service) SelectDiscountByMatchingKeyword(ctx context.Context, productCartMap map[entity.SKU]entity.Product) []entity.Discount {
	discounts := s.repository.GetActiveDiscountsForCheckout(ctx)
	var selectedDiscounts []entity.Discount
	for _, discount := range discounts {
		var matchCounter int
		for _, product := range productCartMap {
			if utility.StringInArray(string(product.SKU), discount.Keywords) {
				matchCounter++
			}
		}
		if matchCounter == len(discount.Keywords) {
			selectedDiscounts = append(selectedDiscounts, discount)
		}
	}

	return selectedDiscounts
}

// DiscountBuyXGetY execute discount logic for "buy x item to get y item"
// for each x item that matched the criteria,
// change y item price to 0 if y item exist, else append y item into product cart map
func (s *service) DiscountBuyXGetY(ctx context.Context, discount entity.Discount, productCartMap map[entity.SKU]entity.Product) (discountedProductCartMap map[entity.SKU]entity.Product, ok bool) {
	discountedProductCartMap = productCartMap
	for productSKU, itemValue := range discount.Rule.ItemValue {
		product, ok := discountedProductCartMap[entity.SKU(productSKU)]
		if !ok && itemValue.Min > 0 {
			return productCartMap, false
		}
		if !ok {
			discountedProductCartMap[entity.SKU(productSKU)] = entity.Product{SKU: product.SKU, BuyQty: 1}
			continue
		}
		if itemValue.Min == 0 {
			if product.BuyQty > 0 {
				tempSKUForDiscount := entity.SKU(fmt.Sprintf("%s-%d", productSKU, time.Now().UnixNano()*rand.Int63()))
				discountedProductCartMap[tempSKUForDiscount] = entity.Product{
					SKU:    product.SKU,
					Name:   "DISCOUNT AMOUNT",
					Price:  -product.Price,
					BuyQty: 1,
				}
			}
		}

		discountedProductCartMap[entity.SKU(productSKU)] = product
	}
	return discountedProductCartMap, true
}

// DiscountBuyN execute discount logic for "buy n item for x amount to get discount"
func (s *service) DiscountBuyN(ctx context.Context, discount entity.Discount, productCartMap map[entity.SKU]entity.Product) (discountedProductCartMap map[entity.SKU]entity.Product, ok bool) {
	discountedProductCartMap = productCartMap
	for productSKU, itemValue := range discount.Rule.ItemValue {
		sku := entity.SKU(productSKU)
		product, ok := discountedProductCartMap[sku]
		if !ok {
			return productCartMap, false
		}
		if itemValue.Min > product.BuyQty {
			return productCartMap, false
		}

		// product.BuyQty -= itemValue.Min
		tempSKUForDiscount := entity.SKU(fmt.Sprintf("%s-%d", productSKU, time.Now().UnixNano()*rand.Int63()))
		if discount.AmountType == entity.DiscountAmountTypeAmount {
			discountedProductCartMap[tempSKUForDiscount] = entity.Product{
				SKU:    tempSKUForDiscount,
				Name:   "DISCOUNT AMOUNT",
				Price:  -discount.Amount,
				BuyQty: 1,
			}
		}
		if discount.AmountType == entity.DiscountAmountTypePercentage {
			discountedProductCartMap[tempSKUForDiscount] = entity.Product{
				SKU:    tempSKUForDiscount,
				Name:   "DISCOUNT AMOUNT",
				Price:  -((100 - discount.Amount) * product.Price / 100),
				BuyQty: product.BuyQty,
			}
		}

		discountedProductCartMap[entity.SKU(productSKU)] = product
	}
	return discountedProductCartMap, true
}
