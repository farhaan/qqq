package product

import (
	"reflect"
	"testing"

	"github.com/farhaan/kuncie/entity"
)

func Test_service_AssignProductsToMap(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		products []entity.Product
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[entity.SKU]entity.Product
	}{
		{
			name: "Merge Same Product",
			fields: fields{
				repository: nil,
			},
			args: args{
				products: []entity.Product{{}, {}},
			},
			want: map[entity.SKU]entity.Product{"": {
				BuyQty: 2,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repository: tt.fields.repository,
			}
			if got := s.AssignProductsToMap(tt.args.products); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.AssignProductsToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
