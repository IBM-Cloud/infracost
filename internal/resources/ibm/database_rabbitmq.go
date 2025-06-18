package ibm

import (
	"fmt"
	"strconv"

	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

func GetRabbitMqCostComponents(r *Database) []*schema.CostComponent {
	if r.Flavor != "" && r.Flavor != "multitenant" {
		return []*schema.CostComponent{
			RabbitMqHostFlavorComponent(r),
			RabbitMqDiskCostComponent(r),
		}
	} else {
		return []*schema.CostComponent{
			RabbitMqRAMCostComponent(r),
			RabbitMqDiskCostComponent(r),
			RabbitMqVirtualProcessorCoreCostComponent(r),
		}
	}
}

func RabbitMqHostFlavorComponent(r *Database) *schema.CostComponent {

	CostUnit := map[string]string{
		"b3c.4x16.encrypted":   "HOST_FOUR_SIXTEEN",
		"b3c.8x32.encrypted":   "HOST_EIGHT_THIRTYTWO",
		"m3c.8x64.encrypted":   "HOST_EIGHT_SIXTYFOUR",
		"b3c.16x64.encrypted":  "HOST_SIXTEEN_SIXTYFOUR",
		"b3c.32x128.encrypted": "HOST_THIRTYTWO_ONEHUNDREDTWENTYEIGHT",
		"m3c.30x240.encrypted": "HOST_THIRTY_TWOHUNDREDFORTY",
	}

	return &schema.CostComponent{
		Name:            fmt.Sprintf("Host Flavour - %s (%s members)", r.Flavor, strconv.FormatInt(r.Members, 10)),
		Unit:            "Instance",
		UnitMultiplier:  decimal.NewFromInt(1), // Final quantity for this cost component will be divided by this amount
		MonthlyQuantity: decimalPtr(decimal.NewFromInt(r.Members)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Location),
			Service:       strPtr("messages-for-rabbitmq"),
			ProductFamily: strPtr("service"),
			AttributeFilters: []*schema.AttributeFilter{
				{
					Key: "planName", Value: strPtr(r.Plan),
				},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(CostUnit[r.Flavor]),
		},
	}
}

func RabbitMqRAMCostComponent(r *Database) *schema.CostComponent {

	return &schema.CostComponent{
		Name:            fmt.Sprintf("RAM (%s members)", strconv.FormatInt(r.Members, 10)),
		Unit:            "GB",
		UnitMultiplier:  decimal.NewFromInt(1),                                                         // Final quantity for this cost component will be divided by this amount
		MonthlyQuantity: decimalPtr(decimal.NewFromFloat(float64(r.Memory*r.Members) / float64(1024))), // Convert to GB
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Location),
			Service:       strPtr("messages-for-rabbitmq"),
			ProductFamily: strPtr("service"),
			AttributeFilters: []*schema.AttributeFilter{
				{
					Key: "planName", Value: strPtr(r.Plan),
				},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("GIGABYTE_MONTHS_RAM"),
		},
	}
}

func RabbitMqDiskCostComponent(r *Database) *schema.CostComponent {
	return &schema.CostComponent{
		Name:            fmt.Sprintf("Disk (%s members)", strconv.FormatInt(r.Members, 10)),
		Unit:            "GB",
		UnitMultiplier:  decimal.NewFromInt(1),                                                       // Final quantity for this cost component will be divided by this amount
		MonthlyQuantity: decimalPtr(decimal.NewFromFloat(float64(r.Disk*r.Members) / float64(1024))), // Convert to GB
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Location),
			Service:       strPtr("messages-for-rabbitmq"),
			ProductFamily: strPtr("service"),
			AttributeFilters: []*schema.AttributeFilter{
				{
					Key: "planName", Value: strPtr(r.Plan),
				},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("GIGABYTE_MONTHS_DISK"),
		},
	}
}

func RabbitMqVirtualProcessorCoreCostComponent(r *Database) *schema.CostComponent {
	return &schema.CostComponent{
		Name:            fmt.Sprintf("Virtual Processor Cores (%s members)", strconv.FormatInt(r.Members, 10)),
		Unit:            "CPU",
		UnitMultiplier:  decimal.NewFromInt(1), // Final quantity for this cost component will be divided by this amount
		MonthlyQuantity: decimalPtr(decimal.NewFromInt(r.CPU * r.Members)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Location),
			Service:       strPtr("messages-for-rabbitmq"),
			ProductFamily: strPtr("service"),
			AttributeFilters: []*schema.AttributeFilter{
				{
					Key: "planName", Value: strPtr(r.Plan),
				},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("VIRTUAL_PROCESSOR_CORES"),
		},
	}
}
