package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// Graduated Tier pricing model
const DNS_SERVICES_PROGRAMMATIC_PLAN_NAME string = "standard-dns"

func GetDNSServicesCostComponents(r *ResourceInstance) []*schema.CostComponent {
	if r.Plan == DNS_SERVICES_PROGRAMMATIC_PLAN_NAME {
		return []*schema.CostComponent{
			// DNSServicesZonesCostComponents(r),
			// DNSServicesPoolsPerHourCostComponents(r),
			// DNSServicesGLBInstancesPerHourCostComponents(r),
			// DNSServicesHealthChecksCostComponents(r),
			// DNSServicesCustomResolverLocationsPerHourCostComponents(r),
			// DNSServicesMillionCustomResolverExternalQueriesCostComponents(r),
			// DNSServicesMillionDNSQueriesCostComponents(r)
		}
	} else {
		costComponent := schema.CostComponent{
			Name:            fmt.Sprintf("Plan %s with customized pricing", r.Plan),
			UnitMultiplier:  decimal.NewFromInt(1), // Final quantity for this cost component will be divided by this amount
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
		return []*schema.CostComponent{
			&costComponent,
		}
	}
}

// Unit: ITEMS (Linear Tier)
func DNSServicesZonesCostComponents(r *ResourceInstance) *schema.CostComponent {

	var quantity *decimal.Decimal = decimalPtr(decimal.NewFromFloat(*r.DNSServices_Zones)) // Quantity of current cost component (i.e. Number of zones)

	return &schema.CostComponent{
		Name:            "Zones",
		Unit:            "Zones",
		UnitMultiplier:  decimal.NewFromFloat(1), // Final quantity for this cost component will be divided by this amount
		MonthlyQuantity: quantity,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("ITEMS"),
		},
	}
}

// Unit: NUMBERGLB (Linear Tier)
func DNSServicesPoolsPerHourCostComponents(r *ResourceInstance) *schema.CostComponent {

	var quantity *decimal.Decimal = decimalPtr(decimal.NewFromFloat(*r.DNSServices_PoolsPerHour)) // Quantity of current cost component (i.e. Number of zones)

	return &schema.CostComponent{
		Name:            "Pools Per Hour",
		Unit:            "Pools Per Hour",
		UnitMultiplier:  decimal.NewFromFloat(1), // Final quantity for this cost component will be divided by this amount
		MonthlyQuantity: quantity,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("NUMBERGLB"),
		},
	}
}

// Unit: NUMBERPOOLS (Linear Tier)
func DNSServicesGLBInstancesPerHourCostComponents(r *ResourceInstance) *schema.CostComponent {

	var quantity *decimal.Decimal = decimalPtr(decimal.NewFromFloat(*r.DNSServices_GLBInstancesPerHour)) // Quantity of current cost component (i.e. Number of zones)

	return &schema.CostComponent{
		Name:            "GLB Instances Per Hour",
		Unit:            "GLB Instances Per Hour",
		UnitMultiplier:  decimal.NewFromFloat(1), // Final quantity for this cost component will be divided by this amount
		MonthlyQuantity: quantity,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("NUMBERPOOLS"),
		},
	}
}

// Unit: NUMBERHEALTHCHECK (Linear Tier)
func DNSServicesHealthChecksCostComponents(r *ResourceInstance) *schema.CostComponent {

	var quantity *decimal.Decimal = decimalPtr(decimal.NewFromFloat(*r.DNSServices_HealthChecks)) // Quantity of current cost component (i.e. Number of zones)

	return &schema.CostComponent{
		Name:            "Health Checks",
		Unit:            "Health Checks",
		UnitMultiplier:  decimal.NewFromFloat(1), // Final quantity for this cost component will be divided by this amount
		MonthlyQuantity: quantity,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("NUMBERHEALTHCHECK"),
		},
	}
}

// Unit: RESOLVERLOCATIONS (Linear Tier)
func DNSServicesCustomResolverLocationsPerHourCostComponents(r *ResourceInstance) *schema.CostComponent {

	var quantity *decimal.Decimal = decimalPtr(decimal.NewFromFloat(*r.DNSServices_CustomResolverLocationsPerHour)) // Quantity of current cost component (i.e. Number of zones)

	return &schema.CostComponent{
		Name:            "Custom Resolver Locations Per Hour",
		Unit:            "Custom Resolver Locations Per Hour",
		UnitMultiplier:  decimal.NewFromFloat(1), // Final quantity for this cost component will be divided by this amount
		MonthlyQuantity: quantity,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("RESOLVERLOCATIONS"),
		},
	}
}

// // Unit: MILLION_ITEMS_CREXTERNALQUERIES (Graduated Tier)
// func DNSServicesMillionCustomResolverExternalQueriesCostComponents(r *ResourceInstance) *schema.CostComponent {
// }

// // Unit: MILLION_ITEMS (Graduated Tier)
// func DNSServicesMillionDNSQueriesCostComponents(r *ResourceInstance) *schema.CostComponent {
// }
