package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getSkuPriceList() (skuPriceList SkuPriceList) {
	skuPriceList = NewSkuPriceList()
	for key, value := range SkuPriceMap {
		skuPriceList.AddItem(key, value)
	}
	return
}

func getSkuSpecialPriceList() (skuSpecialPriceList SkuSpecialPriceList) {
	skuSpecialPriceList = NewSkuSpecialPriceList()
	for key1 := range SkuSpecialPriceMap {
		if len(SkuSpecialPriceMap[key1]) == 0 {
			continue
		}

		for key2, value2 := range SkuSpecialPriceMap[key1] {
			skuSpecialPriceList.AddItem(key1, key2, value2)
		}
	}
	return
}

func TestScanSkusAndGetTotalPriceSingle(t *testing.T) {
	var checkout = NewCheckout(getSkuPriceList(), getSkuSpecialPriceList())

	for s, expected := range TestScenarios {
		t.Run(s, func(t *testing.T) {
			items := strings.Split(s, "")
			for _, item := range items {
				require.NoError(t, checkout.Scan(item, 1))
			}

			assert.Equal(t, expected, checkout.GetTotalPrice())
		})
	}
}

func TestScanSkusAndGetTotalPriceMultiple(t *testing.T) {
	var checkout = NewCheckout(getSkuPriceList(), getSkuSpecialPriceList())

	for s, expected := range TestScenariosMultiple {
		t.Run(s, func(t *testing.T) {
			items := strings.Split(s, "")
			require.NoError(t, checkout.Scan(items[0], len(items)))

			assert.Equal(t, expected, checkout.GetTotalPrice())
		})
	}
}
