package main

import (
	"log"
	"net/http"

	"github.com/farhaan/kuncie/handler"
	"github.com/farhaan/kuncie/internal/checkout"
	"github.com/farhaan/kuncie/internal/discount"
	"github.com/farhaan/kuncie/internal/product"
)

func main() {
	// initiate repositories
	productRepository := product.NewRepository()
	discountRepository := discount.NewRepository()

	// initiate services
	productSvc := product.NewService(productRepository)
	discountSvc := discount.NewService(discountRepository)
	checkoutSvc := checkout.NewService(productSvc, discountSvc)

	// initiate handler
	checkoutHandler := handler.NewCheckoutHandler(checkoutSvc)

	http.HandleFunc("/checkout", checkoutHandler.Checkout)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
