package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// IsFloatingIp struct represents <TODO: cloud service short description>.
//
// <TODO: Add any important information about the resource and links to the
// pricing pages or documentation that might be useful to developers in the future, e.g:>
//
// Resource information: https://cloud.ibm.com/<PATH/TO/RESOURCE>/
// Pricing information: https://cloud.ibm.com/<PATH/TO/PRICING>/
type IsFloatingIp struct {
	Address string
	Region  string
}

// IsFloatingIpUsageSchema defines a list which represents the usage schema of IsFloatingIp.
var IsFloatingIpUsageSchema = []*schema.UsageItem{}

// PopulateUsage parses the u schema.UsageData into the IsFloatingIp.
// It uses the `infracost_usage` struct tags to populate data into the IsFloatingIp.
func (r *IsFloatingIp) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *IsFloatingIp) isFloatingIpCostComponent() *schema.CostComponent {
	return &schema.CostComponent{
		Name:            fmt.Sprintf("Floating Ip %s", r.Region),
		Unit:            "Instance",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			Service:       strPtr("is.floating-ip"),
			ProductFamily: strPtr("service"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("default-nextgen")},
			},
		},
	}
}

// BuildResource builds a schema.Resource from a valid IsFloatingIp struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *IsFloatingIp) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		r.isFloatingIpCostComponent(),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    IsFloatingIpUsageSchema,
		CostComponents: costComponents,
	}
}
