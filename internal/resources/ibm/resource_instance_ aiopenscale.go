package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

/**
* Lite: 'lite' (Free)
* Essentials: 'essentials' ($0.60 USD/RU)
 */
func GetWGOVCostComponents(r *ResourceInstance) []*schema.CostComponent {
	if r.Plan == "essentials" {
		// TODO
		// Note: Global Catalog page only has one metric for "RESOURCE_UNITS"; it does not differentiate between types of models, evaluations, etc. May need to re-think the variables that have been used in the usage file.
	} else if r.Plan == "lite" {
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

/**
* 1 RU ($0.60 USD) for every:
* 1 WGOV_PredictiveModelEvals
* 1 WGOV_FoundationalModelEvals
* 1 WGOV_GlobalExplanations
* 500 WGOV_LocalExplanations
 */
func ResourceUnitCostComponent(r *ResourceInstance) *schema.CostComponent {
	var quantity *decimal.Decimal
	// TODO
}
