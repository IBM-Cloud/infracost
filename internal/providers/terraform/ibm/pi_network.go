package ibm

import (
	"github.com/infracost/infracost/internal/resources/ibm"
	"github.com/infracost/infracost/internal/schema"
)

func getPiNetworkRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "ibm_pi_network",
		RFunc: newPiNetwork,
	}
}

func newPiNetwork(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	region := d.Get("region").String()
	name := d.Get("pi_network_name").String()

	r := &ibm.PiNetwork{
		Address: d.Address,
		Region:  region,
		Name:    name,
	}
	r.PopulateUsage(u)

	return r.BuildResource()
}
