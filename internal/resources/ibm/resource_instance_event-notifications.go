package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

const EVENT_NOTIFICATIONS_LITE_PLAN_PROGRAMMATIC_NAME = "lite"
const EVENT_NOTIFICATIONS_STANDARD_PLAN_PROGRAMMATIC_NAME = "standard"

func GetEventNotificationsCostComponents(r *ResourceInstance) []*schema.CostComponent {
	if r.Plan == EVENT_NOTIFICATIONS_LITE_PLAN_PROGRAMMATIC_NAME {
		costComponent := schema.CostComponent{
			Name:            "Lite Plan",
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
			ProductFilter: &schema.ProductFilter{
				VendorName: strPtr("ibm"),
				Region:     strPtr(r.Location),
				Service:    &r.Service,
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "planName", Value: &r.Plan},
				},
			},
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
		return []*schema.CostComponent{
			&costComponent,
		}
	} else if r.Plan == EVENT_NOTIFICATIONS_STANDARD_PLAN_PROGRAMMATIC_NAME {
		return []*schema.CostComponent{
			EventNotificationsInboundIngestedEventsCostComponent(r),
			// EventNotificationsOutboundCustomDomainEmailGBsTransmittedCostComponent(r),
			// EventNotificationsOutboundCustomDomainEmailsCostComponent(r),
			// EventNotificationsOutboundEmailsCostComponent(r),
			// EventNotificationsOutboundHTTPMessagesCostComponent(r),
			// EventNotificationsOutboundPushMessagesCostComponent(r),
			// EventNotificationsOutboundSMSMessagesCostComponent(r),
			// EventNotificationsResourceUnitsMonthlyCostComponent(r),
			// EventNotificationsResourceUnitsSetupCostComponent(r),
		}
	} else {
		costComponent := schema.CostComponent{
			Name:            fmt.Sprintf("Plan %s not found", r.Plan),
			UnitMultiplier:  decimal.NewFromInt(1), // Final quantity for this cost component will be divided by this amount
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
			ProductFilter: &schema.ProductFilter{
				VendorName: strPtr("ibm"),
				Region:     strPtr(r.Location),
				Service:    &r.Service,
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "planName", Value: &r.Plan},
				},
			},
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
		return []*schema.CostComponent{
			&costComponent,
		}
	}
}

func EventNotificationsInboundIngestedEventsCostComponent(r *ResourceInstance) *schema.CostComponent {

	var quantity *decimal.Decimal
	if r.EventNotifications_InboundIngestedEvents != nil {
		quantity = decimalPtr(decimal.NewFromFloat(*r.EventNotifications_InboundIngestedEvents))
	}

	costComponent := schema.CostComponent{
		Name:            "Ingested Events",     // Short descriptive name of the component.
		Unit:            "Million Events",      // Unit of resource component's measurement. For example, it can be hours or 10M requests.
		UnitMultiplier:  decimal.NewFromInt(1), // Used to calculate the cost of component quantity correctly. For example, if a price is $0.02 per 1k requests, and assuming the amount is 10,000, its cost will be calculated as quantity/unitMultiplier * price.
		MonthlyQuantity: quantity,              // HourlyQuantity or MonthlyQuantity attributes specify the quantity of the resource. If the measurement unit is GB, it will be the number of gigabytes. If the unit is hours, it can be 1 as "1 hour"
		ProductFilter: &schema.ProductFilter{ // Helps identify the exact price of the "product." Usually, it's only one, but if there are pricing tiers, its filters can pick the correct value.
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("MILLION_INGESTED_EVENTS"),
		},
	}
	return &costComponent
}

// func EventNotificationsOutboundCustomDomainEmailGBsTransmittedCostComponent(r *ResourceInstance) *schema.CostComponent {

// 	var quantity *decimal.Decimal
// 	if r.EventNotifications_OutboundCustomDomainEmailGBsTransmitted != nil {
// 		quantity = decimalPtr(decimal.NewFromFloat(*r.EventNotifications_OutboundCustomDomainEmailGBsTransmitted))
// 	}

// 	costComponent := schema.CostComponent{
// 		Name:            "Outbound Custom Domain E-mail GB Transmitted",
// 		Unit:            "GB",
// 		UnitMultiplier:  decimal.NewFromInt(1),
// 		MonthlyQuantity: quantity,
// 		ProductFilter: &schema.ProductFilter{
// 			VendorName: strPtr("ibm"),
// 			Region:     strPtr(r.Location),
// 			Service:    &r.Service,
// 			AttributeFilters: []*schema.AttributeFilter{
// 				{Key: "planName", Value: &r.Plan},
// 			},
// 		},
// 		PriceFilter: &schema.PriceFilter{
// 			Unit: strPtr("GIGABYTE_TRANSMITTED_OUTBOUND_CUSTOM_DOMAIN_EMAIL"),
// 		},
// 	}
// 	return &costComponent
// }

