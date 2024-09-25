package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// IsInstance struct represents an IBM virtual server instance.
//
// Pricing information: https://cloud.ibm.com/kubernetes/catalog/about

type IsInstance struct {
	Address     string
	Region      string
	Profile     string // should be values from CLI 'ibmcloud is instance-profiles'
	Zone        string
	IsDedicated bool // will be true if a dedicated_host or dedicated_host_group is specified
	BootVolume  []struct {
		Name string
		Size int64
	}
	MonthlyInstanceHours *float64 `infracost_usage:"monthly_instance_hours"`
}

var IsInstanceUsageSchema = []*schema.UsageItem{
	{Key: "monthly_instance_hours", DefaultValue: 0, ValueType: schema.Float64},
}

// PopulateUsage parses the u schema.UsageData into the IsInstance.
// It uses the `infracost_usage` struct tags to populate data into the IsInstance.
func (r *IsInstance) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *IsInstance) instanceHoursCostComponent() *schema.CostComponent {

	unit := "INSTANCE_HOURS_MULTI_TENANT"
	if r.IsDedicated {
		unit = "INSTANCE_HOURS_DEDICATED_HOST"
	}

	return &schema.CostComponent{
		Name:            fmt.Sprintf("Instance Hours (%s)", r.Profile),
		Unit:            "Hours",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: decimalPtr(decimal.NewFromFloat(*r.MonthlyInstanceHours)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			Service:       strPtr("is.instance"),
			ProductFamily: strPtr("service"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Profile},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}

// BuildResource builds a schema.Resource from a valid IsShare struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *IsInstance) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		r.instanceHoursCostComponent(),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    IsInstanceUsageSchema,
		CostComponents: costComponents,
	}
}
