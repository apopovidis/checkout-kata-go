package main

import (
	"errors"
	"fmt"
	"strings"
)

type SkuSpecialPriceList struct {
	items []SkuSpecialPrice
}

func NewSkuSpecialPriceList() SkuSpecialPriceList {
	return SkuSpecialPriceList{}
}

func (s *SkuSpecialPriceList) AddItem(name string, numberOfUnits int, totalPrice int) error {
	var n = strings.TrimSpace(name)
	if n == "" {
		return errors.New("name cannot be empty")
	}

	if numberOfUnits <= 0 {
		return errors.New("numberOfUnits must be greater than zero")
	}

	if totalPrice <= 0 {
		return errors.New("totalPrice must be greater than zero")
	}

	var sku = s.GetSkus(n, numberOfUnits, totalPrice)
	if sku != nil {
		return fmt.Errorf("sku with name %s, numberOfUnits %d and total price %d already exists", n, numberOfUnits, totalPrice)
	}

	s.items = append(s.items, NewSkuSpecialPrice(n, numberOfUnits, totalPrice))

	return nil
}

// if args are present then the first index should hold the number of units, and the second the total price
func (s *SkuSpecialPriceList) GetSkus(name string, args ...int) (skus []SkuSpecialPrice) {
	for _, i := range s.items {
		if len(args) == 1 {
			if i.GetName() == name && i.GetNumberOfUnits() == args[0] {
				skus = append(skus, i)
			}
		} else if len(args) == 2 {
			if i.GetName() == name && i.GetNumberOfUnits() == args[0] && i.GetTotalPrice() == args[1] {
				skus = append(skus, i)
			}
		} else {
			if i.GetName() == name {
				skus = append(skus, i)
			}
		}
	}
	return
}

func (s *SkuSpecialPriceList) GetSkuSpecialPrice(name string, numberOfUnits int) (int, error) {
	var skus = s.GetSkus(name, numberOfUnits)
	if skus == nil {
		return 0, fmt.Errorf("sku with name %s and numberOfUnits %d does not exist", name, numberOfUnits)
	}
	return skus[0].GetTotalPrice(), nil
}
