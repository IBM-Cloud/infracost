package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// EnSubscriptionCos struct
//
// Resource information: https://cloud.ibm.com/catalog/services/event-notifications#about
// Pricing information: https://cloud.ibm.com/catalog/services/event-notifications
type EnSubscriptionCos struct {
	Address                                string
	Region                                 string
	Name                                   string
	Plan                                   string
	EnSubscriptionCOS_OutboundHTTPMessages *int64 `infracost_usage:"event-notifications_OUTBOUND_DIGITAL_MESSAGES_HTTP"`
}

// EnSubscriptionCosUsageSchema defines a list which represents the usage schema of EnSubscriptionCos.
var EnSubscriptionCosUsageSchema = []*schema.UsageItem{
	{Key: "event-notifications_OUTBOUND_DIGITAL_MESSAGES_HTTP", DefaultValue: 0, ValueType: schema.Int64},
}

// PopulateUsage parses the u schema.UsageData into the EnSubscriptionCos.
// It uses the `infracost_usage` struct tags to populate data into the EnSubscriptionCos.
func (r *EnSubscriptionCos) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid EnSubscriptionCos struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *EnSubscriptionCos) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		EnSubscriptionCOSOutboundHTTPMessagesCostComponent(r),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    EnSubscriptionCosUsageSchema,
		CostComponents: costComponents,
	}
}

func EnSubscriptionCOSOutboundHTTPMessagesCostComponent(r *EnSubscriptionCos) *schema.CostComponent {
	var costComponent schema.CostComponent
	var quantity *decimal.Decimal

	component_name := "Outbound Cloud Object Storage HTTP Messages"
	component_unit := "Messages"

	if r.EnSubscriptionCOS_OutboundHTTPMessages != nil {
		quantity = decimalPtr(decimal.NewFromInt(*r.EnSubscriptionCOS_OutboundHTTPMessages))
	}

	if r.Plan == "lite" {

		costComponent = schema.CostComponent{
			Name:            fmt.Sprintf("%s (Lite plan)", component_name),
			Unit:            component_unit,
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: quantity,
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
