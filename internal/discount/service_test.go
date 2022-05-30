package discount

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/farhaan/kuncie/entity"
)

func Test_service_DiscountBuyXGetY(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		ctx            context.Context
		discount       entity.Discount
		productCartMap map[entity.SKU]entity.Product
	}
	tests := []struct {
		name                         string
		fields                       fields
		args                         args
		wantDiscountedProductCartMap map[entity.SKU]entity.Product
		wantOk                       bool
	}{
		{
			name:   "Buy Nokia to Get iPhone",
			fields: fields{},
			args: args{
				ctx: nil,
				discount: entity.Discount{
					Type:     1,
					Keywords: []string{"nokia"},
					Rule: entity.DiscountRule{
						StartDate:    time.Now(),
						EndDate:      time.Now().Add(time.Hour),
						Accumulative: false,
						ItemValue: map[string]entity.DiscountRuleItemCriteria{
							"nokia": entity.DiscountRuleItemCriteria{
								Min: 1,
								Max: 1,
							},
							"iPhone": entity.DiscountRuleItemCriteria{
								Min: 0,
								Max: 1,
							},
						},
					},
					AmountType: 3,
					Amount:     1000.00,
				},
				productCartMap: map[entity.SKU]entity.Product{"nokia": {
					SKU:      "nokia",
					Name:     "nokia",
					Price:    5000.00,
					StockQty: 10,
					BuyQty:   1,
				}},
			},
			wantDiscountedProductCartMap: map[entity.SKU]entity.Product{
				"nokia":  {SKU: "nokia", Name: "nokia", Price: 5000.00, StockQty: 10, BuyQty: 1},
				"iPhone": {BuyQty: 1}},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repository: tt.fields.repository,
			}
			gotDiscountedProductCartMap, gotOk := s.DiscountBuyXGetY(tt.args.ctx, tt.args.discount, tt.args.productCartMap)
			if !reflect.DeepEqual(gotDiscountedProductCartMap, tt.wantDiscountedProductCartMap) {
				t.Errorf("service.DiscountBuyXGetY() gotDiscountedProductCartMap = %v, want %v", gotDiscountedProductCartMap, tt.wantDiscountedProductCartMap)
			}
			if gotOk != tt.wantOk {
				t.Errorf("service.DiscountBuyXGetY() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_service_DiscountBuyN(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		ctx            context.Context
		discount       entity.Discount
		productCartMap map[entity.SKU]entity.Product
	}
	tests := []struct {
		name                         string
		fields                       fields
		args                         args
		wantDiscountedProductCartMap map[entity.SKU]entity.Product
		wantOk                       bool
	}{
		{
			name:   "Buy 3 iPhone to get 10% bonus",
			fields: fields{},
			args: args{
				ctx: nil,
				discount: entity.Discount{
					Type:     1,
					Keywords: []string{"iPhone"},
					Rule: entity.DiscountRule{
						StartDate:    time.Now(),
						EndDate:      time.Now().Add(time.Hour),
						Accumulative: false,
						ItemValue: map[string]entity.DiscountRuleItemCriteria{
							"iPhone": {
								Min: 3,
								Max: 100,
							},
						},
					},
					AmountType: 3,
					Amount:     1000.00,
				},
				productCartMap: map[entity.SKU]entity.Product{"iPhone": {
					SKU:      "iPhone",
					Name:     "iPhone",
					Price:    1000.00,
					StockQty: 10,
					BuyQty:   3,
				}},
			},
			wantDiscountedProductCartMap: map[entity.SKU]entity.Product{"iPhone": {
				SKU:      "iPhone",
				Name:     "iPhone",
				Price:    1000.00,
				StockQty: 10,
				BuyQty:   3,
			}},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repository: tt.fields.repository,
			}
			gotDiscountedProductCartMap, gotOk := s.DiscountBuyN(tt.args.ctx, tt.args.discount, tt.args.productCartMap)
			if !reflect.DeepEqual(gotDiscountedProductCartMap, tt.wantDiscountedProductCartMap) {
				t.Errorf("service.DiscountBuyN() gotDiscountedProductCartMap = %v, want %v", gotDiscountedProductCartMap, tt.wantDiscountedProductCartMap)
			}
			if gotOk != tt.wantOk {
				t.Errorf("service.DiscountBuyN() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
