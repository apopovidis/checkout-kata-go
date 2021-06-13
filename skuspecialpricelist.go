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
	if strings.TrimSpace(name) == "" {
		return errors.New("name cannot be empty")
	}

	if numberOfUnits <= 0 {
		return errors.New("numberOfUnits must be greater than zero")
	}

	if totalPrice <= 0 {
		return errors.New("totalPrice must be greater than zero")
	}

	var sku = s.GetSkus(name, numberOfUnits)
	if sku == nil {
		s.items = append(s.items, NewSkuSpecialPrice(name, numberOfUnits, totalPrice))
	}

	return fmt.Errorf("sku with name %s and numberOfUnits %d already exists", name, numberOfUnits)
}

func (s *SkuSpecialPriceList) GetSkus(name string, numberOfUnits ...int) (skus []SkuSpecialPrice) {
	for _, i := range s.items {
		if numberOfUnits != nil && numberOfUnits[0] > 0 {
			if i.GetName() == name && i.GetNumberOfUnits() == numberOfUnits[0] {
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
