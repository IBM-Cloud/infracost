package ibm

import (
	"github.com/infracost/infracost/internal/resources/ibm"
	"github.com/infracost/infracost/internal/schema"
)

func getCodeEngineBuildRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "ibm_code_engine_build",
		RFunc: newCodeEngineBuild,
	}
}

func newCodeEngineBuild(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	region := d.Get("region").String()
	strategysize := d.Get("strategy_size").String()
	r := &ibm.CodeEngineBuild{
		Address:      d.Address,
		Region:       region,
		StrategySize: strategysize,
	}
	r.PopulateUsage(u)

	configuration := make(map[string]any)
	configuration["region"] = region
	configuration["strategysize"] = strategysize

	SetCatalogMetadata(d, d.Type, configuration)

	return r.BuildResource()
}
