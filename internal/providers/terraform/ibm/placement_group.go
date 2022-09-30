package ibm

import (
	"github.com/infracost/infracost/internal/resources/ibm"
	"github.com/infracost/infracost/internal/schema"
)

func getPlacementGroupRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "ibm_placement_group",
		RFunc: newPlacementGroup,
	}
}

func newPlacementGroup(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	region := d.Get("region").String()
	name := d.Get("pi_placement_group_name").String()

	r := &ibm.PlacementGroup{
		Address: d.Address,
		Region:  region,
		Name:    name,
	}
	r.PopulateUsage(u)

	return r.BuildResource()
}
