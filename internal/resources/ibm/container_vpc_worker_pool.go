package ibm

import (
	"fmt"
	"strings"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// ContainerVpcWorkerPool struct represents an IBM Gen2 VPC worker pool, IKS(Kubernetes) or ROKS(OpenShift), depending upon flavor.
// IKS
// Catalog: https://cloud.ibm.com/kubernetes/catalog/create - at infrastructure select VPC
// Pricing: https://cloud.ibm.com/kubernetes/catalog/about/#pricing
// Docs: https://cloud.ibm.com/docs/containers?topic=containers-getting-started
//
// ROKS
// Catalog: https://cloud.ibm.com/kubernetes/catalog/create?platformType=openshift - at infrastructure select VPC
// Pricing: https://cloud.ibm.com/kubernetes/catalog/about?platformType=openshift#pricing
// Docs: https://cloud.ibm.com/docs/openshift?topic=openshift-getting-started
//
// VPC Gen2 Flavors: https://cloud.ibm.com/docs/containers?topic=containers-vpc-gen2-flavors
type ContainerVpcWorkerPool struct {
	Address              string
	Region               string
	KubeVersion          string
	Flavor               string
	WorkerCount          int64
	Zones                []Zone
	Entitlement          bool
	MonthlyInstanceHours *float64 `infracost_usage:"monthly_instance_hours"`
	Name                 string
}

// ContainerVpcWorkerPoolUsageSchema defines a list which represents the usage schema of ContainerVpcWorkerPool.
var ContainerVpcWorkerPoolUsageSchema = []*schema.UsageItem{
	{Key: "monthly_instance_hours", DefaultValue: 0, ValueType: schema.Float64},
}

// PopulateUsage parses the u schema.UsageData into the ContainerVpcWorkerPool.
// It uses the `infracost_usage` struct tags to populate data into the ContainerVpcWorkerPool.
func (r *ContainerVpcWorkerPool) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid ContainerVpcWorkerPool struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *ContainerVpcWorkerPool) BuildResource() *schema.Resource {
	isOpenshift := strings.HasSuffix(strings.ToLower(r.KubeVersion), "openshift")
	operatingSystem := "UBUNTU"
	useOcpPrices := false
	if isOpenshift {
		operatingSystem = "RHEL"
		// if an entitlement is specified, then ocp licensing is already covered. use pricing that
		// does not include ocp charges.
		if !r.Entitlement {
			useOcpPrices = true
		}
	}
	var attributeFilters = []*schema.AttributeFilter{
		{Key: "provider", Value: strPtr("vpc-gen2")},
		{Key: "flavor", Value: strPtr(r.Flavor)},
		{Key: "serverType", Value: strPtr("virtual")},
		{Key: "isolation", Value: strPtr("public")},
		{Key: "catalogRegion", Value: strPtr(r.Region)},
		{Key: "operatingSystem", ValueRegex: strPtr(fmt.Sprintf("/%s/i", operatingSystem))},
	}
	if useOcpPrices {
		attributeFilters = append(attributeFilters, &schema.AttributeFilter{
			Key: "ocpIncluded", Value: strPtr("true"),
		})
	} else {
		attributeFilters = append(attributeFilters, &schema.AttributeFilter{
			Key: "ocpIncluded", Value: strPtr("false"),
		})
	}
	WorkerCount := decimalPtr(decimal.NewFromInt(1))
	if r.WorkerCount != 0 {
		WorkerCount = decimalPtr(decimal.NewFromInt(r.WorkerCount))
	}

	instanceHours := decimalPtr(decimal.NewFromInt(1))
	if r.MonthlyInstanceHours != nil {
		instanceHours = decimalPtr(decimal.NewFromFloat(*r.MonthlyInstanceHours))
	}

	quantity := WorkerCount.Mul(*instanceHours)

	costComponents := []*schema.CostComponent{}

	// If this is for a "default" workerpool, then the cost will already be covered by the cluster creation
	// so do not include it
	if r.Name != "default" {
		for _, zone := range r.Zones {
			zoneCostComponent := &schema.CostComponent{
				Name:            fmt.Sprintf("VPC Container Work Zone flavor: (%s) region: (%s) name: (%s) x(%d) workers", r.Flavor, r.Region, zone.Name, r.WorkerCount),
				Unit:            "hours",
				UnitMultiplier:  decimal.NewFromInt(1),
				MonthlyQuantity: decimalPtr(quantity),
				ProductFilter: &schema.ProductFilter{
					VendorName:       strPtr("ibm"),
					Service:          strPtr("containers-kubernetes"),
					AttributeFilters: attributeFilters,
				},
			}
			costComponents = append(costComponents, zoneCostComponent)
		}
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    ContainerVpcWorkerPoolUsageSchema,
		CostComponents: costComponents,
	}
}
