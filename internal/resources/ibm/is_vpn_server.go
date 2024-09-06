package ibm

import (
	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
)

// IsVpnServer struct represents <TODO: cloud service short description>.
//
// <TODO: Add any important information about the resource and links to the
// pricing pages or documentation that might be useful to developers in the future, e.g:>
//
// Resource information: https://cloud.ibm.com/<PATH/TO/RESOURCE>/
// Pricing information: https://cloud.ibm.com/<PATH/TO/PRICING>/
type IsVpnServer struct {
	Address string
	Region  string
}

// IsVpnServerUsageSchema defines a list which represents the usage schema of IsVpnServer.
var IsVpnServerUsageSchema = []*schema.UsageItem{
}

// PopulateUsage parses the u schema.UsageData into the IsVpnServer.
// It uses the `infracost_usage` struct tags to populate data into the IsVpnServer.
func (r *IsVpnServer) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid IsVpnServer struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *IsVpnServer) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		// TODO: add cost components
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    IsVpnServerUsageSchema,
		CostComponents: costComponents,
	}
}

