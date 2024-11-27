package ibm

import (
	"github.com/infracost/infracost/internal/resources/ibm"
	"github.com/infracost/infracost/internal/schema"
)

func getEnDestinationSafariRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:                "ibm_en_destination_safari",
		RFunc:               newEnDestinationSafari,
		ReferenceAttributes: []string{"instance_guid"},
	}
}

func newEnDestinationSafari(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {

	region := d.Get("region").String()
	name := d.Get("name").String()
	is_pre_prod := d.Get("config.0.params.0.pre_prod").Bool()

	var plan string
	enReferenceAttributes := d.References("instance_guid")
	if len(enReferenceAttributes) > 0 {
		plan = enReferenceAttributes[0].Get("plan").String()
	}

	r := &ibm.EnDestinationSafari{
		Address:   d.Address,
		IsPreProd: is_pre_prod,
		Name:      name,
		Plan:      plan,
		Region:    region,
	}
	r.PopulateUsage(u)

	configuration := make(map[string]any)
	configuration["name"] = name
	configuration["plan"] = plan
	configuration["pre-prod"] = is_pre_prod
	configuration["region"] = region

	SetCatalogMetadata(d, d.Type, configuration)

	return r.BuildResource()
}
