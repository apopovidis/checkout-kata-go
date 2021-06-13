package main

type SkuPrice struct {
	name  string
	price int
}

func NewSkuPrice(
	name string,
	price int,
) SkuPrice {
	return SkuPrice{
		name,
		price,
	}
}

func (s *SkuPrice) GetName() string {
	return s.name
}

func (s *SkuPrice) GetPrice() int {
	return s.price
}
