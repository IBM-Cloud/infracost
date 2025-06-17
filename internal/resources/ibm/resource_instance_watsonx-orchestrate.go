package ibm

import (
	"fmt"
	"math"

	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

const ESSENTIAL_MAU_PER_INSTANCE float64 = 4000
const STANDARD_MAU_PER_INSTANCE float64 = 40000
const ADDITIONAL_MAU_PER_1000_USERS float64 = 1000

/*
 * https://cloud.ibm.com/catalog/services/watsonx-orchestrate
 * Trial = lite
 * Essentials Plan = essentials
 * Standard Plan = standard
 */
func GetWOCostComponents(r *ResourceInstance) []*schema.CostComponent {
	if (r.Plan == "essentials") || (r.Plan == "standard") {
		return []*schema.CostComponent{
			WOInstanceCostComponent(r),
			WOMonthlyActiveUsersCostComponent(r),
			WOMonthlyVoiceUsersCostComponent(r),
			WOSkillRunsCostComponent(r),
			WOClass1RUCostComponent(r),
			WOClass2RUCostComponent(r),
			WOClass3RUCostComponent(r),
		}
	} else if r.Plan == "lite" {
		costComponent := schema.CostComponent{
			Name:            "Trial plan",
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

func WOInstanceCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	var name string
	if r.WO_Instance != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WO_Instance))
	} else {
		q = decimalPtr(decimal.NewFromInt(1))
	}
	if r.Plan == "essentials" {
		name = "Instance (4000 MAUs included)"
	} else {
		name = "Instance (40000 MAUs included)"
	}
	return &schema.CostComponent{
		Name:            name,
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

func WOMonthlyActiveUsersCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	var users_per_block float64
	var included_allotment float64
	var unit string

	users_per_block = ADDITIONAL_MAU_PER_1000_USERS
	unit = "1K MAU"

	if r.Plan == "essentials" {
		included_allotment = ESSENTIAL_MAU_PER_INSTANCE
	} else {
		included_allotment = STANDARD_MAU_PER_INSTANCE
	}

	// if there are more active users than the monthly allotment of users included in the instance price, then create
	// a cost component for the additional users
	if r.WO_mau != nil {
		additional_users := *r.WO_mau - included_allotment
		if additional_users > 0 {
			// price for additional users charged is per 1k quantity, rounded up,
			// so 1001 additional users will equal 2 blocks of additional users
			q = decimalPtr(decimal.NewFromFloat(math.Ceil(additional_users / users_per_block)))
		}
	}

	return &schema.CostComponent{
		Name:            "Additional Monthly Active Users",
		Unit:            unit,
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
			Unit: strPtr("THOUSAND_MONTHLY_ACTIVE_USERS"),
		},
	}
}

func WOMonthlyVoiceUsersCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	var users_per_block float64
	var unit string

	users_per_block = ADDITIONAL_MAU_PER_1000_USERS
	unit = "1K MAU"

	// price for voice users charged is per 1k quantity, rounded up,
	// so 1001 active users that used voice will equal 2 blocks of voice users
	if r.WO_vu != nil {
		voice_users := math.Ceil(*r.WO_vu / users_per_block)
		q = decimalPtr(decimal.NewFromFloat(voice_users))
	}

	return &schema.CostComponent{
		Name:            "Monthly Active Users using voice",
		Unit:            unit,
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
			Unit: strPtr("THOUSAND_MONTHLY_ACTIVE_VOICE_USERS"),
		},
	}
}

func WOSkillRunsCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WO_skillruns != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WO_skillruns))
	}

	return &schema.CostComponent{
		Name:            "Skill Runs",
		Unit:            "10K Skill Runs",
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
			Unit: strPtr("TEN_THOUSAND_SKILL_RUNS"),
		},
	}
}

func WOClass1RUCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WO_Class1RU != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WO_Class1RU))
	}
	return &schema.CostComponent{
		Name:            "Class 1 Resource Units",
		Unit:            "1K RU",
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
			Unit: strPtr("THOUSAND_CLASS_ONE_RESOURCE_UNITS"),
		},
	}
}

func WOClass2RUCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WO_Class2RU != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WO_Class2RU))
	}
	return &schema.CostComponent{
		Name:            "Class 2 Resource Units",
		Unit:            "1K RU",
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
			Unit: strPtr("THOUSAND_CLASS_TWO_RESOURCE_UNITS"),
		},
	}
}

func WOClass3RUCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WO_Class3RU != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WO_Class3RU))
	}
	return &schema.CostComponent{
		Name:            "Class 3 Resource Units",
		Unit:            "1K RU",
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
			Unit: strPtr("THOUSAND_CLASS_THREE_RESOURCE_UNITS"),
		},
	}
}
