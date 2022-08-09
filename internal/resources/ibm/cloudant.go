package ibm

import (
	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
)

// Cloudant struct represents <TODO: cloud service short description>.
//
// <TODO: Add any important information about the resource and links to the
// pricing pages or documentation that might be useful to developers in the future, e.g:>
//
// Resource information: https://cloud.ibm.com/<PATH/TO/RESOURCE>/
// Pricing information: https://cloud.ibm.com/<PATH/TO/PRICING>/
type Cloudant struct {
	Address string
	Region  string
}

// CloudantUsageSchema defines a list which represents the usage schema of Cloudant.
var CloudantUsageSchema = []*schema.UsageItem{
}

// PopulateUsage parses the u schema.UsageData into the Cloudant.
// It uses the `infracost_usage` struct tags to populate data into the Cloudant.
func (r *Cloudant) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid Cloudant struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *Cloudant) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		// TODO: add cost components
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    CloudantUsageSchema,
		CostComponents: costComponents,
	}
}

