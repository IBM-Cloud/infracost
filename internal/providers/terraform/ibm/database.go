package ibm

import (
	"github.com/infracost/infracost/internal/resources/ibm"
	"github.com/infracost/infracost/internal/schema"
)

func getDatabaseRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "ibm_database",
		RFunc: newDatabase,
	}
}

func newDatabase(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	plan := d.Get("plan").String()
	location := d.Get("location").String()
	service := d.Get("service").String()
	name := d.Get("name").String()
	disk := d.Get("disk.allocation_mb").Int()
	var flavor string
	var memory int64
	var cpu int64
	var members int64

	for _, g := range d.Get("group").Array() {
		if g.Get("group_id").String() == "member" {
			flavor = d.Get("host_flavor.id").String()
			memory = d.Get("memory.allocation_mb").Int()
			cpu = d.Get("cpu.allocation_mb").Int()
			members = d.Get("members.allocation_count").Int()
		}
	}

	r := &ibm.Database{
		Name:     name,
		Address:  d.Address,
		Service:  service,
		Plan:     plan,
		Location: location,
		Group:    d.RawValues,
		Flavor:   flavor,
		Disk:     disk,
		Memory:   memory,
		CPU:      cpu,
		Members:  members,
	}
	r.PopulateUsage(u)

	configuration := make(map[string]any)
	configuration["service"] = service
	configuration["plan"] = plan
	configuration["location"] = location
	configuration["disk"] = disk
	configuration["members"] = members

	if flavor != "" {
		configuration["flavor"] = flavor
	} else {
		configuration["memory"] = memory
		configuration["cpu"] = cpu
	}

	SetCatalogMetadata(d, service, configuration)

	return r.BuildResource()
}
