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
			// DNSServicesNumHealthChecksCostComponents(r),
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
func DNSServicesZonesCostComponents(r *ResourceInstance) []*schema.CostComponent {
}

// Unit: NUMBERGLB (Linear Tier)
func DNSServicesPoolsPerHourCostComponents(r *ResourceInstance) []*schema.CostComponent {
}

// Unit: NUMBERPOOLS (Linear Tier)
func DNSServicesGLBInstancesPerHourCostComponents(r *ResourceInstance) []*schema.CostComponent {
}

// Unit: NUMBERHEALTHCHECK (Linear Tier)
func DNSServicesNumHealthChecksCostComponents(r *ResourceInstance) []*schema.CostComponent {
}

// Unit: RESOLVERLOCATIONS (Linear Tier)
func DNSServicesCustomResolverLocationsPerHourCostComponents(r *ResourceInstance) []*schema.CostComponent {
}

// Unit: MILLION_ITEMS_CREXTERNALQUERIES (Graduated Tier)
func DNSServicesMillionCustomResolverExternalQueriesCostComponents(r *ResourceInstance) []*schema.CostComponent {
}

// Unit: MILLION_ITEMS (Graduated Tier)
func DNSServicesMillionDNSQueriesCostComponents(r *ResourceInstance) []*schema.CostComponent {
}
