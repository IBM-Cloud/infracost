package ibm

import (
	"fmt"
	"strconv"

	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// Plans: Lite, Basic, Standard, Enterprise
const AppRappLitePlanProgrammaticName string = "lite"
const AppRappBasicPlanProgrammaticName string = "basic"
const AppRappStandardPlanProgrammaticName string = "standardv2"
const AppRappEnterprisePlanProgrammaticName string = "enterprise"

func GetAppRappCostComponents(r *ResourceInstance) []*schema.CostComponent {
	switch r.Plan {
	case AppRappBasicPlanProgrammaticName:
		return []*schema.CostComponent{
			AppRappActiveEntityIDCostComponent(r),
			AppRappHundredThousandAPICallsCostComponent(r),
		}
	case AppRappStandardPlanProgrammaticName:
		return []*schema.CostComponent{
			AppRappInstanceCostComponent(r),
			AppRappActiveEntityIDCostComponent(r),
			AppRappHundredThousandAPICallsCostComponent(r),
		}
	case AppRappEnterprisePlanProgrammaticName:
		return []*schema.CostComponent{
			AppRappInstanceCostComponent(r),
			AppRappActiveEntityIDCostComponent(r),
			AppRappHundredThousandAPICallsCostComponent(r),
		}
	default:
		// Plan not found — could be Lite — set to $0

		var componentName string
		if r.Plan == AppRappLitePlanProgrammaticName {
			componentName = "Lite plan"
		} else {
			componentName = fmt.Sprintf("Plan %s not found", r.Plan)
		}

		costComponent := schema.CostComponent{
			Name:            componentName,
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

// Active Entity ID cost is based on the number of unique entities that interact with App Configuration instance during the month
func AppRappActiveEntityIDCostComponent(r *ResourceInstance) *schema.CostComponent {
	var ActiveEntityIDsUnitName string = "ACTIVE_ENTITY_IDS"
	var quantity *decimal.Decimal
	var activeEntityIDsIncluded int = 0
	var costComponentName string = "Active Entity IDs"
	var costComponentunitName string = "Active Entity IDs"

	if r.Plan == AppRappStandardPlanProgrammaticName || r.Plan == AppRappEnterprisePlanProgrammaticName { // For Standard, the monthly instance price includes 100,000 API calls
		switch r.Plan {
		case AppRappStandardPlanProgrammaticName: // For Standard, the monthly instance price includes 1000 active entity IDs
			activeEntityIDsIncluded = 1000
		case AppRappEnterprisePlanProgrammaticName: // For Enterprise, the monthly instance price includes 10,000 active entity IDs
			activeEntityIDsIncluded = 10000
		}
		costComponentName = "Additional " + costComponentName + " (first " + strconv.Itoa(activeEntityIDsIncluded) + " included)"
	}

	if r.AppRapp_Active_Entity_IDs != nil {
		quantity = decimalPtr(decimal.NewFromInt(max(*r.AppRapp_Active_Entity_IDs-int64(activeEntityIDsIncluded), 0)))
	} else {
		quantity = decimalPtr(decimal.NewFromInt(0))
	}

	if quantity.IntPart() != 1 {
		costComponentunitName += "s"
	}

	return &schema.CostComponent{
		Name:            costComponentName,
		Unit:            costComponentunitName,
		UnitMultiplier:  decimal.NewFromInt(1), // Final quantity for this cost component will be divided by this amount
		MonthlyQuantity: quantity,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(ActiveEntityIDsUnitName),
		},
	}
}

func AppRappHundredThousandAPICallsCostComponent(r *ResourceInstance) *schema.CostComponent {
	var HundredThousandApiCallsUnitName string = "HUNDRED_THOUSAND_API_CALLS"
	var quantity *decimal.Decimal
	var apiCallsIncluded int = 0
	var costComponentName string = "API Calls"
	var costComponentUnitName string = "100k API Calls"

	if r.Plan == AppRappStandardPlanProgrammaticName || r.Plan == AppRappEnterprisePlanProgrammaticName { // For Standard, the monthly instance price includes 100,000 API calls
		switch r.Plan {
		case AppRappStandardPlanProgrammaticName: // For Standard, the monthly instance price includes 100,000 API calls
			apiCallsIncluded = 100000
		case AppRappEnterprisePlanProgrammaticName: // For Enterprise, the monthly instance price includes 1,000,000 API calls
			apiCallsIncluded = 1000000
		}
		costComponentName = "Additional " + costComponentName + " (first " + strconv.Itoa(apiCallsIncluded) + " included)"
	}

	if r.AppRapp_API_Calls != nil {
		quantity = decimalPtr(decimal.NewFromFloat(max(float64(*r.AppRapp_API_Calls-int64(apiCallsIncluded))/100000, 0)))
	} else {
		quantity = decimalPtr(decimal.NewFromInt(0))
	}

	return &schema.CostComponent{
		Name:            costComponentName,
		Unit:            costComponentUnitName,
		UnitMultiplier:  decimal.NewFromInt(1), // Final quantity for this cost component will be divided by this amount
		MonthlyQuantity: quantity,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(HundredThousandApiCallsUnitName),
		},
	}
}

func AppRappInstanceCostComponent(r *ResourceInstance) *schema.CostComponent {
	var instancesUnitName string = "APPLICATION_INSTANCES"
	var quantity *decimal.Decimal = decimalPtr(decimal.NewFromInt(1))

	return &schema.CostComponent{
		Name:            "Instance",
		Unit:            "Instance",
		UnitMultiplier:  decimal.NewFromInt(1), // Final quantity for this cost component will be divided by this amount
		MonthlyQuantity: quantity,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(instancesUnitName),
		},
	}
}
