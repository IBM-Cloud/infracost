package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

var ENTERPRISE_PLAN_PROGRAMMATIC_NAME string = "enterprise"
var PLUS_PLAN_PROGRAMMATIC_NAME string = "plus"

/*
 * Plus = 'plus'
 * Enterprise = 'enterprise'
 * Premium = 'premium' (not applicable, need to call for pricing)
 */
func GetWDCostComponents(r *ResourceInstance) []*schema.CostComponent {
	if r.Plan == PLUS_PLAN_PROGRAMMATIC_NAME {
		return []*schema.CostComponent{
			WDInstanceCostComponent(r),
			WDMonthlyDocumentsCostComponent(r),
			WDMonthlyQueriesCostComponent(r),
		}
	} else if r.Plan == ENTERPRISE_PLAN_PROGRAMMATIC_NAME {
		return []*schema.CostComponent{
			WDInstanceCostComponent(r),
			WDMonthlyDocumentsCostComponent(r),
			WDMonthlyQueriesCostComponent(r),
			WDMonthlyCustomModelsCostComponent(r),
			WDMonthlyCollectionsCostComponent(r),
		}
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

/*
 * Instance
 * - Plus: $USD/instance/month
 * - Enterprise: $USD/instance/month
 */
func WDInstanceCostComponent(r *ResourceInstance) *schema.CostComponent {

	var quantity *decimal.Decimal
	var instances_unit_name string

	if r.WD_Instance != nil {
		quantity = decimalPtr(decimal.NewFromFloat(*r.WD_Instance))
	} else {
		quantity = decimalPtr(decimal.NewFromInt(1))
	}

	if r.Plan == PLUS_PLAN_PROGRAMMATIC_NAME {
		instances_unit_name = "PLUS_SERVICE_INSTANCES_PER_MONTH"
	} else {
		instances_unit_name = "ENTERPRISE_SERVICE_INSTANCES_PER_MONTH"
	}

	return &schema.CostComponent{
		Name:            "Instance",
		Unit:            "Instance",
		UnitMultiplier:  decimal.NewFromInt(1),
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
			Unit: strPtr(instances_unit_name),
		},
	}
}

/*
 * Documents:
 * - Plus: $USD/documents/month. Includes 10,000 documents per month; $USD for every additional 1,000 documents.
 * - Enterprise: $USD/documents/month. Includes 100,000 documents per month; $USD for every additional 1,000 documents.
 */
func WDMonthlyDocumentsCostComponent(r *ResourceInstance) *schema.CostComponent {
	var quantity *decimal.Decimal

}

/*
 * Queries:
 * - Plus: $USD/queries/month. Includes 10,000 queries per month; $USD for every additional 1,000 queries.
 * - Enterprise: $USD/queries/month. Includes 100,000 queries per month; $USD for every additional 1,000 queries.
 */
func WDMonthlyQueriesCostComponent(r *ResourceInstance) *schema.CostComponent {
	var quantity *decimal.Decimal

}

/*
 * Custom Models:
 * - Enterprise: $USD/additional custom models/month. Includes 3 custom models per month.
 */
func WDMonthlyCustomModelsCostComponent(r *ResourceInstance) *schema.CostComponent {
	var quantity *decimal.Decimal

}

/*
 * Collections
 * - Enterprise: $USD/additional collections/month. Includes 300 collections per month; $USD for every additional 100 collections.
 */
func WDMonthlyCollectionsCostComponent(r *ResourceInstance) *schema.CostComponent {
	var quantity *decimal.Decimal

}
