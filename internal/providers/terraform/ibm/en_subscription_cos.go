package ibm

import (
	"github.com/infracost/infracost/internal/resources/ibm"
	"github.com/infracost/infracost/internal/schema"
)

func getEnSubscriptionCosRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:                "ibm_en_subscription_cos",
		RFunc:               newEnSubscriptionCos,
		ReferenceAttributes: []string{"instance_guid"},
	}
}

func newEnSubscriptionCos(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {

	region := d.Get("region").String()
	name := d.Get("name").String()

	var plan string
	enReferenceAttributes := d.References("instance_guid")
	if len(enReferenceAttributes) > 0 {
		plan = enReferenceAttributes[0].Get("plan").String()
	}

	r := &ibm.EnSubscriptionCos{
		Address: d.Address,
		Region:  region,
		Name:    name,
		Plan:    plan,
	}
	r.PopulateUsage(u)

	configuration := make(map[string]any)
	configuration["name"] = name
	configuration["plan"] = plan
	configuration["region"] = region

	SetCatalogMetadata(d, d.Type, configuration)

	return r.BuildResource()
}