// func EventNotificationsOutboundCustomDomainEmailsCostComponent(r *ResourceInstance) *schema.CostComponent {

// 	var quantity *decimal.Decimal
// 	if r.EventNotifications_OutboundCustomDomainEmails != nil {
// 		quantity = decimalPtr(decimal.NewFromInt(*r.EventNotifications_OutboundCustomDomainEmails))
// 	}

// 	costComponent := schema.CostComponent{
// 		Name:            "Outbound Custom Domain E-mails",
// 		Unit:            "E-mails",
// 		UnitMultiplier:  decimal.NewFromInt(1),
// 		MonthlyQuantity: quantity,
// 		ProductFilter: &schema.ProductFilter{
// 			VendorName: strPtr("ibm"),
// 			Region:     strPtr(r.Location),
// 			Service:    &r.Service,
// 			AttributeFilters: []*schema.AttributeFilter{
// 				{Key: "planName", Value: &r.Plan},
// 			},
// 		},
// 		PriceFilter: &schema.PriceFilter{
// 			Unit: strPtr("OUTBOUND_DIGITAL_MESSAGE_CUSTOM_DOMAIN_EMAIL"),
// 		},
// 	}
// 	return &costComponent
// }

// func EventNotificationsOutboundEmailsCostComponent(r *ResourceInstance) *schema.CostComponent {

// 	var quantity *decimal.Decimal
// 	if r.EventNotifications_OutboundEmails != nil {
// 		quantity = decimalPtr(decimal.NewFromInt(*r.EventNotifications_OutboundEmails))
// 	}

// 	costComponent := schema.CostComponent{
// 		Name:            "Outbound E-mails",
// 		Unit:            "E-mails",
// 		UnitMultiplier:  decimal.NewFromInt(1),
// 		MonthlyQuantity: quantity,
// 		ProductFilter: &schema.ProductFilter{
// 			VendorName: strPtr("ibm"),
// 			Region:     strPtr(r.Location),
// 			Service:    &r.Service,
// 			AttributeFilters: []*schema.AttributeFilter{
// 				{Key: "planName", Value: &r.Plan},
// 			},
// 		},
// 		PriceFilter: &schema.PriceFilter{
// 			Unit: strPtr("OUTBOUND_DIGITAL_MESSAGES_EMAILS"),
// 		},
// 	}
// 	return &costComponent
// }

// func EventNotificationsOutboundHTTPMessagesCostComponent(r *ResourceInstance) *schema.CostComponent {

// 	var quantity *decimal.Decimal
// 	if r.EventNotifications_OutboundHTTPMessages != nil {
// 		quantity = decimalPtr(decimal.NewFromInt(*r.EventNotifications_OutboundHTTPMessages))
// 	}

// 	costComponent := schema.CostComponent{
// 		Name:            "Outbound HTTP Messages",
// 		Unit:            "Messages",
// 		UnitMultiplier:  decimal.NewFromInt(1),
// 		MonthlyQuantity: quantity,
// 		ProductFilter: &schema.ProductFilter{
// 			VendorName: strPtr("ibm"),
// 			Region:     strPtr(r.Location),
// 			Service:    &r.Service,
// 			AttributeFilters: []*schema.AttributeFilter{
// 				{Key: "planName", Value: &r.Plan},
// 			},
// 		},
// 		PriceFilter: &schema.PriceFilter{
// 			Unit: strPtr("OUTBOUND_DIGITAL_MESSAGES_HTTP"),
// 		},
// 	}
// 	return &costComponent
// }

// func EventNotificationsOutboundPushMessagesCostComponent(r *ResourceInstance) *schema.CostComponent {

// 	var quantity *decimal.Decimal
// 	if r.EventNotifications_OutboundPushMessages != nil {
// 		quantity = decimalPtr(decimal.NewFromInt(*r.EventNotifications_OutboundPushMessages))
// 	}

// 	costComponent := schema.CostComponent{
// 		Name:            "Outbound Push Messages",
// 		Unit:            "Messages",
// 		UnitMultiplier:  decimal.NewFromInt(1),
// 		MonthlyQuantity: quantity,
// 		ProductFilter: &schema.ProductFilter{
// 			VendorName: strPtr("ibm"),
// 			Region:     strPtr(r.Location),
// 			Service:    &r.Service,
// 			AttributeFilters: []*schema.AttributeFilter{
// 				{Key: "planName", Value: &r.Plan},
// 			},
// 		},
// 		PriceFilter: &schema.PriceFilter{
// 			Unit: strPtr("OUTBOUND_DIGITAL_MESSAGES_PUSH"),
// 		},
// 	}
// 	return &costComponent
// }

