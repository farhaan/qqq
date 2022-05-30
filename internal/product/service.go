package product

import (
	"context"

	"github.com/farhaan/kuncie/entity"
)

type service struct {
	repository Repository
}
type Service interface {
	GetProductsForCheckoutBySKUs(ctx context.Context, productSKUs []entity.SKU) (products []entity.Product)
	AssignProductsToMap(products []entity.Product) map[entity.SKU]entity.Product
}

func NewService(repository Repository) Service {
	return &service{repository}
}

// GetProductsForCheckoutBySKUs get products by SKUs for checkout things, so it may apply database transaction locking for product stock consistency
func (s *service) GetProductsForCheckoutBySKUs(ctx context.Context, productSKUs []entity.SKU) (products []entity.Product) {
	return s.repository.GetProductsForCheckoutBySKUs(ctx, productSKUs)
}

// assignProductsToMap transform array of product into map of product (overlapped product is merged into one)
func (s *service) AssignProductsToMap(products []entity.Product) map[entity.SKU]entity.Product {
	productMap := map[entity.SKU]entity.Product{}
	for _, product := range products {
		if _, ok := productMap[product.SKU]; !ok {
			productMap[product.SKU] = product
		}
		productMapProduct := productMap[product.SKU]
		productMapProduct.BuyQty += 1
		productMap[product.SKU] = productMapProduct

	}
	return productMap
}
