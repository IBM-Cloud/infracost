package ibm

import (
	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// PiNetwork struct represents <TODO: cloud service short description>.
//
// <TODO: Add any important information about the resource and links to the
// pricing pages or documentation that might be useful to developers in the future, e.g:>
//
// Resource information: https://cloud.ibm.com/<PATH/TO/RESOURCE>/
// Pricing information: https://cloud.ibm.com/<PATH/TO/PRICING>/
type PiNetwork struct {
	Address string
	Region  string
	Name    string
}

// PiNetworkUsageSchema defines a list which represents the usage schema of PiNetwork.
var PiNetworkUsageSchema = []*schema.UsageItem{}

// PopulateUsage parses the u schema.UsageData into the PiNetwork.
// It uses the `infracost_usage` struct tags to populate data into the PiNetwork.
func (r *PiNetwork) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *PiNetwork) NetworkCostComponent() *schema.CostComponent {
	q := decimalPtr(decimal.NewFromInt(1))

	costComponent := schema.CostComponent{
		Name:            r.Name,
		Unit:            "Instance",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
	}

	costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
	return &costComponent
}

// BuildResource builds a schema.Resource from a valid PiNetwork struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *PiNetwork) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		r.NetworkCostComponent(),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    PiNetworkUsageSchema,
		CostComponents: costComponents,
	}
}
