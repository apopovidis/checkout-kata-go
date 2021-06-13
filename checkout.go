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
	return &checkout{
		skuPriceList:        skuPriceList,
		skuSpecialPriceList: skuSpecialPriceList,
	}
}

// Scan is a method that allows to scan items with a given name and count
func (s *checkout) Scan(name string, count int) (err error) {
	if strings.TrimSpace(name) == "" {
		return errors.New("name cannot be empty")
	}

	if count <= 0 {
		return errors.New("count must be greater than 0")
	}

	s.initScannedSkusMap(false)

	if _, ok := s.scannedSkusMap[name]; !ok {
		return errors.Errorf("sku with name %s does not exist", name)
	}

	s.scannedSkusMap[name] += count

	return
}

func (s *checkout) GetTotalPrice() int {
	var res int
	for item := range s.scannedSkusMap {
		res += s.getTotalPriceForItem(item)
	}

	s.initScannedSkusMap(true)

	return res
}

func (s *checkout) initScannedSkusMap(forceInit bool) {
	if (!forceInit) && len(s.scannedSkusMap) > 0 {
		return
	}

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
			total = (remainingCount / skuOffersCounts[i]) * nextOfferPrice
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
