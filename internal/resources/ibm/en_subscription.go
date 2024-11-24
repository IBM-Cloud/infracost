package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// EnSubscription struct
//
// Resource information: https://cloud.ibm.com/catalog/services/event-notifications#about
// Pricing information: https://cloud.ibm.com/catalog/services/event-notifications
type EnSubscription struct {
	Address                                    string
	Name                                       string
	Plan                                       string
	Region                                     string
	EnSubscriptionAndroid_OutboundPushMessages *int64 `infracost_usage:"event-notifications_OUTBOUND_DIGITAL_MESSAGES_PUSH"`
}

// EnSubscriptionUsageSchema defines a list which represents the usage schema of EnSubscription.
var EnSubscriptionUsageSchema = []*schema.UsageItem{
	{Key: "event-notifications_OUTBOUND_DIGITAL_MESSAGES_PUSH", DefaultValue: 0, ValueType: schema.Int64},
}

// PopulateUsage parses the u schema.UsageData into the EnSubscription.
// It uses the `infracost_usage` struct tags to populate data into the EnSubscription.
func (r *EnSubscription) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid EnSubscription struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *EnSubscription) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		EnSubscriptionAndroidOutboundPushMessagesCostComponent(r),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    EnSubscriptionUsageSchema,
		CostComponents: costComponents,
	}
}

func EnSubscriptionAndroidOutboundPushMessagesCostComponent(r *EnSubscription) *schema.CostComponent {
	var costComponent schema.CostComponent
	component_name := "Outbound Android Push Messages"
	component_unit := "Messages"

	if r.Plan == "lite" {
		costComponent = schema.CostComponent{
			Name:            fmt.Sprintf("%s (Lite plan)", component_name),
			Unit:            component_unit,
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

		var quantity *decimal.Decimal
		if r.EnSubscriptionAndroid_OutboundPushMessages != nil {
			quantity = decimalPtr(decimal.NewFromInt(*r.EnSubscriptionAndroid_OutboundPushMessages))
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
				Unit: strPtr("OUTBOUND_DIGITAL_MESSAGES_PUSH"),
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
