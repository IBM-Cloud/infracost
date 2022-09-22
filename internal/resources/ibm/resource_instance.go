package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

// ResourceInstance struct represents a resource instance
//
// This terraform resource is opaque and can handle a number of services, provided with the right parameters
type ResourceInstance struct {
	Address    string
	Service    string
	Plan       string
	Location   string
	Parameters gjson.Result

	// KMS
	// Catalog Link: https://cloud.ibm.com/catalog/services/key-protect
	// Pricing Link: https://cloud.ibm.com/docs/key-protect?topic=key-protect-pricing-plan&interface=ui
	KMS_KeyVersions *int64 `infracost_usage:"kms_key_versions"`
	// Secrets Manager
	// Catalog link: https://cloud.ibm.com/catalog/services/secrets-manager
	SecretsManager_Instance      *int64 `infracost_usage:"secretsmanager_instance"`
	SecretsManager_ActiveSecrets *int64 `infracost_usage:"secretsmanager_active_secrets"`
	// App ID
	// Catalog https://cloud.ibm.com/catalog/services/app-id
	// Pricing https://cloud.ibm.com/docs/appid?topic=appid-pricing
	AppID_Authentications         *int64 `infracost_usage:"appid_authentications"`
	AppID_Users                   *int64 `infracost_usage:"appid_users"`
	AppID_AdvancedAuthentications *int64 `infracost_usage:"appid_advanced_authentications"`
}

// ResourceInstanceUsageSchema defines a list which represents the usage schema of ResourceInstance.
var ResourceInstanceUsageSchema = []*schema.UsageItem{
	{Key: "kms_key_versions", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "secretsmanager_instance", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "secretsmanager_active_secrets", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "appid_authentications", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "appid_users", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "appid_advanced_authentications", DefaultValue: 0, ValueType: schema.Int64},
}

var ResourceCostMap map[string]ResourceCostComponentsFunc = map[string]ResourceCostComponentsFunc{
	"kms":             GetKMSCostComponents,
	"secrets-manager": GetSecretsManagerCostComponents,
	"appid":           GetAppIDCostComponents,
}

// PopulateUsage parses the u schema.UsageData into the ResourceInstance.
// It uses the `infracost_usage` struct tags to populate data into the ResourceInstance.
func (r *ResourceInstance) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

type ResourceCostComponentsFunc func(*ResourceInstance) []*schema.CostComponent

func KMSKeyVersionsFreeCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.KMS_KeyVersions != nil {
		q = decimalPtr(decimal.NewFromInt(*r.KMS_KeyVersions))
		if q.GreaterThan(decimal.NewFromInt(5)) {
			q = decimalPtr(decimal.NewFromInt(5))
		}
	}
	costComponent := schema.CostComponent{
		Name:            fmt.Sprintf("Key versions free allowance (first 5 Key Versions)"),
		Unit:            "Key Versions",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    strPtr("kms"),
		},
	}
	costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
	return &costComponent
}

func KMSKeyVersionCostComponents(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.KMS_KeyVersions != nil {
		q = decimalPtr(decimal.NewFromInt(*r.KMS_KeyVersions))
		if q.LessThanOrEqual(decimal.NewFromInt(5)) {
			q = decimalPtr(decimal.NewFromInt(0))
		} else {
			q = decimalPtr(q.Sub(decimal.NewFromInt(5)))
		}
	}
	costComponent := schema.CostComponent{
		Name:            fmt.Sprintf("Key versions"),
		Unit:            "Key Versions",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
	}
	return &costComponent
}

func GetKMSCostComponents(r *ResourceInstance) []*schema.CostComponent {
	return []*schema.CostComponent{
		KMSKeyVersionsFreeCostComponent(r),
		KMSKeyVersionCostComponents(r),
	}
}

func SecretsManagerInstanceCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.SecretsManager_Instance != nil {
		q = decimalPtr(decimal.NewFromInt(*r.SecretsManager_Instance))
	}
	costComponent := schema.CostComponent{
		Name:            fmt.Sprintf("Instance"),
		Unit:            "Instance",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("INSTANCES"),
		},
	}
	return &costComponent
}

func SecretsManagerActiveSecretsCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.SecretsManager_ActiveSecrets != nil {
		q = decimalPtr(decimal.NewFromInt(*r.SecretsManager_ActiveSecrets))
	}
	costComponent := schema.CostComponent{
		Name:            fmt.Sprintf("Active Secrets"),
		Unit:            "Secrets",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("ACTIVE_SECRETS"),
		},
	}
	return &costComponent
}

func GetSecretsManagerCostComponents(r *ResourceInstance) []*schema.CostComponent {
	if r.Plan == "standard" {
		return []*schema.CostComponent{
			SecretsManagerInstanceCostComponent(r),
			SecretsManagerActiveSecretsCostComponent(r),
		}
	} else {
		costComponent := *&schema.CostComponent{
			Name: fmt.Sprintf("Plan: %s", r.Plan),
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
		return []*schema.CostComponent{
			&costComponent,
		}
	}
}

func AppIDUserCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.AppID_Users != nil {
		q = decimalPtr(decimal.NewFromInt(*r.AppID_Users))
	}
	costComponent := schema.CostComponent{
		Name:            fmt.Sprintf("Users"),
		Unit:            "Users",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("USERS_PER_MONTH"),
		},
	}
	return &costComponent
}

func AppIDAuthenticationCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.AppID_Authentications != nil {
		q = decimalPtr(decimal.NewFromInt(*r.AppID_Authentications))
	}
	costComponent := schema.CostComponent{
		Name:            fmt.Sprintf("Authentications"),
		Unit:            "Authentications",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("AUTHENTICATIONS_PER_MONTH"),
		},
	}
	return &costComponent
}

func AppIDAdvancedAuthenticationCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.AppID_AdvancedAuthentications != nil {
		q = decimalPtr(decimal.NewFromInt(*r.AppID_AdvancedAuthentications))
	}
	costComponent := schema.CostComponent{
		Name:            fmt.Sprintf("Advanced Authentications"),
		Unit:            "Authentications",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr("ADVANCED_AUTHENTICATIONS_PER_MONTH"),
		},
	}
	return &costComponent
}

func GetAppIDCostComponents(r *ResourceInstance) []*schema.CostComponent {
	if r.Plan == "graduated-tier" {
		return []*schema.CostComponent{
			AppIDUserCostComponent(r),
			AppIDAuthenticationCostComponent(r),
			AppIDAdvancedAuthenticationCostComponent(r),
		}
	} else {
		costComponent := *&schema.CostComponent{
			Name: fmt.Sprintf("Plan: %s", r.Plan),
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
		return []*schema.CostComponent{
			&costComponent,
		}
	}
}

// BuildResource builds a schema.Resource from a valid ResourceInstance struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *ResourceInstance) BuildResource() *schema.Resource {
	costComponentsFunc, ok := ResourceCostMap[r.Service]

	if !ok {
		return &schema.Resource{
			Name:        r.Address,
			UsageSchema: ResourceInstanceUsageSchema,
		}
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    ResourceInstanceUsageSchema,
		CostComponents: costComponentsFunc(r),
	}
}
