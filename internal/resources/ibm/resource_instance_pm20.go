package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

const CUH_PER_INSTANCE = 2500

// Map used to generate cost components based on model hosting offerings by a specific GPU/hour.
// The keys are the unit values defined in our PostGres DB. The values are the title to be used in the golden file.
var gpuMap = map[string]string{
	// Categorical GPU selection
	"HOURS_MISTRAL_LARGE":   "Mistral Large Model Hosting Access",
	"HOURS_MISTRAL_ONE_GPU": "Mistral 1 GPU Model Hosting Access",
	"HOURS_MISTRAL_TWO_GPU": "Mistral 2 GPU Model Hosting Access",
	"HOURS_CATEGORY_ONE":    "Small Model Hosting",
	"HOURS_CATEGORY_TWO":    "Medium Model Hosting",
	"HOURS_CATEGORY_THREE":  "Large Model Hosting",
	"HOURS_CATEGORY_FOUR":   "Extra Large Model Hosting",
	"HOURS_CATEGORY_FIVE":   "Extra Small Model Hosting",
	"HOURS_CATEGORY_SIX":    "Very Large Model Hosting",
	// Specific GPU selection
	"HOURS_ONE_L_FORTY_S":       "Model Hosting 1 L40S",
	"HOURS_TWO_L_FORTY_S":       "Model Hosting 2 L40S",
	"HOURS_ONE_A_ONE_HUNDRED":   "Model Hosting 1 A100",
	"HOURS_TWO_A_ONE_HUNDRED":   "Model Hosting 2 A100",
	"HOURS_FOUR_A_ONE_HUNDRED":  "Model Hosting 4 A100",
	"HOURS_EIGHT_A_ONE_HUNDRED": "Model Hosting 8 A100",
	"HOURS_ONE_H_ONE_HUNDRED":   "Model Hosting 1 H100",
	"HOURS_TWO_H_ONE_HUNDRED":   "Model Hosting 2 H100",
	"HOURS_FOUR_H_ONE_HUNDRED":  "Model Hosting 4 H100",
	"HOURS_EIGHT_H_ONE_HUNDRED": "Model Hosting 8 H100",
	"HOURS_ONE_H_TWO_HUNDRED":   "Model Hosting 1 H200",
	"HOURS_TWO_H_TWO_HUNDRED":   "Model Hosting 2 H200",
	"HOURS_FOUR_H_TWO_HUNDRED":  "Model Hosting 4 H200",
	"HOURS_EIGHT_H_TWO_HUNDRED": "Model Hosting 8 H200",
}

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
			WMLMistralLargeOutputResourceUnitsCostComponent(r),
			WMLTextExtractionCatOneResourceUnitsCostComponent(r),
			WMLTextExtractionCatTwoResourceUnitsCostComponent(r, "PAGES_CATEGORY_TWO"),
			WMLIBMModelResourceUnitsCostComponent(r),
			WML3rdPartyModelResourceUnitsCostComponent(r),
			WMLMistralLargeInputResourceUnitsCostComponent(r),
			WMLInstructLabDataResourceUnitsCostComponent(r),
			WMLInstructLabTuningResourceUnitsCostComponent(r),
			// Categorial GPU selection
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_MISTRAL_LARGE"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_MISTRAL_ONE_GPU"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_MISTRAL_TWO_GPU"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_CATEGORY_ONE"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_CATEGORY_TWO"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_CATEGORY_THREE"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_CATEGORY_FOUR"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_CATEGORY_FIVE"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_CATEGORY_SIX"),
			// Specific GPU selection
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_ONE_L_FORTY_S"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_TWO_L_FORTY_S"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_ONE_A_ONE_HUNDRED"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_TWO_A_ONE_HUNDRED"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_FOUR_A_ONE_HUNDRED"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_EIGHT_A_ONE_HUNDRED"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_ONE_H_ONE_HUNDRED"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_TWO_H_ONE_HUNDRED"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_FOUR_H_ONE_HUNDRED"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_EIGHT_H_ONE_HUNDRED"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_ONE_H_TWO_HUNDRED"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_TWO_H_TWO_HUNDRED"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_FOUR_H_TWO_HUNDRED"),
			WMLModelHostingGPUHoursCostComponent(r, "HOURS_EIGHT_H_TWO_HUNDRED"),
		}
	} else if r.Plan == "v2-standard" {
		return []*schema.CostComponent{
			WMLEssentialsCapacityUnitHoursCostComponent(r),
			WMLMistralLargeOutputResourceUnitsCostComponent(r),
			WMLTextExtractionCatOneResourceUnitsCostComponent(r),
			// the unit name is spelled wrong for the standard plan
			WMLTextExtractionCatTwoResourceUnitsCostComponent(r, "PAGES_CATAGORY_TWO"), //nolint:misspell
			WMLIBMModelResourceUnitsCostComponent(r),
			WML3rdPartyModelResourceUnitsCostComponent(r),
			WMLMistralLargeInputResourceUnitsCostComponent(r),
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
	if r.WML_CUHHours != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_CUHHours))
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
	if r.WML_CUHHours != nil {

		// standard plan is billed a fixed amount for each instance, which includes 2500 CUH's per instance.
		// if the used CUH exceeds the included quantity, the overage is charged at a flat rate.
		additional_cuh := *r.WML_CUHHours - (CUH_PER_INSTANCE * instance)
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

func WMLMistralLargeOutputResourceUnitsCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_MistralLargeOutputRU != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_MistralLargeOutputRU))
	}
	return &schema.CostComponent{
		Name:            "Mistral Large Output Resource Unit",
		Unit:            "RU",
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

func WMLTextExtractionCatOneResourceUnitsCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_TextExtractCat1Pages != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_TextExtractCat1Pages))
	}
	return &schema.CostComponent{
		Name:            "Text Extraction Category 1",
		Unit:            "Pages",
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

func WMLTextExtractionCatTwoResourceUnitsCostComponent(r *ResourceInstance, unit string) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_TextExtractCat2Pages != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_TextExtractCat2Pages))
	}
	return &schema.CostComponent{
		Name:            "Text Extraction Category 2",
		Unit:            "Pages",
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
			Unit: strPtr(unit),
		},
	}
}

func WMLIBMModelResourceUnitsCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_IBMModelRU != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_IBMModelRU))
	}
	return &schema.CostComponent{
		Name:            "Resource Units (IBM Models)",
		Unit:            "RU",
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

func WML3rdPartyModelResourceUnitsCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_3rdPartyModelRU != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_3rdPartyModelRU))
	}
	return &schema.CostComponent{
		Name:            "Resource Units (Third Party Models)",
		Unit:            "RU",
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

func WMLMistralLargeInputResourceUnitsCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_MistralLargeInputRU != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_MistralLargeInputRU))
	}
	return &schema.CostComponent{
		Name:            "Mistral Large Input Resource Unit",
		Unit:            "RU",
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

func WMLInstructLabDataResourceUnitsCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_InstructlabDataRU != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_InstructlabDataRU))
	}
	return &schema.CostComponent{
		Name:            "InstructLab Data",
		Unit:            "RU",
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
			Unit: strPtr("INSTRUCTLAB_DATA"),
		},
	}
}

func WMLInstructLabTuningResourceUnitsCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_InstructlabDataRU != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_InstructlabDataRU))
	}
	return &schema.CostComponent{
		Name:            "InstructLab Tuning",
		Unit:            "RU",
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
			Unit: strPtr("INSTRUCTLAB_TUNING"),
		},
	}
}

func WMLModelHostingGPUHoursCostComponent(r *ResourceInstance, unit string) *schema.CostComponent {
	var q *decimal.Decimal
	// Finds the title for the GPU corresponding to the unit
	title := gpuMap[unit]

	if r.WML_ModelHostingHours != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_ModelHostingHours))
	}

	return &schema.CostComponent{
		Name:            title,
		Unit:            "Hours",
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
			Unit: strPtr(unit),
		},
	}
}
