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
 * Essentials Plan = essentials
 * Standard Plan = standard
 */
func GetWOCostComponents(r *ResourceInstance) []*schema.CostComponent {
	switch r.Plan {
	case "essentials", "standard":
		return []*schema.CostComponent{
			WOInstanceCostComponent(r),
			WOMonthlyActiveUsersCostComponent(r),
			WOMonthlyVoiceUsersCostComponent(r),
			WOSkillRunsCostComponent(r),
			WOClass1RUCostComponent(r),
			WOClass2RUCostComponent(r),
			WOClass3RUCostComponent(r),
		}
	case "essentials-agentic-mau":
		return []*schema.CostComponent{
			WOInstanceCostComponent(r),
			WOMonthlyActiveUsersCostComponent(r),
			WOMonthlyVoiceUsersCostComponent(r),
			WOOracleHCMAgentCostComponent(r),
			WOWorkdayHCMAgentCostComponent(r),
			WOSAPAgentCostComponent(r),
			WOSourcingContractMgmtAgentCostComponent(r),
			WOLearningDevAgentCostComponent(r),
			WOPurchasingCoupaAgentCostComponent(r),
			WOInvoiceMgmtAgentCostComponent(r),
			WOSupplierMgmtAgentCostComponent(r),
			WOSalesProspectingAgentCostComponent(r),
		}
	case "standard-agentic-mau", "premium-agentic-mau":
		return []*schema.CostComponent{
			WOInstanceCostComponent(r),
			WOMonthlyActiveUsersCostComponent(r),
			WOMonthlyVoiceUsersCostComponent(r),
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
	case "essentials", "essentials-agentic-mau":
		name = "Instance (4000 MAUs included)"
	case "standard", "standard-agentic-mau":
		name = "Instance (40000 MAUs included)"
	default:
		name = "Instance (50000 MAUs included)"
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
	case "essentials", "essentials-agentic-mau":
		included_allotment = ESSENTIAL_MAU_PER_INSTANCE
	case "standard", "standard-agentic-mau":
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

func WOOracleHCMAgentCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WO_AgentOracleHCM != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WO_AgentOracleHCM))
	}

	return &schema.CostComponent{
		Name:            "Oracle HCM Access",
		Unit:            "Domain Agent",
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
			Unit: strPtr("ORACLE_HCM_ACCESS"),
		},
	}
}

func WOWorkdayHCMAgentCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WO_AgentWorkdayHCM != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WO_AgentWorkdayHCM))
	}

	return &schema.CostComponent{
		Name:            "Workday HCM Access",
		Unit:            "Domain Agent",
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
			Unit: strPtr("WORKDAY_HCM_ACCESS"),
		},
	}
}

func WOSAPAgentCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WO_AgentSAPSuccessFactors != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WO_AgentSAPSuccessFactors))
	}

	return &schema.CostComponent{
		Name:            "SAP SuccessFactors Access",
		Unit:            "Domain Agent",
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
			Unit: strPtr("SAP_SUCCESSFACTORS_ACCESS"),
		},
	}
}

func WOSourcingContractMgmtAgentCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WO_AgentSourcingContract != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WO_AgentSourcingContract))
	}

	return &schema.CostComponent{
		Name:            "Sourcing and Contract Management Access",
		Unit:            "Domain Agent",
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
			Unit: strPtr("SOURCING_AND_CONTRACT_MANAGEMENT_ACCESS"),
		},
	}
}

func WOLearningDevAgentCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WO_AgentLearningDev != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WO_AgentLearningDev))
	}

	return &schema.CostComponent{
		Name:            "Learning and Development for Oracle HCM Access",
		Unit:            "Domain Agent",
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
			Unit: strPtr("LEARNING_AND_DEVELOPMENT_ORACLE_HCM_ACCESS"),
		},
	}
}

func WOPurchasingCoupaAgentCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WO_AgentPurchasing != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WO_AgentPurchasing))
	}

	return &schema.CostComponent{
		Name:            "Purchasing with Coupa Access",
		Unit:            "Domain Agent",
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
			Unit: strPtr("PURCHASING_WITH_COUPA_ACCESS"),
		},
	}
}

func WOInvoiceMgmtAgentCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WO_AgentInvoiceMgmt != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WO_AgentInvoiceMgmt))
	}

	return &schema.CostComponent{
		Name:            "Invoice Management with Coupa Access",
		Unit:            "Domain Agent",
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
			Unit: strPtr("INVOICE_MANAGEMENT_COUPA_ACCESS"),
		},
	}
}

func WOSupplierMgmtAgentCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WO_AgentSupplierMgmt != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WO_AgentSupplierMgmt))
	}

	return &schema.CostComponent{
		Name:            "Supplier Management with Coupa Access",
		Unit:            "Domain Agent",
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
			Unit: strPtr("SUPPLIER_MANAGEMENT_COUPA_ACCESS"),
		},
	}
}

func WOSalesProspectingAgentCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.WO_AgentSalesProspecting != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.WO_AgentSalesProspecting))
	}

	return &schema.CostComponent{
		Name:            "Sales Prospecting Access",
		Unit:            "Domain Agent",
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
			Unit: strPtr("SALES_PROSPECTING_ACCESS"),
		},
	}
}