// func EventNotificationsOutboundSMSMessagesCostComponent(r *ResourceInstance) *schema.CostComponent {

// 	var quantity *decimal.Decimal
// 	if r.EventNotifications_OutboundSMSMessages != nil {
// 		quantity = decimalPtr(decimal.NewFromInt(*r.EventNotifications_OutboundSMSMessages))
// 	}

// 	costComponent := schema.CostComponent{
// 		Name:            "Outbound SMS Messages",
// 		Unit:            "Messages",
// 		UnitMultiplier:  decimal.NewFromInt(1),
// 		MonthlyQuantity: quantity,
// 		ProductFilter: &schema.ProductFilter{
// 			VendorName: strPtr("ibm"),
// 			Region:     strPtr(r.Location),
// 			Service:    &r.Service,
// 			AttributeFilters: []*schema.AttributeFilter{
// 				{Key: "planName", Value: &r.Plan},
// 			},
// 		},
// 		PriceFilter: &schema.PriceFilter{
// 			Unit: strPtr("OUTBOUND_DIGITAL_MESSAGES_SMS_UNITS"),
// 		},
// 	}
// 	return &costComponent
// }

// /*
//  * D0HNPZX - IBM Cloud Event Notifications SMS number monthly fee unit Resource Unit Pay per Use
//  * This part covers the monthly usage fee for a phone number. It recurs monthly. The charge metric is 'resource unit'. Since every country and number type is a different price, the Event Notifications team will send over the proper number of units for the given use case. One unit = $1.
//  */
// func EventNotificationsResourceUnitsMonthlyCostComponent(r *ResourceInstance) *schema.CostComponent {

// 	var quantity *decimal.Decimal
// 	if r.EventNotifications_ResourceUnitsMonthly != nil {
// 		quantity = decimalPtr(decimal.NewFromFloat(*r.EventNotifications_ResourceUnitsMonthly))
// 	}

// 	costComponent := schema.CostComponent{
// 		Name:            "SMS Number Monthly Resource Units",
// 		Unit:            "Resource Units",
// 		UnitMultiplier:  decimal.NewFromInt(1),
// 		MonthlyQuantity: quantity,
// 		ProductFilter: &schema.ProductFilter{
// 			VendorName: strPtr("ibm"),
// 			Region:     strPtr(r.Location),
// 			Service:    &r.Service,
// 			AttributeFilters: []*schema.AttributeFilter{
// 				{Key: "planName", Value: &r.Plan},
// 			},
// 		},
// 		PriceFilter: &schema.PriceFilter{
// 			Unit: strPtr("RESOURCE_UNITS_NUMBER_MONTHLY"),
// 		},
// 	}
// 	return &costComponent
// }

// /*
//  * D0HNQZX - IBM Cloud Event Notifications SMS number setup fee unit Resource Unit Pay per Use.
//  * This part covers any setup fee a phone number might have. It will be a one-time fee. The charge metric is 'resource unit'. Since every country and number type is a different price, the Event Notifications team will send over the proper number of units for the given use case. One unit = $1.
//  */
// func EventNotificationsResourceUnitsSetupCostComponent(r *ResourceInstance) *schema.CostComponent {

// 	var quantity *decimal.Decimal
// 	if r.EventNotifications_ResourceUnitsSetup != nil {
// 		quantity = decimalPtr(decimal.NewFromFloat(*r.EventNotifications_ResourceUnitsSetup))
// 	}

// 	costComponent := schema.CostComponent{
// 		Name:            "SMS Number Setup Resource Units",
// 		Unit:            "Resource Units",
// 		UnitMultiplier:  decimal.NewFromInt(1),
// 		MonthlyQuantity: quantity,
// 		ProductFilter: &schema.ProductFilter{
// 			VendorName: strPtr("ibm"),
// 			Region:     strPtr(r.Location),
// 			Service:    &r.Service,
// 			AttributeFilters: []*schema.AttributeFilter{
// 				{Key: "planName", Value: &r.Plan},
// 			},
// 		},
// 		PriceFilter: &schema.PriceFilter{
// 			Unit: strPtr("RESOURCE_UNITS_NUMBER_SETUP"),
// 		},
// 	}
// 	return &costComponent
// }
