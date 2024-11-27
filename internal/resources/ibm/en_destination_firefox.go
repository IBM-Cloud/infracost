package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// EnDestinationFirefox struct
//
// Resource information: https://cloud.ibm.com/catalog/services/event-notifications#about
// Pricing information: https://cloud.ibm.com/catalog/services/event-notifications
type EnDestinationFirefox struct {
	Address   string
	IsPreProd bool
	Name      string
	Plan      string
	Region    string
}

// EnDestinationFirefoxUsageSchema defines a list which represents the usage schema of EnDestinationFirefox.
var EnDestinationFirefoxUsageSchema = []*schema.UsageItem{}

// PopulateUsage parses the u schema.UsageData into the EnDestinationFirefox.
// It uses the `infracost_usage` struct tags to populate data into the EnDestinationFirefox.
func (r *EnDestinationFirefox) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid EnDestinationFirefox struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *EnDestinationFirefox) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		EnDestinationFirefoxPushDestinationInstancesCostComponent(r),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    EnDestinationFirefoxUsageSchema,
		CostComponents: costComponents,
	}
}

func EnDestinationFirefoxPushDestinationInstancesCostComponent(r *EnDestinationFirefox) *schema.CostComponent {

	var costComponent schema.CostComponent
	component_name := "Push Destination Instances"
	unit := "PUSH_DESTINATION_INSTANCES"

	if r.IsPreProd {
		component_name = "Pre-Prod Push Destination Instances"
		unit = "PUSH_PREPROD_DESTINATION_INSTANCES"
	}

	if r.Plan == "lite" {

		costComponent = schema.CostComponent{
			Name:            fmt.Sprintf("%s (Lite plan)", component_name),
			Unit:            "Instance",
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
			ProductFilter: &schema.ProductFilter{
				VendorName: strPtr("ibm"),
				Region:     strPtr(r.Region),
				Service:    strPtr("event-notifications"),
			},
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
	} else if r.Plan == "standard" {

		costComponent = schema.CostComponent{
			Name:            fmt.Sprintf("%s (Standard plan)", component_name),
			Unit:            "Instance",
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
			ProductFilter: &schema.ProductFilter{
				VendorName: strPtr("ibm"),
				Region:     strPtr(r.Region),
				Service:    strPtr("event-notifications"),
				AttributeFilters: []*schema.AttributeFilter{ // Only standard plan exists
					{Key: "planName", Value: &r.Plan},
				},
			},
			PriceFilter: &schema.PriceFilter{
				Unit: strPtr(unit),
			},
		}

	} else {
		costComponent = schema.CostComponent{
			Name:            fmt.Sprintf("Plan %s not found", r.Plan),
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
			ProductFilter: &schema.ProductFilter{
				VendorName: strPtr("ibm"),
				Region:     strPtr(r.Region),
				Service:    strPtr("event-notifications"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "planName", Value: &r.Plan},
				},
			},
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
	}
	return &costComponent
}
