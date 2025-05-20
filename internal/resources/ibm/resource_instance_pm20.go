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
			WMLMistralLargeOutputResourceUnitsCostComponent(r),
			WMLTextExtractionCatOneCostComponent(r),
			WMLTextExtractionCatTwoCostComponent(r, "PAGES_CATEGORY_TWO"),
			WMLIBMModelResourceUnitsCostComponent(r),
			WML3rdPartyModelResourceUnitsCostComponent(r),
			WMLMistralLargeInputResourceUnitsCostComponent(r),
			WMLSmallModelHostingCostComponent(r),
			WMLMediumModelHostingCostComponent(r),
			WMLLargeModelHostingCostComponent(r),
			WMLExtraLargeModelHostingCostComponent(r),
			WMLExtraSmallModelHostingCostComponent(r),
			WMLVeryLargeModelHostingCostComponent(r),
			WMLMistralLargeModelHostingAccessCostComponent(r),
			WMLInstructLabDataCostComponent(r),
			WMLInstructLabTuningCostComponent(r),
			WMLMistral1GPUModelHostingAccessCostComponent(r),
			WMLMistral2GPUModelHostingAccessCostComponent(r),
		}
	} else if r.Plan == "v2-standard" {
		return []*schema.CostComponent{
			WMLEssentialsCapacityUnitHoursCostComponent(r),
			WMLMistralLargeOutputResourceUnitsCostComponent(r),
			WMLTextExtractionCatOneCostComponent(r),
			WMLTextExtractionCatTwoCostComponent(r, "PAGES_CATAGORY_TWO"),
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

func WMLMistralLargeOutputResourceUnitsCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_MistralLargeOutputRU != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_MistralLargeOutputRU))
	}
	return &schema.CostComponent{
		Name:            "Mistral Large Output Resource Unit",
		Unit:            "RESOURCE_UNITS",
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

func WMLTextExtractionCatOneCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_TextExtractCat1 != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_TextExtractCat1))
	}
	return &schema.CostComponent{
		Name:            "Text Extraction Category 1",
		Unit:            "Page",
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

func WMLTextExtractionCatTwoCostComponent(r *ResourceInstance, unit string) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_TextExtractCat2 != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_TextExtractCat2))
	}
	return &schema.CostComponent{
		Name:            "Text Extraction Category 2",
		Unit:            "Page",
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

func WMLSmallModelHostingCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_SmallModelHosting != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_SmallModelHosting))
	}
	return &schema.CostComponent{
		Name:            "Small Model Hosting",
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
			Unit: strPtr("HOURS_CATEGORY_ONE"),
		},
	}
}

func WMLMediumModelHostingCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_MediumModelHosting != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_MediumModelHosting))
	}
	return &schema.CostComponent{
		Name:            "Medium Model Hosting",
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
			Unit: strPtr("HOURS_CATEGORY_TWO"),
		},
	}
}

func WMLLargeModelHostingCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_LargeModelHosting != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_LargeModelHosting))
	}
	return &schema.CostComponent{
		Name:            "Large Model Hosting",
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
			Unit: strPtr("HOURS_CATEGORY_THREE"),
		},
	}
}

func WMLExtraLargeModelHostingCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_ExtraLargeModelHosting != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_ExtraLargeModelHosting))
	}
	return &schema.CostComponent{
		Name:            "Extra Large Model Hosting",
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
			Unit: strPtr("HOURS_CATEGORY_FOUR"),
		},
	}
}

func WMLExtraSmallModelHostingCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_ExtraSmallModelHosting != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_ExtraSmallModelHosting))
	}
	return &schema.CostComponent{
		Name:            "Extra Small Model Hosting",
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
			Unit: strPtr("HOURS_CATEGORY_FIVE"),
		},
	}
}

func WMLVeryLargeModelHostingCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_VeryLargeModelHosting != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_VeryLargeModelHosting))
	}
	return &schema.CostComponent{
		Name:            "Very Large Model Hosting",
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
			Unit: strPtr("HOURS_CATEGORY_SIX"),
		},
	}
}

func WMLMistralLargeModelHostingAccessCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_MistralLargeModelHostingAccess != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_MistralLargeModelHostingAccess))
	}
	return &schema.CostComponent{
		Name:            "Mistral Large Model Hosting Access",
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
			Unit: strPtr("HOURS_MISTRAL_LARGE"),
		},
	}
}

func WMLInstructLabDataCostComponent(r *ResourceInstance) *schema.CostComponent {
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

func WMLInstructLabTuningCostComponent(r *ResourceInstance) *schema.CostComponent {
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

func WMLMistral2GPUModelHostingAccessCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_Mistral2GPUModelHostingAccess != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_Mistral2GPUModelHostingAccess))
	}
	return &schema.CostComponent{
		Name:            "Mistral 2 GPU Model Hosting Access",
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
			Unit: strPtr("HOURS_MISTRAL_TWO_GPU"),
		},
	}
}

func WMLMistral1GPUModelHostingAccessCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WML_Mistral1GPUModelHostingAccess != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WML_Mistral1GPUModelHostingAccess))
	}
	return &schema.CostComponent{
		Name:            "Mistral 1 GPU Model Hosting Access",
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
			Unit: strPtr("HOURS_MISTRAL_ONE_GPU"),
		},
	}
}
