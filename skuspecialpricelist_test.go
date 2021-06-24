package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type skuSpecialPriceListScenario struct {
	name          string
	numberOfUnits int
	totalPrice    int
	expectErr     bool
}

func TestSkuSpecialPriceListAddItemValidations(t *testing.T) {
	var list = NewSkuSpecialPriceList()

	var scenarios = map[string]skuSpecialPriceListScenario{
		"empty name, any number of units and any total price":                 newSkuSpecialPriceListScenario("", -1, -1, true),
		"name white spaces only, any number of units and any total price":     newSkuSpecialPriceListScenario(" ", -1, -1, true),
		"valid name, negative number of units and any total price":            newSkuSpecialPriceListScenario("A", -1, -1, true),
		"valid name, number of units equal to zero and any total price":       newSkuSpecialPriceListScenario("A", 0, 0, true),
		"valid name, valid number of units and negative total price":          newSkuSpecialPriceListScenario("A", 1, -1, true),
		"valid name, valid number of units and total price equal to zero":     newSkuSpecialPriceListScenario("A", 1, 0, true),
		"valid name, valid number of units and valid total price":             newSkuSpecialPriceListScenario("A", 1, 1, false),
		"name with white spaces, valid number of units and valid total price": newSkuSpecialPriceListScenario(" B ", 2, 2, false),
	}

	for sCase, s := range scenarios {
		t.Run(sCase, func(t *testing.T) {
			var (
				err        = list.AddItem(s.name, s.numberOfUnits, s.totalPrice)
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

func TestSkuSpecialPriceListAddExistingItem(t *testing.T) {
	var (
		list       = NewSkuSpecialPriceList()
		err        = list.AddItem("A", 1, 1) // add sku A for first time
		errMessage string
	)
	if err != nil {
		errMessage = fmt.Sprintf("Expected no error, but instead got: '%s'", err.Error())
	}
	assert.NoError(t, err, errMessage)

	// add sku A again with a different number of units and different total price
	err = list.AddItem("A", 2, 2)
	if err != nil {
		errMessage = fmt.Sprintf("Expected no error, but instead got: '%s'", err.Error())
	}
	assert.NoError(t, err, errMessage)

	// add another sku A again with the same number of units and same total price
	assert.Equal(t, fmt.Errorf("sku with name %s, numberOfUnits %d and total price %d already exists", "A", 2, 2), list.AddItem("A", 2, 2))
}

func TestSkuSpecialPriceListGetSkus(t *testing.T) {
	var list = NewSkuSpecialPriceList()

	require.NoError(t, list.AddItem("A", 1, 1))
	require.NoError(t, list.AddItem("A", 2, 2))
	require.NoError(t, list.AddItem("A", 3, 3))

	// getting all sku A offers
	items := list.GetSkus("A")
	assert.Equal(t, 3, len(items))
	assert.Equal(t, "A", items[0].GetName())
	assert.Equal(t, 1, items[0].GetNumberOfUnits())
	assert.Equal(t, 1, items[0].GetTotalPrice())

	assert.Equal(t, "A", items[1].GetName())
	assert.Equal(t, 2, items[1].GetNumberOfUnits())
	assert.Equal(t, 2, items[1].GetTotalPrice())

	assert.Equal(t, "A", items[2].GetName())
	assert.Equal(t, 3, items[2].GetNumberOfUnits())
	assert.Equal(t, 3, items[2].GetTotalPrice())

	// getting all sku A offers with number of units
	items = list.GetSkus("A", 1)
	assert.Equal(t, 1, len(items))
	assert.Equal(t, "A", items[0].GetName())
	assert.Equal(t, 1, items[0].GetNumberOfUnits())
	assert.Equal(t, 1, items[0].GetTotalPrice())

	// getting all sku A offers with number of units and total price
	items = list.GetSkus("A", 1, 1)
	assert.Equal(t, 1, len(items))
	assert.Equal(t, "A", items[0].GetName())
	assert.Equal(t, 1, items[0].GetNumberOfUnits())
	assert.Equal(t, 1, items[0].GetTotalPrice())
}

func TestSkuSpecialPriceListGetSkuSpecialPrice(t *testing.T) {
	var list = NewSkuSpecialPriceList()

	_, err := list.GetSkuSpecialPrice("A", 1)
	assert.Equal(t, fmt.Errorf("sku with name %s and numberOfUnits %d does not exist", "A", 1), err)

	require.NoError(t, list.AddItem("A", 1, 1))

	totalPrice, err := list.GetSkuSpecialPrice("A", 1)
	require.NoError(t, err)
	assert.Equal(t, 1, totalPrice)
}

func newSkuSpecialPriceListScenario(name string, numberOfUnits int, totalPrice int, expectErr bool) skuSpecialPriceListScenario {
	return skuSpecialPriceListScenario{
		name,
		numberOfUnits,
		totalPrice,
		expectErr,
	}
}
