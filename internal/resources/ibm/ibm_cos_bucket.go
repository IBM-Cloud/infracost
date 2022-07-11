package ibm

import (
	"fmt"
	"strings"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/infracost/infracost/internal/usage"
	"github.com/shopspring/decimal"
)

// IbmCosBucket struct represents IBM Cloud Object Storage instance
//
// Resource information: https://cloud.ibm.com/objectstorage
// Pricing information: https://cloud.ibm.com/objectstorage/create#pricing

type IbmCosBucket struct {
	Address string
	Region  string
}

// IbmCosBucketUsageSchema defines a list which represents the usage schema of IbmCosBucket.
var IbmCosBucketUsageSchema = []*schema.UsageItem{}

// PopulateUsage parses the u schema.UsageData into the IbmCosBucket.
// It uses the `infracost_usage` struct tags to populate data into the IbmCosBucket.
func (r *IbmCosBucket) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid IbmCosBucket struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *IbmCosBucket) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		// TODO: add cost components
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    IbmCosBucketUsageSchema,
		CostComponents: costComponents,
	}
}
