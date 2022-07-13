package ibm

import (
	"fmt"
	"strings"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// IbmCosBucket struct represents IBM Cloud Object Storage instance
//
// Resource information: https://cloud.ibm.com/objectstorage
// Pricing information: https://cloud.ibm.com/objectstorage/create#pricing

type IbmCosBucket struct {
	Address      string
	Region       string
	Location     string
	StorageClass string

	MonthlyAverageCapacity *float64 `infracost_usage:"monthly_average_capacity"`
	PublicStandardEgress   *float64 `infracost_usage:"public_standard_egress"`
	ClassARequestCount     *int64   `infracost_usage:"class_a_request_count"`
	ClassBRequestCount     *int64   `infracost_usage:"class_b_request_count"`
	MonthlyDataRetrieval   *float64 `infracost_usage:"monthly_data_retrieval"`
}

// IbmCosBucketUsageSchema defines a list which represents the usage schema of IbmCosBucket.
var IbmCosBucketUsageSchema = []*schema.UsageItem{
	{Key: "monthly_average_capacity", ValueType: schema.Float64, DefaultValue: 0},
	{Key: "public_standard_egress", ValueType: schema.Float64, DefaultValue: 0},
	{Key: "class_a_request_count", ValueType: schema.Int64, DefaultValue: 0},
	{Key: "class_b_request_count", ValueType: schema.Int64, DefaultValue: 0},
	{Key: "monthly_data_retrieval", ValueType: schema.Int64, DefaultValue: 0},
	{
		Key:          "monthly_egress_data_transfer_gb",
		ValueType:    schema.SubResourceUsage,
		DefaultValue: 0,
	},
}

// PopulateUsage parses the u schema.UsageData into the IbmCosBucket.
// It uses the `infracost_usage` struct tags to populate data into the IbmCosBucket.
func (r *IbmCosBucket) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *IbmCosBucket) MonthlyAverageCapacityCostComponent() *schema.CostComponent {
	return &schema.CostComponent{
		Name:           fmt.Sprintf("Storage-%s-%s", strings.ToLower(r.StorageClass), strings.ToLower(r.Region)),
		Unit:           "hours",
		UnitMultiplier: decimal.NewFromInt(1),
		ProductFilter: &schema.ProductFilter{
			VendorName:       strPtr("ibm"),
			Region:           strPtr(r.Region),
			Service:          strPtr(("cloud-object-storage")),
			ProductFamily:    strPtr("iaas"),
			AttributeFilters: []*schema.AttributeFilter{},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("FLEX_MAX_CAP"),
		},
	}
}

func (r *IbmCosBucket) ClassARequestCountCostComponent() *schema.CostComponent {
	return &schema.CostComponent{
		Name:           fmt.Sprintf("Storage-%s-%s", strings.ToLower(r.StorageClass), strings.ToLower(r.Region)),
		Unit:           "hours",
		UnitMultiplier: decimal.NewFromInt(1),
		ProductFilter: &schema.ProductFilter{
			VendorName:       strPtr("ibm"),
			Region:           strPtr(r.Region),
			Service:          strPtr(("cloud-object-storage")),
			ProductFamily:    strPtr("iaas"),
			AttributeFilters: []*schema.AttributeFilter{},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("VAULT_CLASS_A_CALLS"),
		},
	}
}

func (r *IbmCosBucket) ClassBRequestCountCostComponent() *schema.CostComponent {
	return &schema.CostComponent{
		Name:           fmt.Sprintf("Storage-%s-%s", strings.ToLower(r.StorageClass), strings.ToLower(r.Region)),
		Unit:           "hours",
		UnitMultiplier: decimal.NewFromInt(1),
		ProductFilter: &schema.ProductFilter{
			VendorName:       strPtr("ibm"),
			Region:           strPtr(r.Region),
			Service:          strPtr(("cloud-object-storage")),
			ProductFamily:    strPtr("iaas"),
			AttributeFilters: []*schema.AttributeFilter{},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("VAULT_CLASS_B_CALLS"),
		},
	}
}

func (r *IbmCosBucket) PublicStandardEgressCostComponent() *schema.CostComponent {
	return &schema.CostComponent{
		Name:           fmt.Sprintf("Storage-%s-%s", strings.ToLower(r.StorageClass), strings.ToLower(r.Region)),
		Unit:           "hours",
		UnitMultiplier: decimal.NewFromInt(1),
		ProductFilter: &schema.ProductFilter{
			VendorName:       strPtr("ibm"),
			Region:           strPtr(r.Region),
			Service:          strPtr(("cloud-object-storage")),
			ProductFamily:    strPtr("iaas"),
			AttributeFilters: []*schema.AttributeFilter{},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(""),
		},
	}
}

func (r *IbmCosBucket) MonthlyDataRetrievalCostComponent() *schema.CostComponent {

	retrieval := "FLEX_RETRIEVAL"

	if r.StorageClass == "cold" {
		retrieval = "COLD_VAULT_RETRIEVAL"
	}

	if r.StorageClass == "vault" {
		retrieval = "VAULT_RETRIEVAL"
	}

	return &schema.CostComponent{
		Name:           fmt.Sprintf("Storage-%s-%s", strings.ToLower(r.StorageClass), strings.ToLower(r.Region)),
		Unit:           "hours",
		UnitMultiplier: decimal.NewFromInt(1),
		ProductFilter: &schema.ProductFilter{
			VendorName:       strPtr("ibm"),
			Region:           strPtr(r.Region),
			Service:          strPtr(("cloud-object-storage")),
			ProductFamily:    strPtr("iaas"),
			AttributeFilters: []*schema.AttributeFilter{},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(retrieval),
		},
	}
}

// BuildResource builds a schema.Resource from a valid IbmCosBucket struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *IbmCosBucket) BuildResource() *schema.Resource {

	costComponents := []*schema.CostComponent{
		r.MonthlyAverageCapacityCostComponent(),
		r.ClassARequestCountCostComponent(),
		r.ClassBRequestCountCostComponent(),
	}

	if r.StorageClass == "vault" || r.StorageClass == "cold" || r.StorageClass == "smart" {
		costComponents = append(costComponents, r.MonthlyDataRetrievalCostComponent())
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    IbmCosBucketUsageSchema,
		CostComponents: costComponents,
	}
}
