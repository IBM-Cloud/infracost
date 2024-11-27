package ibm

import (
	"fmt"
	"math"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// EnSubscriptionEmail struct
//
// Resource information: https://cloud.ibm.com/catalog/services/event-notifications#about
// Pricing information: https://cloud.ibm.com/catalog/services/event-notifications
type EnSubscriptionEmail struct {
	Address                                   string
	Region                                    string
	Name                                      string
	Plan                                      string
	EnSubscriptionEmail_OutboundEmailMessages *int64 `infracost_usage:"event-notifications_OUTBOUND_DIGITAL_MESSAGES_EMAILS"`
}

// EnSubscriptionEmailUsageSchema defines a list which represents the usage schema of EnSubscriptionEmail.
var EnSubscriptionEmailUsageSchema = []*schema.UsageItem{
	{Key: "event-notifications_OUTBOUND_DIGITAL_MESSAGES_EMAILS", DefaultValue: 0, ValueType: schema.Int64},
}

// PopulateUsage parses the u schema.UsageData into the EnSubscriptionEmail.
// It uses the `infracost_usage` struct tags to populate data into the EnSubscriptionEmail.
func (r *EnSubscriptionEmail) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid EnSubscriptionEmail struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *EnSubscriptionEmail) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		EnSubscriptionEmailOutboundMessagesCostComponent(r),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    EnSubscriptionEmailUsageSchema,
		CostComponents: costComponents,
	}
}

func EnSubscriptionEmailOutboundMessagesCostComponent(r *EnSubscriptionEmail) *schema.CostComponent {

	var costComponent schema.CostComponent

	component_name := "Outbound E-mail Messages"
	component_unit := "Messages"

	if r.Plan == "lite" {

		quantity := decimalPtr(decimal.NewFromFloat(math.Min(float64(*r.EnSubscriptionEmail_OutboundEmailMessages), float64(20))))

		costComponent = schema.CostComponent{
			Name:            fmt.Sprintf("%s (Lite plan) (Max. 20)", component_name),
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

		var quantity *decimal.Decimal
		if r.EnSubscriptionEmail_OutboundEmailMessages != nil {
			quantity = decimalPtr(decimal.NewFromInt(*r.EnSubscriptionEmail_OutboundEmailMessages))
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
				Unit: strPtr("OUTBOUND_DIGITAL_MESSAGES_EMAILS"),
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
