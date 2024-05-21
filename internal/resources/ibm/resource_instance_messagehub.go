package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

const EVENT_STREAMS_PROGRAMMATIC_ENTERPRISE_PLAN_NAME string = "enterprise-3nodes-2tb"
const EVENT_STREAMS_PROGRAMMATIC_LITE_PLAN_NAME string = "lite" // only available in us-south
const EVENT_STREAMS_PROGRAMMATIC_SATELLITE_PLAN_NAME string = "satellite"
const EVENT_STREAMS_PROGRAMMATIC_STANDARD_PLAN_NAME string = "standard"

func GetEventStreamsCostComponents(r *ResourceInstance) []*schema.CostComponent {

	if r.Plan == EVENT_STREAMS_PROGRAMMATIC_ENTERPRISE_PLAN_NAME {
		return []*schema.CostComponent{
			EventStreamsCapacityUnitHoursCostComponent(r),
			EventStreamsCapacityUnitHoursAdditionalCostComponent(r),
			EventStreamsTerabyteHoursCostComponent(r),
			EventStreamsGBTransmittedOutboundsCostComponent(r),
			EventStreamsCapacityUnitHoursMirroringCostComponent(r),
		}
	} else if r.Plan == EVENT_STREAMS_PROGRAMMATIC_SATELLITE_PLAN_NAME {
		return []*schema.CostComponent{
			EventStreamsCapacityUnitHoursCostComponent(r),
			EventStreamsCapacityUnitHoursAdditionalCostComponent(r),
		}
	} else if r.Plan == EVENT_STREAMS_PROGRAMMATIC_STANDARD_PLAN_NAME {
		return []*schema.CostComponent{
			EventStreamsInstanceHoursCostComponent(r),
			EventStreamsGBTransmittedOutboundsCostComponent(r),
		}
	} else if r.Plan == EVENT_STREAMS_PROGRAMMATIC_LITE_PLAN_NAME {
		costComponent := &schema.CostComponent{
			Name:            "Lite plan",
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
		return []*schema.CostComponent{costComponent}
	} else {
		costComponent := schema.CostComponent{
			Name:            fmt.Sprintf("Plan %s with customized pricing", r.Plan),
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
		return []*schema.CostComponent{
			&costComponent,
		}
	}

}

// Charged by Standard plan only
func EventStreamsInstanceHoursCostComponent(r *ResourceInstance) *schema.CostComponent {

	// var quantity *decimal.Decimal

	// if r.<VAR> != nil {
	// 	quantity = ???
	// }

	// costComponent := schema.CostComponent{
	// 	Name:            "",
	// 	Unit:            "",
	// 	UnitMultiplier:  decimal.NewFromFloat(1), // Final quantity for this cost component will be divided by this amount
	// 	MonthlyQuantity: quantity,
	// 	ProductFilter: &schema.ProductFilter{
	// 		VendorName: strPtr("ibm"),
	// 		Region:     strPtr(r.Location),
	// 		Service:    &r.Service,
	// 		AttributeFilters: []*schema.AttributeFilter{
	// 			{Key: "planName", Value: &r.Plan},
	// 		},
	// 	},
	// 	PriceFilter: &schema.PriceFilter{
	// 		Unit: strPtr("INSTANCE_HOURS"),
	// 	},
	// }
	// return &costComponent
}

func EventStreamsGBTransmittedOutboundsCostComponent(r *ResourceInstance) *schema.CostComponent {

	// var quantity *decimal.Decimal

	// if r.<VAR> != nil {
	// 	quantity = ???
	// }

	// costComponent := schema.CostComponent{
	// 	Name:            "",
	// 	Unit:            "",
	// 	UnitMultiplier:  decimal.NewFromFloat(1), // Final quantity for this cost component will be divided by this amount
	// 	MonthlyQuantity: quantity,
	// 	ProductFilter: &schema.ProductFilter{
	// 		VendorName: strPtr("ibm"),
	// 		Region:     strPtr(r.Location),
	// 		Service:    &r.Service,
	// 		AttributeFilters: []*schema.AttributeFilter{
	// 			{Key: "planName", Value: &r.Plan},
	// 		},
	// 	},
	// 	PriceFilter: &schema.PriceFilter{
	// 		Unit: strPtr("GIGABYTE_TRANSMITTED_OUTBOUNDS"),
	// 	},
	// }
	// return &costComponent
}

func EventStreamsCapacityUnitHoursCostComponent(r *ResourceInstance) *schema.CostComponent {

	// var quantity *decimal.Decimal

	// if r.<VAR> != nil {
	// 	quantity = ???
	// }

	// costComponent := schema.CostComponent{
	// 	Name:            "",
	// 	Unit:            "",
	// 	UnitMultiplier:  decimal.NewFromFloat(1), // Final quantity for this cost component will be divided by this amount
	// 	MonthlyQuantity: quantity,
	// 	ProductFilter: &schema.ProductFilter{
	// 		VendorName: strPtr("ibm"),
	// 		Region:     strPtr(r.Location),
	// 		Service:    &r.Service,
	// 		AttributeFilters: []*schema.AttributeFilter{
	// 			{Key: "planName", Value: &r.Plan},
	// 		},
	// 	},
	// 	PriceFilter: &schema.PriceFilter{
	// 		Unit: strPtr("CAPACITY_UNIT_HOURS"),
	// 	},
	// }
	// return &costComponent
}

func EventStreamsCapacityUnitHoursAdditionalCostComponent(r *ResourceInstance) *schema.CostComponent {

	// var quantity *decimal.Decimal

	// if r.<VAR> != nil {
	// 	quantity = ???
	// }

	// costComponent := schema.CostComponent{
	// 	Name:            "",
	// 	Unit:            "",
	// 	UnitMultiplier:  decimal.NewFromFloat(1), // Final quantity for this cost component will be divided by this amount
	// 	MonthlyQuantity: quantity,
	// 	ProductFilter: &schema.ProductFilter{
	// 		VendorName: strPtr("ibm"),
	// 		Region:     strPtr(r.Location),
	// 		Service:    &r.Service,
	// 		AttributeFilters: []*schema.AttributeFilter{
	// 			{Key: "planName", Value: &r.Plan},
	// 		},
	// 	},
	// 	PriceFilter: &schema.PriceFilter{
	// 		Unit: strPtr("CAPACITY_UNIT_HOURS_ADDITIONAL"),
	// 	},
	// }
	// return &costComponent
}

func EventStreamsTerabyteHoursCostComponent(r *ResourceInstance) *schema.CostComponent {

	// var quantity *decimal.Decimal

	// if r.<VAR> != nil {
	// 	quantity = ???
	// }

	// costComponent := schema.CostComponent{
	// 	Name:            "",
	// 	Unit:            "",
	// 	UnitMultiplier:  decimal.NewFromFloat(1), // Final quantity for this cost component will be divided by this amount
	// 	MonthlyQuantity: quantity,
	// 	ProductFilter: &schema.ProductFilter{
	// 		VendorName: strPtr("ibm"),
	// 		Region:     strPtr(r.Location),
	// 		Service:    &r.Service,
	// 		AttributeFilters: []*schema.AttributeFilter{
	// 			{Key: "planName", Value: &r.Plan},
	// 		},
	// 	},
	// 	PriceFilter: &schema.PriceFilter{
	// 		Unit: strPtr("TERABYTE_HOURS"),
	// 	},
	// }
	// return &costComponent
}

func EventStreamsCapacityUnitHoursMirroringCostComponent(r *ResourceInstance) *schema.CostComponent {

	// var quantity *decimal.Decimal

	// if r.<VAR> != nil {
	// 	quantity = ???
	// }

	// costComponent := schema.CostComponent{
	// 	Name:            "",
	// 	Unit:            "",
	// 	UnitMultiplier:  decimal.NewFromFloat(1), // Final quantity for this cost component will be divided by this amount
	// 	MonthlyQuantity: quantity,
	// 	ProductFilter: &schema.ProductFilter{
	// 		VendorName: strPtr("ibm"),
	// 		Region:     strPtr(r.Location),
	// 		Service:    &r.Service,
	// 		AttributeFilters: []*schema.AttributeFilter{
	// 			{Key: "planName", Value: &r.Plan},
	// 		},
	// 	},
	// 	PriceFilter: &schema.PriceFilter{
	// 		Unit: strPtr("CAPACITY_UNIT_HOURS_MIRRORING"),
	// 	},
	// }
	// return &costComponent
}
