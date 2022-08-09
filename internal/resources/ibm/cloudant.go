package ibm

import (
	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// Cloudant struct represents <TODO: cloud service short description>.
//
// <TODO: Add any important information about the resource and links to the
// pricing pages or documentation that might be useful to developers in the future, e.g:>
//
// Resource information: https://cloud.ibm.com/<PATH/TO/RESOURCE>/
// Pricing information: https://cloud.ibm.com/<PATH/TO/PRICING>/
type Cloudant struct {
	Address  string
	Region   string
	Plan     string
	Capacity string

	Storage *int64 `infracost_usage:"storage"`
}

// CloudantUsageSchema defines a list which represents the usage schema of Cloudant.
var CloudantUsageSchema = []*schema.UsageItem{
	{Key: "storage", ValueType: schema.Int64, DefaultValue: 0},
}

// PopulateUsage parses the u schema.UsageData into the Cloudant.
// It uses the `infracost_usage` struct tags to populate data into the Cloudant.
func (r *Cloudant) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (r *Cloudant) cloudantInstanceCostComponent() *schema.CostComponent {
	paidRegions := []string{"br-sao", "ca-tor", "jp-osa", "in-che"}

	purchaseOption := "100"
	planType := "paygo"
	planName := "standard"

	if len(r.Plan) > 0 {
		planName = r.Plan
	}

	if planName == "dedicated-hardware" {
		purchaseOption = "1"
	}

	if contains(paidRegions, r.Region) {
		planType = "Paid"
	}

	return &schema.CostComponent{
		Name:            "Cloudant",
		Unit:            "Instance",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			Service:       strPtr("cloudantnosqldb"),
			ProductFamily: strPtr("service"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planType", Value: strPtr(planType)},
				{Key: "planName", Value: strPtr(planName)},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr(purchaseOption),
		},
	}
}

func (r *Cloudant) cloudantFreeStorageCostComponent() *schema.CostComponent {
	var q *decimal.Decimal
	if r.Storage != nil {
		q = decimalPtr(decimal.NewFromInt(*r.Storage))
		if q.GreaterThan(decimal.NewFromInt(20)) {
			q = decimalPtr(decimal.NewFromInt(20))
		}
	}

	costComponent := schema.CostComponent{
		Name:            "Free Estimated storage",
		Unit:            "GB",
		MonthlyQuantity: q,
		UnitMultiplier:  decimal.NewFromInt(1),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			Service:       strPtr("cloudantnosqldb"),
			ProductFamily: strPtr("service"),
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("GB_STORAGE_ACCRUED_PER_MONTH"),
		},
	}

	costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))

	return &costComponent
}

func (r *Cloudant) cloudantStorageCostComponent() *schema.CostComponent {

	var q *decimal.Decimal
	if r.Storage != nil {
		q = decimalPtr(decimal.NewFromInt(*r.Storage))
		if q.LessThanOrEqual(decimal.NewFromInt(20)) {
			q = decimalPtr(decimal.NewFromInt(0))
		} else {
			q = decimalPtr(q.Sub(decimal.NewFromInt(20)))
		}
	}

	costComponent := schema.CostComponent{
		Name:            "Estimated storage",
		Unit:            "GB",
		MonthlyQuantity: q,
		UnitMultiplier:  decimal.NewFromInt(1),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			Service:       strPtr("cloudantnosqldb"),
			ProductFamily: strPtr("service"),
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("GB_STORAGE_ACCRUED_PER_MONTH"),
		},
	}

	return &costComponent
}

// BuildResource builds a schema.Resource from a valid Cloudant struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *Cloudant) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		r.cloudantInstanceCostComponent(),
		r.cloudantFreeStorageCostComponent(),
		r.cloudantStorageCostComponent(),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    CloudantUsageSchema,
		CostComponents: costComponents,
	}
}
