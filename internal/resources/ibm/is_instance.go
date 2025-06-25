package ibm

import (
	"fmt"
	"regexp"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// IsInstance struct represents an IBM virtual server instance.
//
// Pricing information: https://cloud.ibm.com/kubernetes/catalog/about

type IsInstance struct {
	Address         string
	Region          string
	OperatingSystem int64
	Image           string
	Profile         string // should be values from CLI 'ibmcloud is instance-profiles'
	Zone            string
	IsDedicated     bool // will be true if a dedicated_host or dedicated_host_group is specified
	BootVolume      struct {
		Name string
		Size int64
	}
	MonthlyInstanceHours *float64 `infracost_usage:"monthly_instance_hours"`
}

var IsInstanceUsageSchema = []*schema.UsageItem{
	{Key: "monthly_instance_hours", DefaultValue: 0, ValueType: schema.Float64},
}

// PopulateUsage parses the u schema.UsageData into the IsInstance.
// It uses the `infracost_usage` struct tags to populate data into the IsInstance.
func (r *IsInstance) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *IsInstance) instanceHoursCostComponent() *schema.CostComponent {

	service := "is.reservation"
	planNamePrefix := "instance-"
	unit := "RESERVATION_HOURS_HOURLY"

	isConfidentialProfile, _ := regexp.MatchString("^.*c-.*$", r.Profile)
	if isConfidentialProfile {
		service = "is.instance"
		planNamePrefix = ""
		unit = "INSTANCE_HOURS_MULTI_TENANT"
	}

	planName := fmt.Sprintf("%s%s", planNamePrefix, r.Profile)
	unitMultiplier := int64(1)
	var q *decimal.Decimal

	if r.MonthlyInstanceHours != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.MonthlyInstanceHours))
	}
	if r.IsDedicated {
		q = decimalPtr(decimal.NewFromFloat(0))
		unitMultiplier = 0
	}

	return &schema.CostComponent{
		Name:            fmt.Sprintf("Instance Hours (%s)", r.Profile),
		Unit:            "Hours",
		UnitMultiplier:  decimal.NewFromInt(unitMultiplier),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			Service:       strPtr(service),
			ProductFamily: strPtr("service"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &planName},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}

func (r *IsInstance) bootVolumeCostComponent() *schema.CostComponent {

	var q *decimal.Decimal
	if r.MonthlyInstanceHours != nil {
		q = decimalPtr(decimal.NewFromFloat(float64(r.BootVolume.Size) * (*r.MonthlyInstanceHours)))
	}

	return &schema.CostComponent{
		Name:            fmt.Sprintf("Boot volume (%s, %d GB)", r.BootVolume.Name, r.BootVolume.Size),
		Unit:            "Hours",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			ProductFamily: strPtr("service"),
			Service:       strPtr("is.volume"),
			Region:        strPtr(r.Region),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", ValueRegex: regexPtr(("gen2-volume-general-purpose"))},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("GIGABYTE_HOURS"),
		},
	}
}

