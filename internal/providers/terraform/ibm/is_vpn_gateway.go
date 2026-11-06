package ibm

import (
	"github.com/infracost/infracost/internal/resources/ibm"
	"github.com/infracost/infracost/internal/schema"
)

func getIsVpnGatewayRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "ibm_is_vpn_gateway",
		RFunc: newIsVpnGateway,
	}
}

func newIsVpnGateway(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	region := d.Get("region").String()
	availabilityMode := d.Get("availability_mode").String()
	r := &ibm.IsVpnGateway{
		Address:          d.Address,
		Region:           region,
		AvailabilityMode: availabilityMode,
	}
	r.PopulateUsage(u)

	configuration := make(map[string]any)
	configuration["region"] = region
	if availabilityMode != "" {
		configuration["availability_mode"] = availabilityMode
	}

	SetCatalogMetadata(d, d.Type, configuration)

	return r.BuildResource()
}
