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
			DNSServicesZonesCostComponents(r),
			DNSServicesPoolsPerHourCostComponents(r),
			DNSServicesGLBInstancesPerHourCostComponents(r),
			DNSServicesHealthChecksCostComponents(r),
			DNSServicesCustomResolverLocationsPerHourCostComponents(r),
			DNSServicesMillionCustomResolverExternalQueriesCostComponents(r),
			DNSServicesMillionDNSQueriesCostComponents(r),
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

	var zones_included int = 1
	var quantity *decimal.Decimal

	additional_zones := *r.DNSServices_Zones - int64(zones_included)
	if additional_zones > 0 {
		quantity = decimalPtr(decimal.NewFromInt(additional_zones))
	} else {
		quantity = decimalPtr(decimal.NewFromInt(0))
	}

	costComponent := schema.CostComponent{
		Name:            "Additional Zones",
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
	return &costComponent
}

// Unit: NUMBERPOOLS (Linear Tier)
func DNSServicesPoolsPerHourCostComponents(r *ResourceInstance) *schema.CostComponent {

	var quantity *decimal.Decimal = decimalPtr(decimal.NewFromFloat(float64(*r.DNSServices_PoolsPerHour) * *r.DNSServices_PoolHours))

	costComponent := schema.CostComponent{
		Name:            "Pool Hours",
		Unit:            "Hours",
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
	return &costComponent
}

// Unit: NUMBERGLB (Linear Tier)
func DNSServicesGLBInstancesPerHourCostComponents(r *ResourceInstance) *schema.CostComponent {

	var quantity *decimal.Decimal = decimalPtr(decimal.NewFromFloat(float64(*r.DNSServices_GLBInstancesPerHour) * *r.DNSServices_GLBInstanceHours))

	costComponent := schema.CostComponent{
		Name:            "GLB Instance Hours",
		Unit:            "Hours",
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
	return &costComponent
}

// Unit: NUMBERHEALTHCHECK (Linear Tier)
func DNSServicesHealthChecksCostComponents(r *ResourceInstance) *schema.CostComponent {

	var quantity *decimal.Decimal = decimalPtr(decimal.NewFromInt(*r.DNSServices_HealthChecks))

	costComponent := schema.CostComponent{
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
	return &costComponent
}

// Unit: RESOLVERLOCATIONS (Linear Tier)
func DNSServicesCustomResolverLocationsPerHourCostComponents(r *ResourceInstance) *schema.CostComponent {

	var quantity *decimal.Decimal = decimalPtr(decimal.NewFromFloat(float64(*r.DNSServices_CustomResolverLocationsPerHour) * *r.DNSServices_CustomResolverLocationHours))

	costComponent := schema.CostComponent{
		Name:            "Custom Resolver Location Hours",
		Unit:            "Hours",
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
	return &costComponent
}

// Unit: MILLION_ITEMS_CREXTERNALQUERIES (Graduated Tier)
func DNSServicesMillionCustomResolverExternalQueriesCostComponents(r *ResourceInstance) *schema.CostComponent {

	var quantity *decimal.Decimal = decimalPtr(decimal.NewFromInt(*r.DNSServices_CustomResolverExternalQueries))

	costComponent := schema.CostComponent{
		Name:            "Million Custom Resolver External Queries",
		Unit:            "Million Queries",
		UnitMultiplier:  decimal.NewFromInt(1),
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
			Unit: strPtr("MILLION_ITEMS_CREXTERNALQUERIES"),
		},
	}
	return &costComponent
}

// Unit: MILLION_ITEMS (Graduated Tier)
func DNSServicesMillionDNSQueriesCostComponents(r *ResourceInstance) *schema.CostComponent {

	var million_dns_queries_included float32 = 1
	var quantity *decimal.Decimal

	additional_million_dns_queries := *r.DNSServices_DNSQueries - int64(million_dns_queries_included)
	if additional_million_dns_queries > 0 {
		quantity = decimalPtr(decimal.NewFromInt(additional_million_dns_queries))
	} else {
		quantity = decimalPtr(decimal.NewFromInt(0))
	}

	costComponent := schema.CostComponent{
		Name:            "Additional Million DNS Queries",
		Unit:            "Million Queries",
		UnitMultiplier:  decimal.NewFromInt(1),
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
			Unit: strPtr("MILLION_ITEMS"),
		},
	}
	return &costComponent
}
