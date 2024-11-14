package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

var CODE_ENGINE_BUILD_BUILDS_SIZE = map[string]map[string]float64{
	"small": {
		"CPU": 0.5,
		"Memory": 2,
	},
	"medium": {
		"CPU": 1,
		"Memory": 4,
	},
	"large": {
		"CPU": 2,
		"Memory": 8,
	},
	"xlarge": {
		"CPU": 4,
		"Memory": 16,
	},
	"xxlarge": {
		"CPU": 12,
		"Memory": 48,
	},
}

// CodeEngineBuild struct represents <TODO: cloud service short description>.
//
// <TODO: Add any important information about the resource and links to the
// pricing pages or documentation that might be useful to developers in the future, e.g:>
//
// Resource information: https://cloud.ibm.com/docs/codeengine?topic=codeengine-getting-started
// Pricing information: https://cloud.ibm.com/docs/codeengine?topic=codeengine-pricing
type CodeEngineBuild struct {
	Address string
	Region  string
	StrategySize string

	InstanceHours *float64 `infracost_usage:"instance_hours"`
}

// CodeEngineBuildUsageSchema defines a list which represents the usage schema of CodeEngineBuild.
var CodeEngineBuildUsageSchema = []*schema.UsageItem{
	{Key: "instance_hours", DefaultValue: 1, ValueType: schema.Float64},
}

// PopulateUsage parses the u schema.UsageData into the CodeEngineBuild.
// It uses the `infracost_usage` struct tags to populate data into the CodeEngineBuild.
func (r *CodeEngineBuild) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *CodeEngineBuild) CodeEngineBuildVirtualProcessorCoreCostComponent() *schema.CostComponent {
	var ss string = r.StrategySize
	if r.StrategySize == "" {
		ss = "medium"
	}

	var sscpu float64 = CODE_ENGINE_BUILD_BUILDS_SIZE[ss]["CPU"]

	var hours *decimal.Decimal
	if (r.InstanceHours != nil) {
		hours = decimalPtr(decimal.NewFromFloat(*r.InstanceHours * sscpu))
	}
	
	return &schema.CostComponent{
		Name:			fmt.Sprintf("Virtual Processor Cores (%s build)", ss),
		Unit:			"vCPU Hours",
		UnitMultiplier:	decimal.NewFromInt(1),
		MonthlyQuantity: hours,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region: 	strPtr(r.Region),
			Service: 	strPtr("codeengine"),
			AttributeFilters: []*schema.AttributeFilter{
				{
					Key: "planName", Value: strPtr("standard"),
				},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("VIRTUAL_PROCESSOR_CORE_HOURS"),
		},
	}
}

func (r *CodeEngineBuild) CodeEngineBuildRAMCostComponent() *schema.CostComponent {
	var ss string = r.StrategySize
	if r.StrategySize == "" {
		ss = "medium"
	}

	var ssmem float64 = CODE_ENGINE_BUILD_BUILDS_SIZE[ss]["Memory"]

	var hours *decimal.Decimal
	if (r.InstanceHours != nil) {
		hours = decimalPtr(decimal.NewFromFloat(*r.InstanceHours * ssmem))
	}
	
	return &schema.CostComponent{
		Name:			fmt.Sprintf("RAM (%s build)", ss),
		Unit:			"GB Hours",
		UnitMultiplier:	decimal.NewFromInt(1),
		MonthlyQuantity: hours,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region: 	strPtr(r.Region),
			Service: 	strPtr("codeengine"),
			AttributeFilters: []*schema.AttributeFilter{
				{
					Key: "planName", Value: strPtr("standard"),
				},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("GIGABYTE_HOURS"),
		},
	}
}


// BuildResource builds a schema.Resource from a valid CodeEngineBuild struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *CodeEngineBuild) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		r.CodeEngineBuildVirtualProcessorCoreCostComponent(),
		r.CodeEngineBuildRAMCostComponent(),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    CodeEngineBuildUsageSchema,
		CostComponents: costComponents,
	}
}
