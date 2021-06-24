package main

import (
	"github.com/pkg/errors"

	"sort"
	"strings"
)

// ICheckout is an interface that defines two major methods
// needed in order to support the basic functionality of checkout
type ICheckout interface {
	Scan(item string, count int) (err error)
	GetTotalPrice() int
}

type checkout struct {
	skuPriceList        SkuPriceList
	skuSpecialPriceList SkuSpecialPriceList
	scannedSkusMap      map[string]int
}

// NewCheckout is constructor for the checkout struct
func NewCheckout(
	skuPriceList SkuPriceList,
	skuSpecialPriceList SkuSpecialPriceList,
) ICheckout {
	// init the map in the constructor function as the struct has not been created yet
	var scannedSkusMap = make(map[string]int)
	for _, item := range skuPriceList.GetItems() {
		scannedSkusMap[item.GetName()] = 0
	}
	return &checkout{
		skuPriceList,
		skuSpecialPriceList,
		scannedSkusMap,
	}
}

// Scan is a method that allows to scan items with a given name and count
func (s *checkout) Scan(name string, count int) (err error) {
	var n = strings.TrimSpace(name)
	if n == "" {
		return errors.New("name cannot be empty")
	}

	if count <= 0 {
		return errors.New("count must be greater than 0")
	}

	if _, ok := s.scannedSkusMap[n]; !ok {
		return errors.Errorf("sku with name %s does not exist", n)
	}

	s.scannedSkusMap[n] += count

	return
}

func (s *checkout) GetTotalPrice() int {
	var res int
	for item := range s.scannedSkusMap {
		res += s.getTotalPriceForItem(item)
	}

	s.initScannedSkusMap()

	return res
}

func (s *checkout) initScannedSkusMap() {
	s.scannedSkusMap = make(map[string]int)
	for _, item := range s.skuPriceList.GetItems() {
		s.scannedSkusMap[item.GetName()] = 0
	}
}

func (s *checkout) getTotalPriceForItem(name string) int {
	var (
		total int
		count = s.scannedSkusMap[name]
	)

	var skuPriceForName = s.skuPriceList.GetSku(name).GetPrice()

	// if there is no offer for the name just do a total and return it	var skuPriceListForName = s.skuSpecialPriceList.GetSkus(name)
	var skuPriceListForName = s.skuSpecialPriceList.GetSkus(name)
	if len(skuPriceListForName) == 0 {
		return count * skuPriceForName
	}

	var (
		index           int
		skuOffersCounts = make([]int, len(skuPriceListForName))
	)
	for _, c := range skuPriceListForName {
		skuOffersCounts[index] = c.GetNumberOfUnits()
		index++
	}

	// sort the sku offer counts so that we start from the highest offer
	sort.Ints(skuOffersCounts)

	var remainingCount = count
	for i := len(skuOffersCounts) - 1; i >= 0; i-- {
		if remainingCount < skuOffersCounts[i] {
			continue
		}

		var (
			skuPriceListForNameAndNumberOfUnits = s.skuSpecialPriceList.GetSkus(name, skuOffersCounts[i])
			nextOfferPrice                      = skuPriceListForNameAndNumberOfUnits[0].GetTotalPrice()
			reminder                            = remainingCount % skuOffersCounts[i]
		)

		// the largest offer is satisfied so return
		if reminder == 0 {
			total += (remainingCount / skuOffersCounts[i]) * nextOfferPrice
			remainingCount = 0
			break
		}

		// any next offer
		total += ((remainingCount - reminder) / skuOffersCounts[i]) * nextOfferPrice
		remainingCount = reminder
	}

	// if we have remaining count still then add it to totals
	if remainingCount > 0 {
		total += remainingCount * skuPriceForName
	}

	return total
}
