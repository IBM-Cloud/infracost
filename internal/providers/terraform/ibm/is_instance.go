package ibm

import (
	"github.com/infracost/infracost/internal/resources/ibm"
	"github.com/infracost/infracost/internal/schema"

	"strings"
)

func getIsInstanceRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "ibm_is_instance",
		RFunc: newIsInstance,
	}
}

// valid profile values https://cloud.ibm.com/docs/vpc?topic=vpc-profiles&interface=ui
// profile names in Global Catalog contain dots instead of dashes
func newIsInstance(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	region := d.Get("region").String()
	profile := d.Get("profile").String()
	zone := d.Get("zone").String()
	dedicatedHost := strings.TrimSpace(d.Get("dedicated_host").String())
	dedicatedHostGroup := strings.TrimSpace(d.Get("dedicated_host_group").String())
	isDedicated := !((dedicatedHost == "") && (dedicatedHostGroup == ""))
	name := d.Get("name").String()

	bootVolume := make([]struct {
		Name string
		Size int64
	}, 0)

	bootVolumeParse := d.Get("boot_volume").Array()
	if len(bootVolumeParse) > 0 {
		for _, volume := range bootVolumeParse {
			name := volume.Get("name").String()
			if name == "" {
				name = "Unnamed boot volume"
			}
			size := volume.Get("size").Int()
			if size == 0 {
				size = 100
			}
			bootVolume = append(bootVolume, struct {
				Name string
				Size int64
			}{Name: name, Size: size})
		}
	}

	r := &ibm.IsInstance{
		Address:     d.Address,
		Region:      region,
		Profile:     profile,
		Zone:        zone,
		IsDedicated: isDedicated,
		BootVolume:  bootVolume,
	}

	r.PopulateUsage(u)

	configuration := make(map[string]any)
	configuration["name"] = name
	configuration["on_dedicated_host"] = isDedicated
	configuration["profile"] = profile
	configuration["region"] = region

	SetCatalogMetadata(d, d.Type, configuration)

	return r.BuildResource()
}
