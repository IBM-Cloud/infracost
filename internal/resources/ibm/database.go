package ibm

import (
	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

// Database struct represents a database instance
//
// This terraform resource is opaque and can handle multiple databases, provided with the right parameters
type Database struct {
	Name     string
	Address  string
	Service  string
	Plan     string
	Location string
	Group    gjson.Result

	// Databases For PostgreSQL
	// Catalog Link: https://cloud.ibm.com/catalog/services/databases-for-postgresql
	// Pricing Link: https://cloud.ibm.com/docs/databases-for-postgresql?topic=databases-for-postgresql-pricing
	PostgreSQL_Ram     *int64 `infracost_usage:"postgresql_database_ram_mb"`
	PostgreSQL_Disk    *int64 `infracost_usage:"postgresql_database_disk_mb"`
	PostgreSQL_Core    *int64 `infracost_usage:"postgresql_database_core"`
	PostgreSQL_Members *int64 `infracost_usage:"postgresql_database_members"`

	// Databases For ElasticSearch
	// Catalog Link: https://cloud.ibm.com/catalog/services/databases-for-elasticsearch
	// Pricing Link: https://cloud.ibm.com/docs/databases-for-elasticsearch?topic=databases-for-elasticsearch-pricing
	ElasticSearch_Ram     *int64 `infracost_usage:"elasticsearch_database_ram_mb"`
	ElasticSearch_Disk    *int64 `infracost_usage:"elasticsearch_database_disk_mb"`
	ElasticSearch_Core    *int64 `infracost_usage:"elasticsearch_database_core"`
	ElasticSearch_Members *int64 `infracost_usage:"elasticsearch_database_members"`
}

type DatabaseCostComponentsFunc func(*Database) []*schema.CostComponent

// PopulateUsage parses the u schema.UsageData into the Database.
// It uses the `infracost_usage` struct tags to populate data into the Database.
func (r *Database) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// DatabaseUsageSchema defines a list which represents the usage schema of Database.
var DatabaseUsageSchema = []*schema.UsageItem{
	{Key: "postgresql_database_ram_mb", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "postgresql_database_disk_mb", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "postgresql_database_core", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "postgresql_database_members", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "elasticsearch_database_ram_mb", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "elasticsearch_database_disk_mb", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "elasticsearch_database_core", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "elasticsearch_database_members", DefaultValue: 0, ValueType: schema.Int64},
}

var DatabaseCostMap map[string]DatabaseCostComponentsFunc = map[string]DatabaseCostComponentsFunc{
	"databases-for-postgresql": GetPostgresCostComponents,
	// "databases-for-etcd":
	// "databases-for-redis":
	"databases-for-elasticsearch": GetElasticSearchCostComponents,
	// "messages-for-rabbitmq":
	// "databases-for-mongodb":
	// "databases-for-mysql":
	// "databases-for-cassandra":
	// "databases-for-enterprisedb"
}

func ConvertMBtoGB(d decimal.Decimal) decimal.Decimal {
	return d.Div(decimal.NewFromInt(1024))
}

// BuildResource builds a schema.Resource from a valid Database struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *Database) BuildResource() *schema.Resource {
	costComponentsFunc, ok := DatabaseCostMap[r.Service]

	if !ok {
		return &schema.Resource{
			Name:        r.Address,
			UsageSchema: DatabaseUsageSchema,
		}
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    DatabaseUsageSchema,
		CostComponents: costComponentsFunc(r),
	}
}
