package ibm

import (
	"fmt"
	"math"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// EnSubscriptionSms struct
//
// Resource information: https://cloud.ibm.com/catalog/services/event-notifications#about
// Pricing information: https://cloud.ibm.com/catalog/services/event-notifications
type EnSubscriptionSms struct {
	Address                                   string
	Region                                    string
	Name                                      string
	Plan                                      string
	EnSubscriptionSMS_OutboundSMSMessageUnits *int64 `infracost_usage:"event-notifications_notifications_OUTBOUND_DIGITAL_MESSAGES_SMS_UNITS"`
}

// EnSubscriptionSmsUsageSchema defines a list which represents the usage schema of EnSubscriptionSms.
var EnSubscriptionSmsUsageSchema = []*schema.UsageItem{
	{Key: "event-notifications_notifications_OUTBOUND_DIGITAL_MESSAGES_SMS_UNITS", DefaultValue: 0, ValueType: schema.Int64},
}

// PopulateUsage parses the u schema.UsageData into the EnSubscriptionSms.
// It uses the `infracost_usage` struct tags to populate data into the EnSubscriptionSms.
func (r *EnSubscriptionSms) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid EnSubscriptionSms struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *EnSubscriptionSms) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		EnSubscriptionSMSOutboundSMSMessageUnitsCostComponent(r),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    EnSubscriptionSmsUsageSchema,
		CostComponents: costComponents,
	}
}

func EnSubscriptionSMSOutboundSMSMessageUnitsCostComponent(r *EnSubscriptionSms) *schema.CostComponent {
	var costComponent schema.CostComponent
	component_name := "Outbound IBM Cloud SMS Message Units"
	component_unit := "Message Units"

	if r.Plan == "lite" {

		quantity := math.Min(float64(*r.EnSubscriptionSMS_OutboundSMSMessageUnits), float64(20))

		costComponent = schema.CostComponent{
			Name:            fmt.Sprintf("%s (Lite plan) (Max. 20)", component_name),
			Unit:            component_unit,
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromFloat(quantity)),
			ProductFilter: &schema.ProductFilter{
				VendorName: strPtr("ibm"),
				Region:     strPtr(r.Region),
				Service:    strPtr("event-notifications"),
			},
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))

	} else if r.Plan == "standard" {

		var quantity *decimal.Decimal
		if r.EnSubscriptionSMS_OutboundSMSMessageUnits != nil {
			quantity = decimalPtr(decimal.NewFromInt(*r.EnSubscriptionSMS_OutboundSMSMessageUnits))
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
				Unit: strPtr("OUTBOUND_DIGITAL_MESSAGES_SMS_UNITS"),
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
