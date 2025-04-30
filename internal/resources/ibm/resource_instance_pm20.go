package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

const CUH_PER_INSTANCE = 2500

/*
 * v2-professional = "Standard" pricing plan
 * v2-standard == "Essentials" pricing plan
 * lite = "Lite" free plan
 */
func GetWMLCostComponents(r *ResourceInstance) []*schema.CostComponent {
	if r.Plan == "v2-professional" {
		return []*schema.CostComponent{
			WMLInstanceCostComponent(r),
			WMLStandardCapacityUnitHoursCostComponent(r),
			WMLClass1ResourceUnitsCostComponent(r),
			WMLClass2ResourceUnitsCostComponent(r),
			WMLClass3ResourceUnitsCostComponent(r),
		}
	} else if r.Plan == "v2-standard" {
		return []*schema.CostComponent{
			WMLEssentialsCapacityUnitHoursCostComponent(r),
			WMLClass1ResourceUnitsCostComponent(r),
			WMLClass2ResourceUnitsCostComponent(r),
			WMLClass3ResourceUnitsCostComponent(r),
		}
	} else if r.Plan == "lite" {
		costComponent := schema.CostComponent{
			Name:            "Lite plan",
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
	} else {
		costComponent := schema.CostComponent{
			Name:            fmt.Sprintf("Plan %s not found", r.Plan),
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
	}
}

func WMLInstanceCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_Instance != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_Instance))
	} else {
		q = decimalPtr(decimal.NewFromInt(1))
	}
	return &schema.CostComponent{
		Name:            "Instance (2500 CUH included)",
		Unit:            "Instance",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("INSTANCES"),
		},
	}
}

func WMLEssentialsCapacityUnitHoursCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_CUH != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_CUH))
	}
	return &schema.CostComponent{
		Name:            "Capacity Unit-Hours",
		Unit:            "CUH",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("CAPACITY_UNIT_HOURS"),
		},
	}
}

func WMLStandardCapacityUnitHoursCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	var instance float64

	if r.WML_Instance != nil {
		instance = *r.WML_Instance
	} else {
		instance = 1
	}
	if r.WML_CUH != nil {

		// standard plan is billed a fixed amount for each instance, which includes 2500 CUH's per instance.
		// if the used CUH exceeds the included quantity, the overage is charged at a flat rate.
		additional_cuh := *r.WML_CUH - (CUH_PER_INSTANCE * instance)
		if additional_cuh > 0 {
			q = decimalPtr(decimal.NewFromFloat(additional_cuh))
		}
	}

	return &schema.CostComponent{
		Name:            "Additional Capacity Unit-Hours",
		Unit:            "CUH",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("CAPACITY_UNIT_HOURS"),
		},
	}
}

func WMLStandardHoursCategoryOneCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_Hours_Category_One != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_Hours_Category_One))
	}
	return &schema.CostComponent{
		Name:            "Small Model Hosting",
		Unit:            "HOURS_CATEGORY_ONE",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("HOURS_CATEGORY_ONE"),
		},
	}
}

func WMLStandardHoursCategoryTwoCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_Hours_Category_Two != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_Hours_Category_Two))
	}
	return &schema.CostComponent{
		Name:            "Medium Model Hosting",
		Unit:            "HOURS_CATEGORY_TWO",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("HOURS_CATEGORY_TWO"),
		},
	}
}

func WMLStandardHoursCategoryThreeCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_Hours_Category_Three != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_Hours_Category_Three))
	}
	return &schema.CostComponent{
		Name:            "Large Model Hosting",
		Unit:            "HOURS_CATEGORY_THREE",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("HOURS_CATEGORY_THREE"),
		},
	}
}


func WMLStandardHoursCategoryFourCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_Hours_Category_Four != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_Hours_Category_Four))
	}
	return &schema.CostComponent{
		Name:            "Extra Large Model Hosting",
		Unit:            "HOURS_CATEGORY_FOUR",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("HOURS_CATEGORY_FOUR"),
		},
	}
}

func WMLStandardHoursCategoryFiveCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_Hours_Category_Five != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_Hours_Category_Five))
	}
	return &schema.CostComponent{
		Name:            "Very Small Model Hosting",
		Unit:            "HOURS_CATEGORY_FIVE",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("HOURS_CATEGORY_FIVE"),
		},
	}
}

func WMLStandardHoursCategorySixCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_Hours_Category_Six != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_Hours_Category_Six))
	}
	return &schema.CostComponent{
		Name:            "Very Large Model Hosting",
		Unit:            "HOURS_CATEGORY_SIX",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("HOURS_CATEGORY_SIX"),
		},
	}
}

func WMLStandardHoursMistralLargeCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_Hours_Mistral_Large != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_Hours_Mistral_Large))
	}
	return &schema.CostComponent{
		Name:            "Mistral Large Model Hosting Access",
		Unit:            "HOURS_MISTRAL_LARGE",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("HOURS_MISTRAL_LARGE"),
		},
	}
}

func WMLStandardPagesCategoryOneCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_Pages_Category_One != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_Pages_Category_One))
	}
	return &schema.CostComponent{
		Name:            "Text Extraction Category 1",
		Unit:            "PAGES_CATEGORY_ONE",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("PAGES_CATEGORY_ONE"),
		},
	}
}

func WMLStandardPagesCategoryTwoCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_Pages_Category_Two != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_Pages_Category_Two))
	}
	return &schema.CostComponent{
		Name:            "Text Extraction Category 2",
		Unit:            "PAGES_CATEGORY_TWO",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("PAGES_CATEGORY_TWO"),
		},
	}
}

func WMLStandardModelInferenceIbmCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_Model_Inference_Ibm != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_Model_Inference_Ibm))
	}
	return &schema.CostComponent{
		Name:            "Resource Units IBM Models",
		Unit:            "MODEL_INFERENCE_IBM",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("MODEL_INFERENCE_IBM"),
		},
	}
}

func WMLStandardMistralLargeInputResourceUnitsCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_Mistral_Large_Input_Resource_Units != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_Mistral_Large_Input_Resource_Units))
	}
	return &schema.CostComponent{
		Name:            "Mistral Large Input Resource Unit",
		Unit:            "MISTRAL_LARGE_INPUT_RESOURCE_UNITS",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("MISTRAL_LARGE_INPUT_RESOURCE_UNITS"),
		},
	}
}

func WMLStandardMistralLargeResourceUnitsCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_Mistral_Large_Resource_Units != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_Mistral_Large_Resource_Units))
	}
	return &schema.CostComponent{
		Name:            "Mistral Large Output Resource Unit",
		Unit:            "MISTRAL_LARGE_RESOURCE_UNITS",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("MISTRAL_LARGE_RESOURCE_UNITS"),
		},
	}
}

func WMLStandardModelInferenceThirdPartyCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_Model_Inference_Third_Party != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_Model_Inference_Third_Party))
	}
	return &schema.CostComponent{
		Name:            "Resource Units (Third Party Models)",
		Unit:            "MODEL_INFERENCE_THIRD_PARTY",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("MODEL_INFERENCE_THIRD_PARTY"),
		},
	}
}



// func WMLClass1ResourceUnitsCostComponent(r *ResourceInstance) *schema.CostComponent {
// 	var q *decimal.Decimal
// 	if r.WML_Class1RU != nil {
// 		q = decimalPtr(decimal.NewFromFloat(*r.WML_Class1RU))
// 	}
// 	return &schema.CostComponent{
// 		Name:            "Class 1 Resource Units",
// 		Unit:            "RU",
// 		UnitMultiplier:  decimal.NewFromInt(1),
// 		MonthlyQuantity: q,
// 		ProductFilter: &schema.ProductFilter{
// 			VendorName: strPtr("ibm"),
// 			Region:     strPtr(r.Location),
// 			Service:    &r.Service,
// 			AttributeFilters: []*schema.AttributeFilter{
// 				{Key: "planName", Value: &r.Plan},
// 			},
// 		},
// 		PriceFilter: &schema.PriceFilter{
// 			Unit: strPtr("CLASS_ONE_RESOURCE_UNITS"),
// 		},
// 	}
// }

// func WMLClass2ResourceUnitsCostComponent(r *ResourceInstance) *schema.CostComponent {
// 	var q *decimal.Decimal
// 	if r.WML_Class1RU != nil {
// 		q = decimalPtr(decimal.NewFromFloat(*r.WML_Class2RU))
// 	}
// 	return &schema.CostComponent{
// 		Name:            "Class 2 Resource Units",
// 		Unit:            "RU",
// 		UnitMultiplier:  decimal.NewFromInt(1),
// 		MonthlyQuantity: q,
// 		ProductFilter: &schema.ProductFilter{
// 			VendorName: strPtr("ibm"),
// 			Region:     strPtr(r.Location),
// 			Service:    &r.Service,
// 			AttributeFilters: []*schema.AttributeFilter{
// 				{Key: "planName", Value: &r.Plan},
// 			},
// 		},
// 		PriceFilter: &schema.PriceFilter{
// 			Unit: strPtr("CLASS_TWO_RESOURCE_UNITS"),
// 		},
// 	}
// }

// func WMLClass3ResourceUnitsCostComponent(r *ResourceInstance) *schema.CostComponent {
// 	var q *decimal.Decimal
// 	if r.WML_Class1RU != nil {
// 		q = decimalPtr(decimal.NewFromFloat(*r.WML_Class3RU))
// 	}
// 	return &schema.CostComponent{
// 		Name:            "Class 3 Resource Units",
// 		Unit:            "RU",
// 		UnitMultiplier:  decimal.NewFromInt(1),
// 		MonthlyQuantity: q,
// 		ProductFilter: &schema.ProductFilter{
// 			VendorName: strPtr("ibm"),
// 			Region:     strPtr(r.Location),
// 			Service:    &r.Service,
// 			AttributeFilters: []*schema.AttributeFilter{
// 				{Key: "planName", Value: &r.Plan},
// 			},
// 		},
// 		PriceFilter: &schema.PriceFilter{
// 			Unit: strPtr("CLASS_THREE_RESOURCE_UNITS"),
// 		},
// 	}
// }
