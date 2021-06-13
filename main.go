package main

import (
	"fmt"
	"log"
	"strings"
)

var SkuSpecialPriceMap = map[string]map[int]int{
	"A": map[int]int{
		3: 130,
		4: 140, // extra
	},
	"B": map[int]int{
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

// TestScenarios defines a map of scenarios for single sku scanning, with the key being the name of the
// test itself and the value the expected result
var TestScenarios = map[string]int{
	"A":     50,
	"B":     30,
	"C":     20,
	"D":     15,
	"AAA":   130,
	"BB":    45,
	"BAB":   95,
	"AAAAA": 190,
	"DDAA":  130,
}

// TestScenariosMultiple defines a map of scenarios for multi sku scanning, with the key being the name of the
// test itself and the value the expected result
var TestScenariosMultiple = map[string]int{
	"C":  20,
	"CC": 40,
}

func main() {
	var skuPriceList = NewSkuPriceList()
	for key, value := range SkuPriceMap {
		skuPriceList.AddItem(key, value)
	}

	var skuSpecialPriceList = NewSkuSpecialPriceList()
	for key1 := range SkuSpecialPriceMap {
		if len(SkuSpecialPriceMap[key1]) == 0 {
			continue
		}

		for key2, value2 := range SkuSpecialPriceMap[key1] {
			skuSpecialPriceList.AddItem(key1, key2, value2)
		}
	}

	var checkout = NewCheckout(skuPriceList, skuSpecialPriceList)

	var pass = true
	for s, expected := range TestScenarios {
		items := strings.Split(s, "")
		for _, item := range items {
			if err := checkout.Scan(item, 1); err != nil {
				log.Fatal(err)
			}
		}

		var result = checkout.GetTotalPrice()
		if expected != result {
			fmt.Printf("Failed at scenario %s. Expected %d, but got %d\n", s, expected, result)
			pass = false
		}
	}

	if pass {
		fmt.Println("PASS")
	}
}
