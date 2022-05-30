package product

import (
	"context"
	"encoding/json"
	"os"

	"github.com/farhaan/kuncie/entity"
)

type repository struct {
	productMap map[entity.SKU]entity.Product
}
type Repository interface {
	GetProductsForCheckoutBySKUs(ctx context.Context, productSKUs []entity.SKU) (products []entity.Product)
}

func NewRepository() Repository {
	return &repository{productMap: initiate()}
}

// GetProductsForCheckoutBySKUs get products by SKUs for checkout things, so it may apply database transaction locking for product stock consistency
func (r *repository) GetProductsForCheckoutBySKUs(ctx context.Context, productSKUs []entity.SKU) (products []entity.Product) {
	// TODO: implement transaction lock to avoid race condition that caused by product availability
	for _, sku := range productSKUs {
		if product, ok := r.productMap[sku]; ok {
			products = append(products, product)
		}
	}

	return products
}

// initiate initiate product data using fake json
func initiate() map[entity.SKU]entity.Product {
	jsonFile, err := os.Open("database/products.json")
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	var productMap map[entity.SKU]entity.Product

	jsonDecoder := json.NewDecoder(jsonFile)
	jsonDecoder.UseNumber()

	if err := jsonDecoder.Decode(&productMap); err != nil {
		panic(err)
	}

	return productMap
}
