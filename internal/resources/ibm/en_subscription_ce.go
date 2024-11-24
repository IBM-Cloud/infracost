package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// EnSubscriptionCe struct represents <TODO: cloud service short description>.
//
// <TODO: Add any important information about the resource and links to the
// pricing pages or documentation that might be useful to developers in the future, e.g:>
//
// Resource information: https://cloud.ibm.com/<PATH/TO/RESOURCE>/
// Pricing information: https://cloud.ibm.com/<PATH/TO/PRICING>/
type EnSubscriptionCe struct {
	Address                                       string
	Region                                        string
	Name                                          string
	Plan                                          string
	EnSubscriptionCodeEngine_OutboundHTTPMessages *int64 `infracost_usage:"event-notifications_OUTBOUND_DIGITAL_MESSAGES_HTTP"`
}

// EnSubscriptionCeUsageSchema defines a list which represents the usage schema of EnSubscriptionCe.
var EnSubscriptionCeUsageSchema = []*schema.UsageItem{
	{Key: "event-notifications_OUTBOUND_DIGITAL_MESSAGES_HTTP", DefaultValue: 0, ValueType: schema.Int64},
}

// PopulateUsage parses the u schema.UsageData into the EnSubscriptionCe.
// It uses the `infracost_usage` struct tags to populate data into the EnSubscriptionCe.
func (r *EnSubscriptionCe) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid EnSubscriptionCe struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *EnSubscriptionCe) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		EnSubscriptionCodeEngineOutboundHTTPMessagesCostComponent(r),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    EnSubscriptionCeUsageSchema,
		CostComponents: costComponents,
	}
}

func EnSubscriptionCodeEngineOutboundHTTPMessagesCostComponent(r *EnSubscriptionCe) *schema.CostComponent {
	var costComponent schema.CostComponent
	component_name := "Outbound Code Engine HTTP Messages"
	component_unit := "Messages"

	if r.Plan == "standard" {

		var quantity *decimal.Decimal
		if r.EnSubscriptionCodeEngine_OutboundHTTPMessages != nil {
			quantity = decimalPtr(decimal.NewFromInt(*r.EnSubscriptionCodeEngine_OutboundHTTPMessages))
		}

		costComponent = schema.CostComponent{
			Name:            fmt.Sprintf("%s (Standard plan)", component_name),
			Unit:            component_unit,
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: quantity,
			ProductFilter: &schema.ProductFilter{
				VendorName: strPtr("ibm"),
				Region:     strPtr(r.Region),
				Service:    strPtr("event-notifications"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "planName", Value: &r.Plan},
				},
			},
			PriceFilter: &schema.PriceFilter{
				Unit: strPtr("OUTBOUND_DIGITAL_MESSAGES_HTTP"),
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
