package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var SkuSpecialPriceMap = map[string]map[int]int{
	"A": {
		3: 130,
	},
	"B": {
		2: 45,
	},
	"C": nil,
	"D": nil,
}

var SkuSpecialPriceMapExtraOffer = map[string]map[int]int{
	"A": {
		3: 130,
		4: 140, // extra
	},
	"B": {
		2: 45,
	},
	"C": nil,
	"D": nil,
}

// SkuPriceMap defines a map of skus and their corresponding prices
var SkuPriceMap = map[string]int{
	"A": 50,
	"B": 30,
	"C": 20,
	"D": 15,
}

// TestScenarios defines a map of standard scenarios for single sku scanning, with the key being the name of the
// test itself and the value the expected result
var TestScenarios = map[string]int{
	"A":    50,
	"B":    30,
	"C":    20,
	"D":    15,
	"BB":   45,
	"BAB":  95,
	"DDAA": 130,
	// the differences
	"AAA":        130,
	"AAAA":       180,
	"AAAAA":      230,
	"AAAAAA":     260,
	"AAAAAAA":    310,
	"AAAAAAAA":   360,
	"AAAAAAAAB":  390,
	"AAAAAAAABB": 405,
}

// TestScenariosExtraOffer defines a map of standard scenarios along with some scenarios affected by the extra special offer
// - for single sku scanning, with the key being the name of the test itself and the value the expected result
var TestScenariosExtraOffer = map[string]int{
	"A":    50,
	"B":    30,
	"C":    20,
	"D":    15,
	"BB":   45,
	"BAB":  95,
	"DDAA": 130,
	// the differences
	"AAA":        130,
	"AAAA":       140,
	"AAAAA":      190,
	"AAAAAA":     240,
	"AAAAAAA":    270,
	"AAAAAAAAB":  310,
	"AAAAAAAABB": 325,
}

// TestScenariosMultiple defines a map of scenarios for multi sku scanning, with the key being the name of the
// test itself and the value the expected result
var TestScenariosMultiple = map[string]int{
	"C":  20,
	"CC": 40,
}

type scanScenario struct {
	name      string
	count     int
	expectErr bool
}

var newScanScenario = func(name string, count int, expectErr bool) scanScenario {
	return scanScenario{
		name,
		count,
		expectErr,
	}
}

func TestScanValidations(t *testing.T) {
	var checkout = NewCheckout(GetSkuPriceList(SkuPriceMap), GetSkuSpecialPriceList(SkuSpecialPriceMap))

	var scenarios = map[string]scanScenario{
		"empty sku name and any count":           newScanScenario("", 1, true),
		"sku name only spaces and any count":     newScanScenario(" ", 1, true),
		"valid sku name and negative count":      newScanScenario("A", -1, true),
		"valid sku name and count equal to zero": newScanScenario("A", 0, true),
		"valid sku name and valid count":         newScanScenario("A", 1, false),
	}

	for sCase, s := range scenarios {
		t.Run(sCase, func(t *testing.T) {
			var (
				err        = checkout.Scan(s.name, s.count)
				errMessage string
			)
			if err != nil {
				errMessage = fmt.Sprintf("Expected no error, but instead got: '%s' for scenario '%s'", err.Error(), sCase)
			}
			switch {
			case s.expectErr:
				assert.Error(t, err, errMessage)
			default:
				assert.NoError(t, err, errMessage)
			}
		})
	}
}

func TestGetTotalPriceWhenNoSkuIsScannedYet(t *testing.T) {
	var checkout = NewCheckout(GetSkuPriceList(SkuPriceMap), GetSkuSpecialPriceList(SkuSpecialPriceMap))
	assert.Equal(t, 0, checkout.GetTotalPrice())
}

func TestScanSkusAndGetTotalPriceForSingleSku(t *testing.T) {
	runTestScenarios(t, SkuPriceMap, SkuSpecialPriceMap, TestScenarios)
	runTestScenarios(t, SkuPriceMap, SkuSpecialPriceMapExtraOffer, TestScenariosExtraOffer)
}

func TestScanSkusAndGetTotalPriceForMultipleSkus(t *testing.T) {
	var checkout = NewCheckout(GetSkuPriceList(SkuPriceMap), GetSkuSpecialPriceList(SkuSpecialPriceMap))

	for s, expected := range TestScenariosMultiple {
		t.Run(s, func(t *testing.T) {
			items := strings.Split(s, "")
			require.NoError(t, checkout.Scan(items[0], len(items)))

			assert.Equal(t, expected, checkout.GetTotalPrice())
		})
	}
}

func runTestScenarios(t *testing.T, skuPriceMap map[string]int, skuSpecialPriceMap map[string]map[int]int, testScenarios map[string]int) {
	t.Helper()

	var checkout = NewCheckout(GetSkuPriceList(skuPriceMap), GetSkuSpecialPriceList(skuSpecialPriceMap))

	for s, expected := range testScenarios {
		t.Run(s, func(t *testing.T) {
			items := strings.Split(s, "")
			for _, item := range items {
				require.NoError(t, checkout.Scan(item, 1))
			}

			assert.Equal(t, expected, checkout.GetTotalPrice())
		})
	}
}

func GetSkuPriceList(skuPriceMap map[string]int) (skuPriceList SkuPriceList) {
	skuPriceList = NewSkuPriceList()
	for key, value := range skuPriceMap {
		skuPriceList.AddItem(key, value)
	}
	return
}

func GetSkuSpecialPriceList(skuSpecialPriceMap map[string]map[int]int) (skuSpecialPriceList SkuSpecialPriceList) {
	skuSpecialPriceList = NewSkuSpecialPriceList()
	for key1 := range skuSpecialPriceMap {
		if len(skuSpecialPriceMap[key1]) == 0 {
			continue
		}

		for key2, value2 := range skuSpecialPriceMap[key1] {
			skuSpecialPriceList.AddItem(key1, key2, value2)
		}
	}
	return
}
