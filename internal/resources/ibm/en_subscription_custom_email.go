package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// EnSubscriptionCustomEmail struct
//
// Resource information: https://cloud.ibm.com/catalog/services/event-notifications#about
// Pricing information: https://cloud.ibm.com/catalog/services/event-notifications
type EnSubscriptionCustomEmail struct {
	Address                                               string
	Region                                                string
	Name                                                  string
	Plan                                                  string
	EnSubscriptionEmail_OutboundCustomDomainEmailMessages *int64   `infracost_usage:"event-notifications_OUTBOUND_DIGITAL_MESSAGE_CUSTOM_DOMAIN_EMAIL"`
	EnSubscriptionEmail_OutboundTransmittedGB             *float64 `infracost_usage:"event-notifications_GIGABYTE_TRANSMITTED_OUTBOUND_CUSTOM_DOMAIN_EMAIL"`
}

// EnSubscriptionCustomEmailUsageSchema defines a list which represents the usage schema of EnSubscriptionCustomEmail.
var EnSubscriptionCustomEmailUsageSchema = []*schema.UsageItem{
	{Key: "event-notifications_OUTBOUND_DIGITAL_MESSAGE_CUSTOM_DOMAIN_EMAIL", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "event-notifications_GIGABYTE_TRANSMITTED_OUTBOUND_CUSTOM_DOMAIN_EMAIL", DefaultValue: 0, ValueType: schema.Int64},
}

// PopulateUsage parses the u schema.UsageData into the EnSubscriptionCustomEmail.
// It uses the `infracost_usage` struct tags to populate data into the EnSubscriptionCustomEmail.
func (r *EnSubscriptionCustomEmail) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid EnSubscriptionCustomEmail struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *EnSubscriptionCustomEmail) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		EnSubscriptionEmailOutboundCustomDomainEmailMessagesCostComponent(r),
		EnSubscriptionEmail_OutboundTransmittedGBCostComponent(r),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    EnSubscriptionCustomEmailUsageSchema,
		CostComponents: costComponents,
	}
}

func EnSubscriptionEmailOutboundCustomDomainEmailMessagesCostComponent(r *EnSubscriptionCustomEmail) *schema.CostComponent {

	var costComponent schema.CostComponent
	var quantity *decimal.Decimal

	component_name := "Outbound Custom Domain E-mail Messages"
	component_unit := "Messages"

	if r.EnSubscriptionEmail_OutboundCustomDomainEmailMessages != nil {
		quantity = decimalPtr(decimal.NewFromInt(*r.EnSubscriptionEmail_OutboundCustomDomainEmailMessages))
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
				Unit: strPtr("OUTBOUND_DIGITAL_MESSAGE_CUSTOM_DOMAIN_EMAIL"),
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

func EnSubscriptionEmail_OutboundTransmittedGBCostComponent(r *EnSubscriptionCustomEmail) *schema.CostComponent {

	var costComponent schema.CostComponent
	var quantity *decimal.Decimal

	component_name := "Outbound Transmitted E-mail Messages"
	component_unit := "GB"

	if r.EnSubscriptionEmail_OutboundTransmittedGB != nil {
		quantity = decimalPtr(decimal.NewFromFloat(*r.EnSubscriptionEmail_OutboundTransmittedGB))
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
				Unit: strPtr("GIGABYTE_TRANSMITTED_OUTBOUND_CUSTOM_DOMAIN_EMAIL"),
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