func (r *IsInstance) imageHoursCostComponent() *schema.CostComponent {

	imageMap := map[string]string{
		"r006-88da7a09-2f59-4324-ac85-e3165f6323f5": "REDHAT_VCPU_HOURS", //redhat
		"r006-20ba59a0-19f6-457b-ba6f-80ddc988ef10": "REDHAT_VCPU_HOURS",
		"r006-3bcf9a79-b5b2-45a8-ae8d-1bbb338836ef": "REDHAT_VCPU_HOURS",
		"r006-4b854390-a503-4393-8e5d-59e21816e727": "REDHAT_VCPU_HOURS",
		"r006-3133dc35-e3db-498d-b415-c8145cea2e44": "REDHAT_VCPU_HOURS",
		"r006-987002d0-978e-46f3-bc9d-5c44dac7e4ab": "REDHAT_VCPU_HOURS",
		"r006-33fc9c08-2838-4220-ac26-6591e6e1c73b": "REDHAT_VCPU_HOURS",
		"r006-d0af4fc6-23d0-42c8-b5ea-67c61506bf86": "REDHAT_VCPU_HOURS",
		"r006-14d5062e-d4dc-4af7-9b68-248afefbacd4": "REDHAT_VCPU_HOURS",
		"r006-fc0dbd40-b40a-4b54-bc1a-5251f4b5990d": "REDHAT_VCPU_HOURS",
		"r006-04fe51ed-bcb8-4cd3-8981-637ba174c55a": "REDHAT_VCPU_HOURS",
		"r006-95301a19-6b2b-4870-a64b-22f374ab78d6": "REDHAT_VCPU_HOURS",
		"r006-491e1f7f-3193-477b-9e69-2a1839700bd4": "REDHAT_VCPU_HOURS",
		"r006-f62eaf9e-0a50-47f2-9b0e-b9189c9ba81a": "REDHAT_VCPU_HOURS",
		"r006-868b9c19-1cc4-4456-a1ef-9be5b801d4af": "REDHAT_VCPU_HOURS",
		"r006-5f1399ec-0974-46f3-b466-bed6f843593e": "REDHAT_VCPU_HOURS",
		"r006-ef09be7b-01af-44cc-9503-8a07a5183964": "REDHAT_VCPU_HOURS",
		"r006-c717e22f-3c6c-4bd3-b8a2-3342aefe6611": "SUSE_INSTANCE_HOURS", //sles
		"r010-bbe1011a-ab30-460f-863c-75c5c005d659": "SUSE_INSTANCE_HOURS",
		"r006-e19bf9bf-a454-4ba8-a2b5-86eb30a33b96": "SUSE_INSTANCE_HOURS",
		"r006-7154d75e-adb7-4c7b-99fc-9af25e0b8485": "SUSE_INSTANCE_HOURS",
		"r006-91e4ffcd-292c-4283-a66e-776f3b33a04d": "SUSE_INSTANCE_HOURS",
		"r006-e2f7bfd9-7768-44bf-9f4b-95cdf25b214a": "SUSE_INSTANCE_HOURS",
		"r006-d9a00e92-c04b-4a26-b244-e22cd5fd5b1b": "SUSE_INSTANCE_HOURS",
		"r006-8fe7a011-c636-437c-b5a7-b892581d4788": "SUSE_INSTANCE_HOURS",
		"r006-6f37720a-5a4b-437f-8e5f-6002e95ccbc1": "SUSE_INSTANCE_HOURS",
		"r006-5796a93a-e45d-46e3-ab7c-ac89f3b928f0": "WINDOWS_VCPU_HOURS", //windows
		"r006-a7e03025-b620-45f1-9094-1be3d0020de3": "WINDOWS_VCPU_HOURS",
		"r006-c09590a5-8a2c-4699-b682-aebd1e787eb7": "WINDOWS_VCPU_HOURS",
		"r006-4b0761e6-2510-4ebe-a359-761ca67c07ab": "WINDOWS_VCPU_HOURS",
		"r006-be892432-a287-4b5b-a249-af30f58d108d": "WINDOWS_VCPU_HOURS", //windows w sql
	}

	var q *decimal.Decimal

	if r.MonthlyInstanceHours != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.MonthlyInstanceHours))
	}
	return &schema.CostComponent{
		Name:            fmt.Sprintf("Image (%s)", r.Image),
		Unit:            "Hours",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			Service:       strPtr("is.instance"),
			ProductFamily: strPtr("service"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("gen2-instance")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(imageMap[r.Image]),
		},
	}
}

func (r *IsInstance) sqlLicenceCostComponent() *schema.CostComponent {
	var q *decimal.Decimal

	if r.MonthlyInstanceHours != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.MonthlyInstanceHours))
	}

	return &schema.CostComponent{
		Name:            fmt.Sprintf("Image (%s)", r.Image),
		Unit:            "Hours",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			Service:       strPtr("is.instance"),
			ProductFamily: strPtr("service"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("gen2-instance")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("MSSQL_LICENSE_HOURS"),
		},
	}
}

// BuildResource builds a schema.Resource from a valid IsShare struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *IsInstance) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		r.instanceHoursCostComponent(),
		r.bootVolumeCostComponent(),
		r.imageHoursCostComponent(),
	}
	if r.Profile == "r006-be892432-a287-4b5b-a249-af30f58d108d" {
		costComponents = append(costComponents, r.sqlLicenceCostComponent())
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    IsInstanceUsageSchema,
		CostComponents: costComponents,
	}
}
