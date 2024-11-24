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
	Address                                    string
	Region                                     string
	Name                                       string
	Plan                                       string
	EnSubscriptionSMS_NumberResourceUnits      *float64 `infracost_usage:"event-notifications_RESOURCE_UNITS_NUMBER_MONTHLY"`
	EnSubscriptionSMS_NumberSetupResourceUnits *float64 `infracost_usage:"event-notifications_RESOURCE_UNITS_NUMBER_SETUP"`
	EnSubscriptionSMS_OutboundMessageUnits     *float64 `infracost_usage:"event-notifications_OUTBOUND_DIGITAL_MESSAGES_SMS_UNITS"`
}

// EnSubscriptionSmsUsageSchema defines a list which represents the usage schema of EnSubscriptionSms.
var EnSubscriptionSmsUsageSchema = []*schema.UsageItem{
	{Key: "event-notifications_OUTBOUND_DIGITAL_MESSAGES_SMS_UNITS", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "event-notifications_RESOURCE_UNITS_NUMBER_SETUP", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "event-notifications_RESOURCE_UNITS_NUMBER_MONTHLY", DefaultValue: 0, ValueType: schema.Float64},
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
		EnSubscriptionSMSNumberSetupResourceUnitsCostComponent(r),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    EnSubscriptionSmsUsageSchema,
		CostComponents: costComponents,
	}
}

func EnSubscriptionSMSNumberSetupResourceUnitsCostComponent(r *EnSubscriptionSms) *schema.CostComponent {

	component_unit := "Resource Units"
	component_name := "SMS Number Setup Resource Units"

	var costComponent schema.CostComponent

	if r.Plan == "lite" {

		var quantity *decimal.Decimal
		if r.EnSubscriptionSMS_NumberSetupResourceUnits != nil {
			quantity = decimalPtr(decimal.NewFromFloat(*r.EnSubscriptionSMS_NumberSetupResourceUnits))
		} else {
			quantity = decimalPtr(decimal.NewFromInt(1))
		}

		costComponent = schema.CostComponent{
			Name:            fmt.Sprintf("%s (Lite plan)", component_name),
			Unit:            component_unit,
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: quantity,
			ProductFilter: &schema.ProductFilter{
				VendorName: strPtr("ibm"),
				Region:     strPtr(r.Region),
				Service:    strPtr("event-notifications"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "planName", Value: strPtr("standard")},
				},
			},
			PriceFilter: &schema.PriceFilter{
				Unit: strPtr("RESOURCE_UNITS_NUMBER_SETUP"),
			},
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))

	} else if r.Plan == "standard" {

		var quantity *decimal.Decimal
		if r.EnSubscriptionSMS_NumberSetupResourceUnits != nil {
			quantity = decimalPtr(decimal.NewFromFloat(*r.EnSubscriptionSMS_NumberSetupResourceUnits))
		} else {
			quantity = decimalPtr(decimal.NewFromInt(1))
		}

		costComponent = schema.CostComponent{
			Name:            component_name,
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
				Unit: strPtr("RESOURCE_UNITS_NUMBER_SETUP"),
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

func EnSubscriptionSMSNumberResourceUnitsCostComponent(r *EnSubscriptionSms) *schema.CostComponent {

	var costComponent schema.CostComponent

	if r.Plan == "standard" {

		var quantity *decimal.Decimal
		if r.EnSubscriptionSMS_NumberSetupResourceUnits != nil {
			quantity = decimalPtr(decimal.NewFromFloat(*r.EnSubscriptionSMS_NumberSetupResourceUnits))
		}

		costComponent = schema.CostComponent{
			Name:            "SMS Number Use Resource Units",
			Unit:            "Resource Units",
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
				Unit: strPtr("RESOURCE_UNITS_NUMBER_MONTHLY"),
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

func EnSubscriptionSMSOutboundSMSMessageUnitsCostComponent(r *EnSubscriptionSms) *schema.CostComponent {
	var costComponent schema.CostComponent
	component_name := "Outbound IBM Cloud SMS Message Units"
	component_unit := "Message Units"

	if r.Plan == "lite" {

		quantity := math.Min(float64(*r.EnSubscriptionSMS_OutboundMessageUnits), float64(20))

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
		if r.EnSubscriptionSMS_OutboundMessageUnits != nil {
			quantity = decimalPtr(decimal.NewFromFloat(*r.EnSubscriptionSMS_OutboundMessageUnits))
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
