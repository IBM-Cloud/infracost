package ibm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// PiInstance struct represents a Virtual Power Systems instance
//
// Resource information: https://www.ibm.com/products/power-virtual-server
// Pricing information: https://cloud.ibm.com/catalog/services/power-systems-virtual-server
// Detailed pricing information: https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-pricing-virtual-server

type PiInstance struct {
	Address                string
	Region                 string
	ProcessorMode          string
	SystemType             string
	StorageType            string
	OperatingSystem        int64
	Memory                 float64
	Cpus                   float64
	LegacyIBMiImageVersion bool

	Storage                   *float64 `infracost_usage:"storage"`
	CloudStorageSolution      *int64   `infracost_usage:"cloud_storage_solution"`
	HighAvailability          *int64   `infracost_usage:"high_availability"`
	DB2WebQuery               *int64   `infracost_usage:"db2_web_query"`
	RationalDevStudioLicences *int64   `infracost_usage:"rational_dev_studio_licenses"`
	Profile                   *string  `infracost_usage:"profile"`
	Epic                      *int64   `infracost_usage:"epic"`
}

// Operating System
const (
	AIX int64 = iota
	IBMI
	RHEL
	SLES
)

const s922 string = "s922"
const e980 string = "e980"
const e1080 string = "e1080"

// PiInstanceUsageSchema defines a list which represents the usage schema of PiInstance.
var PiInstanceUsageSchema = []*schema.UsageItem{
	{Key: "storage", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "cloud_storage_solution", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "high_availability", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "db2_web_query", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "rational_dev_studio_licenses", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "profile", DefaultValue: "", ValueType: schema.String},
	{Key: "epic", DefaultValue: 0, ValueType: schema.Int64},
}

// PopulateUsage parses the u schema.UsageData into the PiInstance.
// It uses the `infracost_usage` struct tags to populate data into the PiInstance.
func (r *PiInstance) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid PiInstance struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *PiInstance) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		r.piInstanceStorageCostComponent(),
	}

	if r.OperatingSystem != SLES && r.Profile == nil {
		costComponents = append(costComponents, r.piInstanceCoresCostComponent(), r.piInstanceMemoryCostComponent())
	}

	if r.OperatingSystem != AIX && r.OperatingSystem != IBMI {
		costComponents = append(costComponents, r.piInstanceLinuxOperatingSystemCostComponent())
	}

	if r.OperatingSystem == AIX {
		costComponents = append(costComponents, r.piInstanceAIXOperatingSystemCostComponent())
	} else if r.OperatingSystem == IBMI {
		costComponents = append(costComponents,
			r.piInstanceIBMiLPPPOperatingSystemCostComponent(),
			r.piInstanceIBMiOSOperatingSystemCostComponent(),
			r.piInstanceCloudStorageSolutionCostComponent(),
			r.piInstanceHighAvailabilityCostComponent(),
			r.piInstanceDB2WebQueryCostComponent(),
			r.piInstanceRationalDevStudioLicensesCostComponent(),
		)
		if r.LegacyIBMiImageVersion {
			costComponents = append(costComponents, r.piInstanceIBMiOperatingSystemServiceExtensionCostComponent())
		}
	} else if r.OperatingSystem == SLES && r.Profile != nil {
		costComponents = append(costComponents, r.piInstanceMemoryHanaProfileCostComponent(), r.piInstanceCoresHanaProfileCostComponent())
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    PiInstanceUsageSchema,
		CostComponents: costComponents,
	}
}

func (r *PiInstance) piInstanceLinuxOperatingSystemCostComponent() *schema.CostComponent {

	c := schema.CostComponent{
		Name:            "Linux Operating System",
		Unit:            "Instance",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			ProductFamily: strPtr("service"),
			Service:       strPtr("power-iaas"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("power-virtual-server-group")},
				{Key: "planType", Value: strPtr("Paid")},
			},
		},
	}

	c.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))

	return &c
}

