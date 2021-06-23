package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type skuPriceListScenario struct {
	name      string
	price     int
	expectErr bool
}

func TestSkuPriceListAddItemValidations(t *testing.T) {
	var list = NewSkuPriceList()

	var scenarios = map[string]skuPriceListScenario{
		"empty name and any price":               newSkuPriceListScenario("", -1, true),
		"name white spaces only and any price":   newSkuPriceListScenario(" ", -1, true),
		"valid name and negative price":          newSkuPriceListScenario("A", -1, true),
		"valid name and price equal to zero":     newSkuPriceListScenario("A", 0, true),
		"valid name and valid price":             newSkuPriceListScenario("A", 1, false),
		"name with white spaces and valid price": newSkuPriceListScenario(" B ", 2, false),
	}

	for sCase, s := range scenarios {
		t.Run(sCase, func(t *testing.T) {
			var (
				err        = list.AddItem(s.name, s.price)
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

func TestSkuPriceListAddExistingItem(t *testing.T) {
	var (
		list       = NewSkuPriceList()
		err        = list.AddItem("A", 1) // add sku A for first time
		errMessage string
	)
	if err != nil {
		errMessage = fmt.Sprintf("Expected no error, but instead got: '%s'", err.Error())
	}
	assert.NoError(t, err, errMessage)

	// add sku A again
	assert.Equal(t, fmt.Errorf("sku with name %s already exists", "A"), list.AddItem("A", 1))
}

func TestSkuPriceListGetItems(t *testing.T) {
	var list = NewSkuPriceList()

	require.NoError(t, list.AddItem("A", 1))
	require.NoError(t, list.AddItem("B", 2))

	items := list.GetItems()
	assert.Equal(t, 2, len(items))
	assert.Equal(t, "A", items[0].GetName())
	assert.Equal(t, 1, items[0].GetPrice())
	assert.Equal(t, "B", items[1].GetName())
	assert.Equal(t, 2, items[1].GetPrice())
}

func newSkuPriceListScenario(name string, price int, expectErr bool) skuPriceListScenario {
	return skuPriceListScenario{
		name,
		price,
		expectErr,
	}
}
