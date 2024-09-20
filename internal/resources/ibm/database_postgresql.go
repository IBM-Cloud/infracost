package ibm

import (
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

func PostgresRAMCostComponent(r *Database) *schema.CostComponent {
	var R decimal.Decimal
	if r.PostgreSQL_Ram != nil {
		R = ConvertMBtoGB(decimal.NewFromInt(*r.PostgreSQL_Ram))
	} else { // set the default
		R = decimal.NewFromInt(1)
	}
	var m decimal.Decimal
	if r.PostgreSQL_Members != nil {
		m = decimal.NewFromInt(*r.PostgreSQL_Members)
	} else { // set the default
		m = decimal.NewFromInt(2)
	}

	cost := R.Mul(m)

	costComponent := schema.CostComponent{
		Name:            "RAM",
		Unit:            "GB-RAM",
		MonthlyQuantity: &cost,
		UnitMultiplier:  decimal.NewFromInt(1),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Location),
			Service:       strPtr("databases-for-postgresql"),
			ProductFamily: strPtr("service"),
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("GIGABYTE_MONTHS_RAM"),
		},
	}
	return &costComponent
}

func PostgresDiskCostComponent(r *Database) *schema.CostComponent {
	var d decimal.Decimal
	if r.PostgreSQL_Disk != nil {
		d = ConvertMBtoGB(decimal.NewFromInt(*r.PostgreSQL_Disk))
	} else { // set the default
		d = decimal.NewFromInt(5)
	}
	var m decimal.Decimal
	if r.PostgreSQL_Members != nil {
		m = decimal.NewFromInt(*r.PostgreSQL_Members)
	} else { // set the default
		m = decimal.NewFromInt(2)
	}

	cost := d.Mul(m)

	costComponent := schema.CostComponent{
		Name:            "Disk",
		Unit:            "GB-DISK",
		MonthlyQuantity: &cost,
		UnitMultiplier:  decimal.NewFromInt(1),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Location),
			Service:       strPtr("databases-for-postgresql"),
			ProductFamily: strPtr("service"),
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("GIGABYTE_MONTHS_DISK"),
		},
	}
	return &costComponent
}

func PostgresCoreCostComponent(r *Database) *schema.CostComponent {
	var c decimal.Decimal
	if r.PostgreSQL_Core != nil {
		c = decimal.NewFromInt(*r.PostgreSQL_Core)
	} else { // set the default
		c = decimal.NewFromInt(0)
	}
	var m decimal.Decimal
	if r.PostgreSQL_Members != nil {
		m = decimal.NewFromInt(*r.PostgreSQL_Members)
	} else { // set the default
		m = decimal.NewFromInt(2)
	}

	cost := c.Mul(m)

	costComponent := schema.CostComponent{
		Name:            "Core",
		Unit:            "Virtual Processor Core",
		MonthlyQuantity: &cost,
		UnitMultiplier:  decimal.NewFromInt(1),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Location),
			Service:       strPtr("databases-for-postgresql"),
			ProductFamily: strPtr("service"),
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("VIRTUAL_PROCESSOR_CORES"),
		},
	}
	return &costComponent
}

func GetPostgresCostComponents(r *Database) []*schema.CostComponent {
	return []*schema.CostComponent{
		PostgresRAMCostComponent(r),
		PostgresDiskCostComponent(r),
		PostgresCoreCostComponent(r),
	}
}
