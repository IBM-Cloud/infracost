package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// EnSubscriptionSn struct
//
// Resource information: https://cloud.ibm.com/catalog/services/event-notifications#about
// Pricing information: https://cloud.ibm.com/catalog/services/event-notifications
type EnSubscriptionSn struct {
	Address                                       string
	Region                                        string
	Name                                          string
	Plan                                          string
	EnSubscriptionServiceNow_OutboundHTTPMessages *int64 `infracost_usage:"event-notifications_OUTBOUND_DIGITAL_MESSAGES_HTTP"`
}

// EnSubscriptionSnUsageSchema defines a list which represents the usage schema of EnSubscriptionSn.
var EnSubscriptionSnUsageSchema = []*schema.UsageItem{
	{Key: "event-notifications_OUTBOUND_DIGITAL_MESSAGES_HTTP", DefaultValue: 0, ValueType: schema.Int64},
}

// PopulateUsage parses the u schema.UsageData into the EnSubscriptionSn.
// It uses the `infracost_usage` struct tags to populate data into the EnSubscriptionSn.
func (r *EnSubscriptionSn) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid EnSubscriptionSn struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *EnSubscriptionSn) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		EnSubscriptionServiceNowOutboundHTTPMessagesCostComponent(r),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    EnSubscriptionSnUsageSchema,
		CostComponents: costComponents,
	}
}

func EnSubscriptionServiceNowOutboundHTTPMessagesCostComponent(r *EnSubscriptionSn) *schema.CostComponent {
	var costComponent schema.CostComponent
	component_name := "Outbound ServiceNow HTTP Messages"
	component_unit := "Messages"

	if r.Plan == "standard" {

		var quantity *decimal.Decimal
		if r.EnSubscriptionServiceNow_OutboundHTTPMessages != nil {
			quantity = decimalPtr(decimal.NewFromInt(*r.EnSubscriptionServiceNow_OutboundHTTPMessages))
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
