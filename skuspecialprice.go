package main

type SkuSpecialPrice struct {
	name          string
	numberOfUnits int
	totalPrice    int
}

func NewSkuSpecialPrice(
	name string,
	numberOfUnits int,
	totalPrice int,
) SkuSpecialPrice {
	return SkuSpecialPrice{
		name,
		numberOfUnits,
		totalPrice,
	}
}

func (s *SkuSpecialPrice) GetName() string {
	return s.name
}

func (s *SkuSpecialPrice) GetNumberOfUnits() int {
	return s.numberOfUnits
}

func (s *SkuSpecialPrice) GetTotalPrice() int {
	return s.totalPrice
}
