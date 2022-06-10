package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

const (
	readCapacity        = 100
	writeCapacity       = 50
	globalQueryCapacity = 5
)

// Resource information: https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/resources/cloudant
// Pricing information: https://www.ibm.com/cloud/cloudant/pricing
type Cloudant struct {
	Address  string
	Region   string
	Plan     string
	Capacity int64

	AdditionalConsumptionStorageGB *int64 `infracost_usage:"additional_consumption_storage_gb"`
}

// PopulateUsage parses the u schema.UsageData into the Cloudant.
// It uses the `infracost_usage` struct tags to populate data into the Cloudant.
func (r *Cloudant) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// CloudantUsageSchema defines a list which represents the usage schema of Cloudant.
var CloudantUsageSchema = []*schema.UsageItem{
	{Key: "additional_consumption_storage_gb", DefaultValue: 0},
}

func (r *Cloudant) BuildResource() *schema.Resource {

	if r.Plan == "lite" {
		return &schema.Resource{
			Name:      r.Address,
			NoPrice:   true,
			IsSkipped: true,
		}
	}

	costComponents := []*schema.CostComponent{
		r.readCapacityCostComponent(),
		r.writeCapacityCostComponent(),
		r.globalQueryCostComponent(),
		r.additionalCostComponent(),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    CloudantUsageSchema,
		CostComponents: costComponents,
	}
}

func (r *Cloudant) readCapacityCostComponent() *schema.CostComponent {
	return &schema.CostComponent{
		Name:            fmt.Sprintf("Read capacity (%d Reads/sec)", r.Capacity*readCapacity),
		Unit:            "capacity",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: decimalPtr(decimal.NewFromInt(r.Capacity)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Cloudant"),
			ProductFamily: strPtr("Databases"),
			AttributeFilters: []*schema.AttributeFilter{
				{
					Key:   "component",
					Value: strPtr("readCapacity"),
				},
			},
		},
	}
}

func (r *Cloudant) writeCapacityCostComponent() *schema.CostComponent {
	return &schema.CostComponent{
		Name:            fmt.Sprintf("Write capacity (%d writes/sec)", r.Capacity*writeCapacity),
		Unit:            "capacity",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: decimalPtr(decimal.NewFromInt(r.Capacity)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Cloudant"),
			ProductFamily: strPtr("Databases"),
			AttributeFilters: []*schema.AttributeFilter{
				{
					Key:   "component",
					Value: strPtr("writeCapacity"),
				},
			},
		},
	}
}

func (r *Cloudant) globalQueryCostComponent() *schema.CostComponent {
	return &schema.CostComponent{
		Name:            fmt.Sprintf("Global query capacity (%d global queries/sec)", r.Capacity*globalQueryCapacity),
		Unit:            "capacity",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: decimalPtr(decimal.NewFromInt(r.Capacity)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Cloudant"),
			ProductFamily: strPtr("Databases"),
			AttributeFilters: []*schema.AttributeFilter{
				{
					Key:   "component",
					Value: strPtr("globalQuery"),
				},
			},
		},
	}
}

// additionalCostComponent returns a cost component for additional consumption-based charges per GB.
func (r *Cloudant) additionalCostComponent() *schema.CostComponent {
	var quantity *decimal.Decimal
	if r.AdditionalConsumptionStorageGB != nil {
		quantity = decimalPtr(decimal.NewFromInt(*r.AdditionalConsumptionStorageGB))
	}

	return &schema.CostComponent{
		Name:            "Additional consumption-based charges",
		Unit:            "GB",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: quantity,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Cloudant"),
			ProductFamily: strPtr("Databases"),
			AttributeFilters: []*schema.AttributeFilter{
				{
					Key:   "component",
					Value: strPtr("additionalCost"),
				},
			},
		},
	}
}
