package ibm

import (
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

func GetLogsCostComponents(r *ResourceInstance) []*schema.CostComponent {
	return []*schema.CostComponent{
		//cost component functions go here
		LogsStoreNSearchCostComponent(r),
		LogsAnalyzeNAlertCostComponent(r),
		LogsPriorityInsightsCostComponent(r, "PRIORITY_INSIGHTS_RETENTION_SEVEN_DAYS", "7"),
		LogsPriorityInsightsCostComponent(r, "PRIORITY_INSIGHTS_RETENTION_FOURTEEN_DAYS", "14"),
		LogsPriorityInsightsCostComponent(r, "PRIORITY_INSIGHTS_RETENTION_THIRTY_DAYS", "30"),
		LogsPriorityInsightsCostComponent(r, "PRIORITY_INSIGHTS_RETENTION_SIXTY_DAYS", "60"),
		LogsPriorityInsightsCostComponent(r, "PRIORITY_INSIGHTS_RETENTION_NINETY_DAYS", "90"),
	}

}

func LogsStoreNSearchCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.Logs_Hours != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.Logs_Hours))
	}

	return &schema.CostComponent{
		Name:            "Store and Search",
		Unit:            "GB/Hour",
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
			Unit: strPtr("STORE_AND_SEARCH"),
		},
	}

}

func LogsAnalyzeNAlertCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.Logs_Hours != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.Logs_Hours))
	}

	return &schema.CostComponent{
		Name:            "Analyze and Alert",
		Unit:            "GB/Hour",
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
			Unit: strPtr("ANALYZE_AND_ALERT"),
		},
	}

}

func LogsPriorityInsightsCostComponent(r *ResourceInstance, unit string, days string) *schema.CostComponent {
	var q *decimal.Decimal
	if r.Logs_Hours != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.Logs_Hours))
	}

	return &schema.CostComponent{
		Name:            "Priority Insights (Retention " + days + " days)",
		Unit:            "GB/Hour",
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
