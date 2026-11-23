package ibm

import (
	"fmt"
	"math"

	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

const ESSENTIAL_MAU_PER_INSTANCE float64 = 4000
const STANDARD_MAU_PER_INSTANCE float64 = 40000
const PREMIUM_MAU_PER_INSTANCE float64 = 50000
const ADDITIONAL_MAU_PER_1000_USERS float64 = 1000

/*
 * https://cloud.ibm.com/catalog/services/watsonx-orchestrate
 * Trial = lite
 * Essentials Agentic with MAU = essentials-agentic-mau
 * Standard Agentic with MAU = standard-agentic-mau
 * Premium Agentic with MAU (data isolation) = premium-agentic-mau
 * Premium Agentic with messages = premium-agentic-messages
 */
func GetWOCostComponents(r *ResourceInstance) []*schema.CostComponent {
	switch r.Plan {
	case "essentials-agentic-mau", "standard-agentic-mau", "premium-agentic-mau":
		return []*schema.CostComponent{
			WOInstanceCostComponent(r),
			WOMonthlyActiveUsersCostComponent(r),
			WOMonthlyVoiceUsersCostComponent(r),
			WOPremierVoiceUsersCostComponent(r),
			WOPremierLLMMessagesCostComponent(r),
			WOThirdPartyMinutesCostComponent(r),
			WOThirdPartyCharactersCostComponent(r),
		}
	case "premium-agentic-messages":
		return []*schema.CostComponent{
			WOInstanceCostComponent(r),
			WOMessagesCostComponent(r),
			WOMinutesCostComponent(r),
		}
	case "lite":
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
	default:
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
	switch r.Plan {
	case "essentials-agentic-mau":
		name = "Instance (4000 MAUs included)"
	case "standard-agentic-mau":
		name = "Instance (40000 MAUs included)"
	case "premium-agentic-mau":
		name = "Instance (50000 MAUs included)"
	default:
		name = "Instance"
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

	switch r.Plan {
	case "essentials-agentic-mau":
		included_allotment = ESSENTIAL_MAU_PER_INSTANCE
	case "standard-agentic-mau":
		included_allotment = STANDARD_MAU_PER_INSTANCE
	default:
		included_allotment = PREMIUM_MAU_PER_INSTANCE
	}

	// if there are more active users than the monthly allotment of users included in the instance price, then create
	// a cost component for the additional users
	if r.WO_mau != nil {
		additional_users := *r.WO_mau - included_allotment
		if additional_users > 0 {
			// price for additional users charged is per 1k quantity, rounded up,
			// so 1001 additional users will equal 2 blocks of additional users
			q = decimalPtr(decimal.NewFromFloat(math.Ceil(additional_users / users_per_block)))
		} else {
			q = decimalPtr(decimal.Zero)
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

// woUsageCostComponent builds a usage-based cost component. The usage value is
// divided by blockSize and rounded up, because IBM charges per block of units.
func woUsageCostComponent(r *ResourceInstance, name, unit, priceUnit string, usage *float64, blockSize float64) *schema.CostComponent {
	var q *decimal.Decimal
	if usage != nil {
		q = decimalPtr(decimal.NewFromFloat(math.Ceil(*usage / blockSize)))
	}

	return &schema.CostComponent{
		Name:            name,
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
			Unit: strPtr(priceUnit),
		},
	}
}

func WOPremierVoiceUsersCostComponent(r *ResourceInstance) *schema.CostComponent {
	return woUsageCostComponent(r, "Monthly Active Users using premier voice", "1K MAU", "PREMIER_VOICE_USERS_PER_MONTH", r.WO_PremierVoiceUsers, 1000)
}

func WOPremierLLMMessagesCostComponent(r *ResourceInstance) *schema.CostComponent {
	return woUsageCostComponent(r, "Premier LLM messages", "10K messages", "PREMIER_LLM_MESSAGES", r.WO_PremierLLMMessages, 10000)
}

func WOThirdPartyMinutesCostComponent(r *ResourceInstance) *schema.CostComponent {
	return woUsageCostComponent(r, "Deepgram STT minutes", "minute", "THIRD_PARTY_MINUTE", r.WO_ThirdPartyMinutes, 1)
}

func WOThirdPartyCharactersCostComponent(r *ResourceInstance) *schema.CostComponent {
	return woUsageCostComponent(r, "Deepgram TTS characters", "1K characters", "THIRD_PARTY_CHARACTER", r.WO_ThirdPartyCharacters, 1000)
}

func WOMessagesCostComponent(r *ResourceInstance) *schema.CostComponent {
	return woUsageCostComponent(r, "Messages", "10K messages", "TEN_THOUSAND_MESSAGES", r.WO_Messages, 10000)
}

func WOMinutesCostComponent(r *ResourceInstance) *schema.CostComponent {
	return woUsageCostComponent(r, "Voice minutes", "minute", "MINUTES", r.WO_Minutes, 1)
}
