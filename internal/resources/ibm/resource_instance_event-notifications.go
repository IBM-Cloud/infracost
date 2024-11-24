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
			Unit:            "Instance",
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
