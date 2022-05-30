# Checkout System

This repository is a checkout system 

## Installation

Install go first using homebrew or by downloading go installation file.

Clone this repository by using
```bash
git clone git@github.com:farhaan/qqq.git
```

Run go mod dependencies
```bash
go mod tidy
```

compile the binary by running the makefile
```bash
make build
```

run the outputted binary
```bash
./bin/app
```

## Directory Structure
```
database/
entity/
handler/                // act as a controller
├─ checkout.go          // controller for checkout path
internal/               // directory to main apps logic
├─ checkout/            
│  ├─ service.go        // actual business logic interface & implementation
├─ discount/
│  ├─ service.go        
│  ├─ repository.go     // (optional) any logic regarding persistence
├─ product/
│  ├─ service.go
│  ├─ repository.go
utility/                // directory to any utility function that accessed by the whole microservice 
main.go                 // entrypoints to build REST API Server
makefile                // helper command to prepare envar / code generated file/
go.mod
```


## Flow
When user checking out the items they sent the data like the following:
```bash

curl --request POST \
  --url http://<URL>/checkout \
  --header 'Content-Type: application/json' \
  --data '[
	"120P90","120P90","120P90"
]'

```

`120P90` is a product sku's, so if user would like to checking out same product, the can simply add the product sku's into the array.


### Checkout
When user invoking the checkout endpoint, router then catch the request and then forward the the request into the handler (checkout.go) and then the handler parsing the request into certain struct and forward the request into internal/checkout/service.go from there will be 
    - getting product cart details from internal/product/service.go
    - deducting and getting the applied discounts from internal/discount/service.go
after that user may retrieve the data like the following
```json
{
	"products": {
		"120P90": {
			"sku": "120P90",
			"name": "Google Home",
			"price": 49.99,
			"stockQty": 10,
			"buyQty": 3
		},
		"120P90--5057722652628842616": {
			"sku": "120P90--5057722652628842616",
			"name": "DISCOUNT AMOUNT",
			"price": -49.99,
			"stockQty": 0,
			"buyQty": 1
		}
	},
	"finalProducts": {
		"120P90": {
			"sku": "120P90",
			"name": "Google Home",
			"price": 49.99,
			"stockQty": 10,
			"buyQty": 3
		},
		"120P90--5057722652628842616": {
			"sku": "120P90--5057722652628842616",
			"name": "DISCOUNT AMOUNT",
			"price": -49.99,
			"stockQty": 0,
			"buyQty": 1
		}
	},
	"discounts": [
		{
			"type": 2,
			"keywords": [
				"120P90"
			],
			"rule": {
				"startDate": "2022-01-31T15:04:05+07:00",
				"endDate": "2022-05-31T15:04:05+07:00",
				"accumulative": false,
				"itemValue": {
					"120P90": {
						"min": 3,
						"max": 3
					}
				}
			},
			"amountType": 2,
			"amount": 49.99
		}
	],
	"totalBasePrice": 99.98,
	"finalPrice": 99.98
}

```

#### How does the discount works

