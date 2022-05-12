package schema

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type CostComponent struct {
	Name                 string
	Unit                 string
	UnitMultiplier       decimal.Decimal
	IgnoreIfMissingPrice bool
	ProductFilter        *ProductFilter
	PriceFilter          *PriceFilter
	HourlyQuantity       *decimal.Decimal
	MonthlyQuantity      *decimal.Decimal
	MonthlyDiscountPerc  float64
	price                decimal.Decimal
	customPrice          *decimal.Decimal
	priceHash            string
	HourlyCost           *decimal.Decimal
	MonthlyCost          *decimal.Decimal
	tierPrices           []decimal.Decimal
	tierPriceHashes      []string
	TierQuantities       []decimal.Decimal
	TierNames            []string
	MonthlyTierCost      []*decimal.Decimal
}

func (c *CostComponent) CalculateCosts() {

	if len(c.tierPrices) > 0 {
		fmt.Printf("using tiered pricing...")
		var runningTotal decimal.Decimal = decimal.Zero
		for i := 0; i < len(c.tierPrices); i++ {
			runningTotal = runningTotal.Add(c.tierPrices[i].Mul(c.TierQuantities[i]))
			c.MonthlyTierCost = append(c.MonthlyTierCost, decimalPtr(c.tierPrices[i].Mul(c.TierQuantities[i])))
		}
		c.MonthlyCost = decimalPtr(runningTotal)
	} else {
		c.fillQuantities()

		if c.HourlyQuantity != nil {
			c.HourlyCost = decimalPtr(c.price.Mul(*c.HourlyQuantity))
		}
		if c.MonthlyQuantity != nil {
			discountMul := decimal.NewFromFloat(1.0 - c.MonthlyDiscountPerc)
			c.MonthlyCost = decimalPtr(c.price.Mul(*c.MonthlyQuantity).Mul(discountMul))
		}
	}
}

func (c *CostComponent) fillQuantities() {
	if c.MonthlyQuantity != nil && c.HourlyQuantity == nil {
		c.HourlyQuantity = decimalPtr(c.MonthlyQuantity.Div(HourToMonthUnitMultiplier))
	} else if c.HourlyQuantity != nil && c.MonthlyQuantity == nil {
		c.MonthlyQuantity = decimalPtr(c.HourlyQuantity.Mul(HourToMonthUnitMultiplier))
	}
}

func (c *CostComponent) SetPrice(price decimal.Decimal) {
	c.price = price
}

func (c *CostComponent) Price() decimal.Decimal {
	return c.price
}

func (c *CostComponent) SetPriceHash(priceHash string) {
	c.priceHash = priceHash
}

func (c *CostComponent) PriceHash() string {
	return c.priceHash
}

func (c *CostComponent) SetCustomPrice(price *decimal.Decimal) {
	c.customPrice = price
}

func (c *CostComponent) CustomPrice() *decimal.Decimal {
	return c.customPrice
}

func (c *CostComponent) SetTierPrices(tierPrices []decimal.Decimal) {
	c.tierPrices = tierPrices
}

func (c *CostComponent) SetTierPriceHashes(tierPriceHashes []string) {
	c.tierPriceHashes = tierPriceHashes
}

func (c *CostComponent) SetTierQuantities(tierQuantities []decimal.Decimal) {
	c.TierQuantities = tierQuantities
}

func (c *CostComponent) SetTierNames(tierNames []string) {
	c.TierNames = tierNames
}

func (c *CostComponent) UnitMultiplierPrice() decimal.Decimal {
	return c.Price().Mul(c.UnitMultiplier)
}

func (c *CostComponent) UnitMultiplierHourlyQuantity() *decimal.Decimal {
	if c.HourlyQuantity == nil {
		return nil
	}
	m := c.HourlyQuantity.Div(c.UnitMultiplier)
	return &m
}

func (c *CostComponent) UnitMultiplierMonthlyQuantity() *decimal.Decimal {
	if c.MonthlyQuantity == nil {
		return nil
	}
	m := c.MonthlyQuantity.Div(c.UnitMultiplier)
	return &m
}
