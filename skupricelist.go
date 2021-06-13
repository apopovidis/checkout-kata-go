package main

import (
	"errors"
	"fmt"
	"strings"
)

type SkuPriceList struct {
	items []SkuPrice
}

func NewSkuPriceList() SkuPriceList {
	return SkuPriceList{}
}

func (s *SkuPriceList) GetItems() []SkuPrice {
	return s.items
}

func (s *SkuPriceList) AddItem(name string, price int) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name cannot be empty")
	}

	if price <= 0 {
		return errors.New("price must be greater than zero")
	}

	var sku = s.GetSku(name)
	if sku == nil {
		s.items = append(s.items, NewSkuPrice(name, price))
	}

	return fmt.Errorf("sku with name %s already exists", name)
}

func (s *SkuPriceList) GetSku(name string) *SkuPrice {
	for _, i := range s.items {
		if i.GetName() == name {
			return &i
		}
	}
	return nil
}

func (s *SkuPriceList) GetSkuPrice(name string) (int, error) {
	var sku = s.GetSku(name)
	if sku == nil {
		return 0, errors.New(fmt.Sprintf("sku with name %s does not exist", name))
	}
	return sku.GetPrice(), nil
}
