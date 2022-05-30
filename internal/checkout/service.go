package checkout

import (
	"context"

	"github.com/farhaan/kuncie/entity"
	"github.com/farhaan/kuncie/internal/discount"
	"github.com/farhaan/kuncie/internal/product"
)

type service struct {
	productSvc  product.Service
	discountSvc discount.Service
}
type Service interface {
	Checkout(ctx context.Context, productSKUs []entity.SKU) (checkout *entity.Checkout)
}

func NewService(productSvc product.Service, discountSvc discount.Service) Service {
	return &service{productSvc, discountSvc}
}

func (s *service) Checkout(ctx context.Context, productSKUs []entity.SKU) (checkout *entity.Checkout) {
	products := s.productSvc.GetProductsForCheckoutBySKUs(ctx, productSKUs)
	productCartMap := s.productSvc.AssignProductsToMap(products)
	discountedProductCartMap, discounts := s.discountSvc.Apply(ctx, productCartMap)

	checkout = &entity.Checkout{
		Products:      productCartMap,
		FinalProducts: discountedProductCartMap,
		Discounts:     discounts,
	}

	checkout.TotalBasePrice, checkout.TotalFinalPrice = s.Deduct(productCartMap, discountedProductCartMap)

	return checkout
}

func (s *service) Deduct(productCartMap, discountedProductCartMap map[entity.SKU]entity.Product) (totalBasePrice, totalFinalPrice float64) {
	for _, productCart := range discountedProductCartMap {
		totalFinalPrice += productCart.Price * float64(productCart.BuyQty)
	}
	for _, productCart := range productCartMap {
		totalBasePrice += productCart.Price * float64(productCart.BuyQty)
	}

	return totalBasePrice, totalFinalPrice
}