func (r *PiInstance) piInstanceAIXOperatingSystemCostComponent() *schema.CostComponent {
	unit := ""

	if r.OperatingSystem == AIX {
		if r.SystemType == s922 {
			unit = "AIX_SMALL_APPLICATION_INSTANCE_HOURS"
		} else if r.SystemType == e980 || r.SystemType == e1080 {
			unit = "AIX_MEDIUM_APPLICATION_INSTANCE_HOURS"
		}
	}

	return &schema.CostComponent{
		Name:           "Operating System",
		Unit:           "Cores",
		UnitMultiplier: schema.HourToMonthUnitMultiplier,
		HourlyQuantity: decimalPtr(decimal.NewFromFloat(r.Cpus)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			ProductFamily: strPtr("service"),
			Service:       strPtr("power-iaas"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("power-virtual-server-group")},
				{Key: "planType", Value: strPtr("Paid")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}

func (r *PiInstance) piInstanceIBMiLPPPOperatingSystemCostComponent() *schema.CostComponent {
	unit := ""

	if r.OperatingSystem == IBMI {
		if r.SystemType == s922 {
			unit = "IBMI_LPP_PTEN_APPLICATION_INSTANCE_HOURS"
		} else if r.SystemType == e980 {
			unit = "IBMI_LPP_PTHIRTY_APPLICATION_INSTANCE_HOURS"
		}
	}

	return &schema.CostComponent{
		Name:           "Operating System IBMi LPP",
		Unit:           "Core",
		UnitMultiplier: schema.HourToMonthUnitMultiplier,
		HourlyQuantity: decimalPtr(decimal.NewFromFloat(r.Cpus)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			ProductFamily: strPtr("service"),
			Service:       strPtr("power-iaas"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("power-virtual-server-group")},
				{Key: "planType", Value: strPtr("Paid")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}

func (r *PiInstance) piInstanceIBMiOSOperatingSystemCostComponent() *schema.CostComponent {
	unit := ""

	if r.OperatingSystem == IBMI {
		if r.SystemType == s922 {
			unit = "IBMI_OS_PTEN_APPLICATION_INSTANCE_HOURS"
		} else if r.SystemType == e980 {
			unit = "IBMI_OS_PTHIRTY_APPLICATION_INSTANCE_HOURS"
		}
	}

	return &schema.CostComponent{
		Name:           "Operating System IBMi OS",
		Unit:           "Core",
		UnitMultiplier: schema.HourToMonthUnitMultiplier,
		HourlyQuantity: decimalPtr(decimal.NewFromFloat(r.Cpus)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			ProductFamily: strPtr("service"),
			Service:       strPtr("power-iaas"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("power-virtual-server-group")},
				{Key: "planType", Value: strPtr("Paid")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}

func (r *PiInstance) piInstanceIBMiOperatingSystemServiceExtensionCostComponent() *schema.CostComponent {
	unit := "IBM_I_OS_PTEN_SRVC_EXT_PER_PROC_CORE_HR"

	if r.SystemType == e980 {
		unit = "IBM_I_SERVICE_EXTENSION_PER_CORE_HOUR"
	}

	return &schema.CostComponent{
		Name:           "Operating System IBMi Service Extension",
		Unit:           "Core",
		UnitMultiplier: schema.HourToMonthUnitMultiplier,
		HourlyQuantity: decimalPtr(decimal.NewFromFloat(r.Cpus)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			ProductFamily: strPtr("service"),
			Service:       strPtr("power-iaas"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("power-virtual-server-group")},
				{Key: "planType", Value: strPtr("Paid")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}

func (r *PiInstance) piInstanceMemoryHanaProfileCostComponent() *schema.CostComponent {
	var memoryAmount int64

	if r.Profile != nil {
		coresAndMemory := strings.Split(*r.Profile, "-")[1]
		memoryString := strings.Split(coresAndMemory, "x")[1]
		memory, err := strconv.Atoi(memoryString)
		if err != nil {
			memoryAmount = 0
		} else {
			memoryAmount = int64(memory)
		}
	}

	unit := "MEMHANA_APPLICATION_INSTANCE_HOURS"

	return &schema.CostComponent{
		Name:           "Linux HANA Memory",
		Unit:           "Memory",
		UnitMultiplier: schema.HourToMonthUnitMultiplier,
		HourlyQuantity: decimalPtr(decimal.NewFromInt(memoryAmount)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			ProductFamily: strPtr("service"),
			Service:       strPtr("power-iaas"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("power-virtual-server-group")},
				{Key: "planType", Value: strPtr("Paid")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}

func (r *PiInstance) piInstanceCoresHanaProfileCostComponent() *schema.CostComponent {
	var coresAmount int64

	if r.Profile != nil {
		coresAndMemory := strings.Split(*r.Profile, "-")[1]
		coresString := strings.Split(coresAndMemory, "x")[0]
		cores, err := strconv.Atoi(coresString)
		if err != nil {
			coresAmount = 0
		} else {
			coresAmount = int64(cores)
		}
	}

	unit := "COREHANA_APPLICATION_INSTANCE_HOURS"

	return &schema.CostComponent{
		Name:           "Linux HANA Cores",
		Unit:           "Core",
		UnitMultiplier: schema.HourToMonthUnitMultiplier,
		HourlyQuantity: decimalPtr(decimal.NewFromInt(coresAmount)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			ProductFamily: strPtr("service"),
			Service:       strPtr("power-iaas"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("power-virtual-server-group")},
				{Key: "planType", Value: strPtr("Paid")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}

func (r *PiInstance) piInstanceCloudStorageSolutionCostComponent() *schema.CostComponent {
	var cloudStorageSolutionAmount int64

	if r.CloudStorageSolution != nil {
		cloudStorageSolutionAmount = int64(*r.CloudStorageSolution)
	}

	unit := "IBMI_CSS_APPLICATION_INSTANCE_HOURS"

	return &schema.CostComponent{
		Name:           "Cloud Storage Solution",
		Unit:           "Core",
		UnitMultiplier: schema.HourToMonthUnitMultiplier,
		HourlyQuantity: decimalPtr(decimal.NewFromFloat(r.Cpus * float64(cloudStorageSolutionAmount))),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			ProductFamily: strPtr("service"),
			Service:       strPtr("power-iaas"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("power-virtual-server-group")},
				{Key: "planType", Value: strPtr("Paid")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}

func (r *PiInstance) piInstanceHighAvailabilityCostComponent() *schema.CostComponent {
	var highAvailabilityAmount int64

	if r.HighAvailability != nil {
		highAvailabilityAmount = int64(*r.HighAvailability)
	}

	unit := "IBMIHA_PTHIRTY_APPLICATION_INSTANCES"

	return &schema.CostComponent{
		Name:           "High Availability",
		Unit:           "Core",
		UnitMultiplier: schema.HourToMonthUnitMultiplier,
		HourlyQuantity: decimalPtr(decimal.NewFromFloat(r.Cpus * float64(highAvailabilityAmount))),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			ProductFamily: strPtr("service"),
			Service:       strPtr("power-iaas"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("power-virtual-server-group")},
				{Key: "planType", Value: strPtr("Paid")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}

func (r *PiInstance) piInstanceDB2WebQueryCostComponent() *schema.CostComponent {
	var db2WebQueryAmount int64

	if r.DB2WebQuery != nil {
		db2WebQueryAmount = int64(*r.DB2WebQuery)
	}

	unit := "IBMI_DBIIWQ_APPLICATION_INSTANCE_HOURS"

	return &schema.CostComponent{
		Name:           "IBM DB2 Web Query",
		Unit:           "Core",
		UnitMultiplier: schema.HourToMonthUnitMultiplier,
		HourlyQuantity: decimalPtr(decimal.NewFromFloat(r.Cpus * float64(db2WebQueryAmount))),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			ProductFamily: strPtr("service"),
			Service:       strPtr("power-iaas"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("power-virtual-server-group")},
				{Key: "planType", Value: strPtr("Paid")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}

func (r *PiInstance) piInstanceRationalDevStudioLicensesCostComponent() *schema.CostComponent {
	var RationalDevStudioLicencesAmount int64

	if r.RationalDevStudioLicences != nil {
		RationalDevStudioLicencesAmount = int64(*r.RationalDevStudioLicences)
	}

	unit := "IBMIRDS_APPLICATION_INSTANCES"

	return &schema.CostComponent{
		Name:           "Rational Dev Studio",
		Unit:           "Instance",
		UnitMultiplier: schema.HourToMonthUnitMultiplier,
		HourlyQuantity: decimalPtr(decimal.NewFromInt(RationalDevStudioLicencesAmount)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			ProductFamily: strPtr("service"),
			Service:       strPtr("power-iaas"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("power-virtual-server-group")},
				{Key: "planType", Value: strPtr("Paid")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}

func (r *PiInstance) piInstanceCoresCostComponent() *schema.CostComponent {
	q := decimalPtr(decimal.NewFromFloat(r.Cpus))

	epicEnabled := r.Epic != nil && *r.Epic == 1

	unit := ""

	if r.ProcessorMode == "shared" {
		if r.SystemType == s922 {
			unit = "SOS_VIRTUAL_PROCESSOR_CORE_HOURS"
		} else if r.SystemType == e980 {
			unit = "ESS_VIRTUAL_PROCESSOR_CORE_HOURS"
		} else if r.SystemType == e1080 {
			unit = "PTEN_ESS_VIRTUAL_PROCESSOR_CORE_HRS"
		}
	} else if r.ProcessorMode == "dedicated" {
		if r.SystemType == s922 {
			unit = "SOD_VIRTUAL_PROCESSOR_CORE_HOURS"
		} else if r.SystemType == e980 {
			if epicEnabled {
				unit = "ESS_VIRTUAL_PROCESSOR_CORE_HOURS"
			} else {
				if r.OperatingSystem == SLES && r.Profile == nil {
					unit = "COREHANA_APPLICATION_INSTANCE_HOURS"
				} else {
					unit = "EDD_VIRTUAL_PROCESSOR_CORE_HOURS"
				}
			}
		} else if r.SystemType == e1080 {
			unit = "PTEN_EDD_VIRTUAL_PROCESSOR_CORE_HRS"
		}
	} else if r.ProcessorMode == "capped" {
		if r.SystemType == s922 {
			unit = "SOC_VIRTUAL_PROCESSOR_CORE_HOURS"
		} else if r.SystemType == e980 {
			unit = "ECC_VIRTUAL_PROCESSOR_CORE_HOURS"
		} else if r.SystemType == e1080 {
			unit = "PTEN_ECC_VIRTUAL_PROCESSOR_CORE_HRS"
		}
	}

	name := "Cores"

	if epicEnabled {
		name = "Cores - Epic enabled"
	}

	return &schema.CostComponent{
		Name:           name,
		Unit:           "Core",
		UnitMultiplier: schema.HourToMonthUnitMultiplier,
		HourlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			ProductFamily: strPtr("service"),
			Service:       strPtr("power-iaas"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("power-virtual-server-group")},
				{Key: "planType", Value: strPtr("Paid")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}

func (r *PiInstance) piInstanceMemoryCostComponent() *schema.CostComponent {
	q := decimalPtr(decimal.NewFromFloat(r.Memory))

	unit := "MS_GIGABYTE_HOURS"

	if r.OperatingSystem == SLES {
		unit = "MEMHANA_APPLICATION_INSTANCE_HOURS"
	}

	return &schema.CostComponent{
		Name:           "Memory",
		Unit:           "GB",
		UnitMultiplier: schema.HourToMonthUnitMultiplier,
		HourlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			ProductFamily: strPtr("service"),
			Service:       strPtr("power-iaas"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("power-virtual-server-group")},
				{Key: "planType", Value: strPtr("Paid")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}

func (r *PiInstance) piInstanceStorageCostComponent() *schema.CostComponent {

	var q *decimal.Decimal

	if r.Storage != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.Storage))
	}

	unit := ""

	if r.StorageType == "tier1" {
		unit = "TIER_ONE_STORAGE_GIGABYTE_HOURS"
	} else if r.StorageType == "tier3" {
		unit = "TIER_THREE_STORAGE_GIGABYTE_HOURS"
	}

	return &schema.CostComponent{
		Name:           fmt.Sprintf("Storage - %s", r.StorageType),
		Unit:           "GB",
		UnitMultiplier: schema.HourToMonthUnitMultiplier,
		HourlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			ProductFamily: strPtr("service"),
			Service:       strPtr("power-iaas"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("power-virtual-server-group")},
				{Key: "planType", Value: strPtr("Paid")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}
