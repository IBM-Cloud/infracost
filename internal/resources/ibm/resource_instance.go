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
	Name       string
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
	// App Connect
	// Catalog https://cloud.ibm.com/catalog/services/app-connect
	// Pricing https://www.ibm.com/products/app-connect/pricing
	AppConnect_GigabyteTransmittedOutbounds *float64 `infracost_usage:"appconnect_gigabyte_transmitted_outbounds"`
	AppConnect_ThousandRuns                 *float64 `infracost_usage:"appconnect_thousand_runs"`
	AppConnect_VCPUHours                    *float64 `infracost_usage:"appconnect_vcpu_hours"`
	// LogDNA
	// Catalog https://cloud.ibm.com/catalog/services/logdna
	LogDNA_GigabyteMonths *float64 `infracost_usage:"logdna_gigabyte_months"`
	// Activity Tracker
	// Catalog https://cloud.ibm.com/catalog/services/logdnaat
	ActivityTracker_GigabyteMonths *float64 `infracost_usage:"activitytracker_gigabyte_months"`
	// Monitoring (Sysdig)
	// Catalog https://cloud.ibm.com/catalog/services/ibm-cloud-monitoring
	// Pricing https://cloud.ibm.com/docs/monitoring?topic=monitoring-pricing_plans
	Monitoring_NodeHour       *float64 `infracost_usage:"sysdig-monitor_NODE_HOURS"`
	Monitoring_ContainerHour  *float64 `infracost_usage:"sysdig-monitor_CONTAINER_HOURS"`
	Monitoring_APICall        *float64 `infracost_usage:"sysdig-monitor_API_CALL_HOURS"`
	Monitoring_TimeSeriesHour *float64 `infracost_usage:"sysdig-monitor_TIME_SERIES_HOURS"`
	// Continuous Delivery
	// Catalog https://cloud.ibm.com/catalog/services/continuous-delivery
	// Pricing https://cloud.ibm.com/docs/ContinuousDelivery?topic=ContinuousDelivery-limitations_usage&interface=ui
	ContinuousDelivery_AuthorizedUsers *int64 `infracost_usage:"continuousdelivery_authorized_users"`
	// Watson Machine Learning
	// https://dataplatform.cloud.ibm.com/docs/content/wsj/getting-started/wml-plans.html?context=cpdaas
	WML_CUHHours             *float64 `infracost_usage:"wml_capacity_unit_hours"`
	WML_Instance             *float64 `infracost_usage:"wml_instance"`
	WML_MistralLargeOutputRU *float64 `infracost_usage:"wml_mistral_unit_output_ru"`
	WML_MistralLargeInputRU  *float64 `infracost_usage:"wml_mistral_unit_input_ru"`
	WML_TextExtractCat1Pages *float64 `infracost_usage:"wml_text_extract_cat_1_pages"`
	WML_TextExtractCat2Pages *float64 `infracost_usage:"wml_text_extract_cat_2_pages"`
	WML_IBMModelRU           *float64 `infracost_usage:"wml_ibm_model_ru"`
	WML_3rdPartyModelRU      *float64 `infracost_usage:"wml_3rd_party_model_ru"`
	WML_InstructlabDataRU    *float64 `infracost_usage:"wml_instructlab_data_ru"`
	WML_InstructlabTuningRU  *float64 `infracost_usage:"wml_instructlab_tuning_ru"`
	WML_ModelHostingHours    *float64 `infracost_usage:"wml_model_hosting_hours"`
	// Watson Assistant
	WA_Instance *float64 `infracost_usage:"wa_instance"`
	WA_mau      *float64 `infracost_usage:"wa_monthly_active_users"`
	WA_vu       *float64 `infracost_usage:"wa_monthly_voice_users"`
	// Watson Discovery
	WD_Instance     *float64 `infracost_usage:"wd_instance"`
	WD_Documents    *float64 `infracost_usage:"wd_documents"`
	WD_Queries      *float64 `infracost_usage:"wd_queries"`
	WD_CustomModels *float64 `infracost_usage:"wd_custom_models"`
	WD_Collections  *float64 `infracost_usage:"wd_collections"`
	// Security and Compliance Center (SCC)
	SCC_Evaluations *float64 `infracost_usage:"scc_evaluations"`
	// Watson Studio
	WS_CUH *float64 `infracost_usage:"data-science-experience_CAPACITY_UNIT_HOURS"`
	// SCC Workload Protection (Sysdig Secure)
	SCCWP_MulticloudCSPMComputeInstances *float64 `infracost_usage:"sysdig-secure_MULTI_CLOUD_CSPM_COMPUTE_INSTANCES"`
	SCCWP_NodeHours                      *float64 `infracost_usage:"sysdig-secure_NODE_HOURS"`
	SCCWP_VMNodeHours                    *float64 `infracost_usage:"sysdig-secure_VM_NODE_HOUR"`
	// Watsonx.governance
	WGOV_ru     *float64 `infracost_usage:"aiopenscale_RESOURCE_UNITS"`
	WGOV_Models *float64 `infracost_usage:"aiopenscale_MODELS_PER_MONTH"`
	// DNS Services
	DNSServices_CustomResolverExternalQueries *int64   `infracost_usage:"dns-svcs_MILLION_ITEMS_CREXTERNALQUERIES"`
	DNSServices_CustomResolverLocationHours   *float64 `infracost_usage:"dns-svcs_RESOLVERLOCATIONS"`
	DNSServices_CustomResolverLocations       *int64   `infracost_usage:"dns-svcs_qty_custom_resolver_locations"`
	DNSServices_DNSQueries                    *int64   `infracost_usage:"dns-svcs_MILLION_ITEMS"`
	DNSServices_GLBInstanceHours              *float64 `infracost_usage:"dns-svcs_NUMBERGLB"`
	DNSServices_GLBInstances                  *int64   `infracost_usage:"dns-svcs_qty_glb_instances"`
	DNSServices_HealthChecks                  *int64   `infracost_usage:"dns-svcs_NUMBERHEALTHCHECK"`
	DNSServices_PoolHours                     *float64 `infracost_usage:"dns-svcs_NUMBERPOOLS"`
	DNSServices_Pools                         *int64   `infracost_usage:"dns-svcs_qty_pools"`
	DNSServices_Zones                         *int64   `infracost_usage:"dns-svcs_ITEMS"`
	// Event Streams
	EventStreams_CapacityUnitHours            *float64 `infracost_usage:"messagehub_CAPACITY_UNIT_HOURS"`
	EventStreams_CapacityUnitHoursAdditional  *float64 `infracost_usage:"messagehub_CAPACITY_UNIT_HOURS_ADDITIONAL"`
	EventStreams_CapacityUnitHoursMirroring   *float64 `infracost_usage:"messagehub_CAPACITY_UNIT_HOURS_MIRRORING"`
	EventStreams_CapacityUnits                *float64 `infracost_usage:"messagehub_qty_capacity_units"`
	EventStreams_CapacityUnitsAdditional      *float64 `infracost_usage:"messagehub_qty_capacity_units_additional"`
	EventStreams_CapacityUnitsMirroring       *float64 `infracost_usage:"messagehub_qty_capacity_units_mirroring"`
	EventStreams_GigabyteTransmittedOutbounds *float64 `infracost_usage:"messagehub_GIGABYTE_TRANSMITTED_OUTBOUNDS"`
	EventStreams_InstanceHours                *float64 `infracost_usage:"messagehub_INSTANCE_HOURS"`
	EventStreams_Instances                    *float64 `infracost_usage:"messagehub_qty_instances"`
	EventStreams_TerabyteHours                *float64 `infracost_usage:"messagehub_TERABYTE_HOURS"`
	EventStreams_Terabytes                    *float64 `infracost_usage:"messagehub_qty_terabytes"`
	// Event Notifications
	EventNotifications_InboundIngestedEvents *float64 `infracost_usage:"event-notifications_MILLION_INGESTED_EVENTS"`
	Logs_Hours                               *float64 `infracost_usage:"logs_hours"`
}

type ResourceCostComponentsFunc func(*ResourceInstance) []*schema.CostComponent

// PopulateUsage parses the u schema.UsageData into the ResourceInstance.
// It uses the `infracost_usage` struct tags to populate data into the ResourceInstance.
func (r *ResourceInstance) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// ResourceInstanceUsageSchema defines a list which represents the usage schema of ResourceInstance.
var ResourceInstanceUsageSchema = []*schema.UsageItem{
	{Key: "kms_key_versions", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "secretsmanager_instance", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "secretsmanager_active_secrets", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "appid_authentications", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "appid_users", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "appid_advanced_authentications", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "appconnect_gigabyte_transmitted_outbounds", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "appconnect_thousand_runs", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "appconnect_vcpu_hours", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "logdna_gigabyte_months", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "activitytracker_gigabyte_months", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "sysdig-monitor_NODE_HOURS", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "sysdig-monitor_CONTAINER_HOURS", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "sysdig-monitor_API_CALL_HOURS", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "sysdig-monitor_TIME_SERIES_HOURS", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "continuousdelivery_authorized_users", DefaultValue: 0, ValueType: schema.Int64},
	{Key: "wml_capacity_unit_hours", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wml_instance", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wml_mistral_unit_output_ru", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wml_mistral_unit_input_ru", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wml_text_extract_cat_1_pages", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wml_text_extract_cat_2_pages", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wml_ibm_model_ru", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wml_3rd_party_model_ru", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wml_instructlab_data_ru", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wsl_instructlab_tuning_ru", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wml_model_hosting_hours", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wa_instance", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wa_monthly_active_users", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wa_monthly_voice_users", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wd_instance", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wd_documents", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wd_queries", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wd_custom_models", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "wd_collections", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "scc_evaluations", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "data-science-experience_CAPACITY_UNIT_HOURS", DefaultValue: 1, ValueType: schema.Float64},
	{Key: "sysdig-secure_MULTI_CLOUD_CSPM_COMPUTE_INSTANCES", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "sysdig-secure_NODE_HOURS", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "sysdig-secure_VM_NODE_HOUR", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "aiopenscale_RESOURCE_UNITS", DefaultValue: 1, ValueType: schema.Float64},
	{Key: "aiopenscale_MODELS_PER_MONTH", DefaultValue: 1, ValueType: schema.Float64},
	{Key: "dns-svcs_ITEMS", DefaultValue: 1, ValueType: schema.Int64},
	{Key: "dns-svcs_MILLION_ITEMS_CREXTERNALQUERIES", DefaultValue: 1, ValueType: schema.Int64},
	{Key: "dns-svcs_MILLION_ITEMS", DefaultValue: 1, ValueType: schema.Int64},
	{Key: "dns-svcs_NUMBERGLB", DefaultValue: 1, ValueType: schema.Float64},
	{Key: "dns-svcs_NUMBERHEALTHCHECK", DefaultValue: 1, ValueType: schema.Int64},
	{Key: "dns-svcs_NUMBERPOOLS", DefaultValue: 1, ValueType: schema.Float64},
	{Key: "dns-svcs_RESOLVERLOCATIONS", DefaultValue: 1, ValueType: schema.Float64},
	{Key: "dns-svcs_qty_custom_resolver_locations", DefaultValue: 1, ValueType: schema.Int64},
	{Key: "dns-svcs_qty_glb_instances", DefaultValue: 1, ValueType: schema.Int64},
	{Key: "dns-svcs_qty_pools", DefaultValue: 1, ValueType: schema.Int64},
	{Key: "messagehub_CAPACITY_UNIT_HOURS", DefaultValue: 1, ValueType: schema.Float64},
	{Key: "messagehub_CAPACITY_UNIT_HOURS_ADDITIONAL", DefaultValue: 1, ValueType: schema.Float64},
	{Key: "messagehub_CAPACITY_UNIT_HOURS_MIRRORING", DefaultValue: 1, ValueType: schema.Float64},
	{Key: "messagehub_qty_capacity_units", DefaultValue: 1, ValueType: schema.Float64},
	{Key: "messagehub_qty_capacity_units_additional", DefaultValue: 1, ValueType: schema.Float64},
	{Key: "messagehub_qty_capacity_units_mirroring", DefaultValue: 1, ValueType: schema.Float64},
	{Key: "messagehub_GIGABYTE_TRANSMITTED_OUTBOUNDS", DefaultValue: 1, ValueType: schema.Float64},
	{Key: "messagehub_INSTANCE_HOURS", DefaultValue: 1, ValueType: schema.Float64},
	{Key: "messagehub_qty_instances", DefaultValue: 1, ValueType: schema.Float64},
	{Key: "messagehub_TERABYTE_HOURS", DefaultValue: 1, ValueType: schema.Float64},
	{Key: "messagehub_qty_terabytes", DefaultValue: 1, ValueType: schema.Float64},
	{Key: "event-notifications_MILLION_INGESTED_EVENTS", DefaultValue: 0, ValueType: schema.Float64},
}

var ResourceInstanceCostMap map[string]ResourceCostComponentsFunc = map[string]ResourceCostComponentsFunc{
	"aiopenscale":             GetWGOVCostComponents,
	"appconnect":              GetAppConnectCostComponents,
	"appid":                   GetAppIDCostComponents,
	"compliance":              GetSCCCostComponents,
	"continuous-delivery":     GetContinuousDeliveryCostComponenets,
	"conversation":            GetWACostComponents,
	"data-science-experience": GetWSCostComponents,
	"discovery":               GetWDCostComponents,
	"dns-svcs":                GetDNSServicesCostComponents,
	"event-notifications":     GetEventNotificationsCostComponents,
	"kms":                     GetKMSCostComponents,
	"logdna":                  GetLogDNACostComponents,
	"logdnaat":                GetActivityTrackerCostComponents,
	"messagehub":              GetEventStreamsCostComponents,
	"pm-20":                   GetWMLCostComponents,
	"power-iaas":              GetPowerCostComponents,
	"secrets-manager":         GetSecretsManagerCostComponents,
	"sysdig-monitor":          GetSysdigCostComponenets,
	"sysdig-secure":           GetSCCWPCostComponents,
	"logs":                    GetLogsCostComponents,
}

func KMSKeyVersionsFreeCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.KMS_KeyVersions != nil {
		q = decimalPtr(decimal.NewFromInt(*r.KMS_KeyVersions))
		if q.GreaterThan(decimal.NewFromInt(5)) {
			q = decimalPtr(decimal.NewFromInt(5))
		}
	}
	costComponent := schema.CostComponent{
		Name:            "Key versions free allowance (first 5 Key Versions)",
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
		Name:            "Key versions",
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
		Name:            "Instance",
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
		Name:            "Active Secrets",
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
	if r.Plan == "trial" {
		costComponent := schema.CostComponent{
			Name:            "Trial plan",
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
			ProductFilter: &schema.ProductFilter{
				VendorName: strPtr("ibm"),
				Region:     strPtr(r.Location),
				Service:    &r.Service,
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "planName", Value: &r.Plan},
				},
			},
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
		return []*schema.CostComponent{
			&costComponent,
		}
	} else if r.Plan == "standard" {
		return []*schema.CostComponent{
			SecretsManagerInstanceCostComponent(r),
			SecretsManagerActiveSecretsCostComponent(r),
		}
	} else {
		costComponent := schema.CostComponent{
			Name:            fmt.Sprintf("Plan %s not found", r.Plan),
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
			ProductFilter: &schema.ProductFilter{
				VendorName: strPtr("ibm"),
				Region:     strPtr(r.Location),
				Service:    &r.Service,
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "planName", Value: &r.Plan},
				},
			},
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
		return []*schema.CostComponent{
			&costComponent,
		}
	}
}

func GetPowerCostComponents(r *ResourceInstance) []*schema.CostComponent {

	costComponent := schema.CostComponent{
		Name:            "Workspace for Power Virtual Server",
		Unit:            "Workspace",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
	}
	costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
	return []*schema.CostComponent{
		&costComponent,
	}
}

func AppIDUserCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.AppID_Users != nil {
		q = decimalPtr(decimal.NewFromInt(*r.AppID_Users))
	}
	costComponent := schema.CostComponent{
		Name:            "Users",
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
		Name:            "Authentications",
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
		Name:            "Advanced Authentications",
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
		costComponent := schema.CostComponent{
			Name:            fmt.Sprintf("Plan %s not found", r.Plan),
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
			ProductFilter: &schema.ProductFilter{
				VendorName: strPtr("ibm"),
				Region:     strPtr(r.Location),
				Service:    &r.Service,
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "planName", Value: &r.Plan},
				},
			},
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
		return []*schema.CostComponent{
			&costComponent,
		}
	}
}

func AppConnectFlowsRunsCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.AppConnect_ThousandRuns != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.AppConnect_ThousandRuns))
	}
	return &schema.CostComponent{
		Name:            "Flow runs",
		Unit:            "1k runs",
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
			Unit: strPtr("ITEMS_PER_MONTH"),
		},
	}
}

func AppConnectEgressCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.AppConnect_GigabyteTransmittedOutbounds != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.AppConnect_GigabyteTransmittedOutbounds))
	}
	return &schema.CostComponent{
		Name:            "Gigabytes transmitted outbounds",
		Unit:            "GB",
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
			Unit: strPtr("GIGABYTE_TRANSMITTED_OUTBOUND"),
		},
	}
}

func AppConnectCpuCostComponent(r *ResourceInstance) *schema.CostComponent {
	var q *decimal.Decimal
	if r.AppConnect_VCPUHours != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.AppConnect_VCPUHours))
	}
	return &schema.CostComponent{
		Name:            "VCPU",
		Unit:            "VCPU hours",
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
			Unit: strPtr("VIRTUAL_PROCESSOR_CORE_HOURS"),
		},
	}
}

func GetAppConnectCostComponents(r *ResourceInstance) []*schema.CostComponent {
	if r.Plan == "appconnectplanprofessional" {
		return []*schema.CostComponent{
			AppConnectFlowsRunsCostComponent(r),
			AppConnectEgressCostComponent(r),
		}
	} else if r.Plan == "appconnectplanenterprise" {
		return []*schema.CostComponent{
			AppConnectFlowsRunsCostComponent(r),
			AppConnectEgressCostComponent(r),
			AppConnectCpuCostComponent(r),
		}
	} else if r.Plan == "lite" {
		costComponent := schema.CostComponent{
			Name:            "Lite plan",
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
			ProductFilter: &schema.ProductFilter{
				VendorName: strPtr("ibm"),
				Region:     strPtr(r.Location),
				Service:    &r.Service,
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "planName", Value: &r.Plan},
				},
			},
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
		return []*schema.CostComponent{
			&costComponent,
		}
	} else {
		costComponent := schema.CostComponent{
			Name:            fmt.Sprintf("Plan %s not found", r.Plan),
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
			ProductFilter: &schema.ProductFilter{
				VendorName: strPtr("ibm"),
				Region:     strPtr(r.Location),
				Service:    &r.Service,
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "planName", Value: &r.Plan},
				},
			},
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
		return []*schema.CostComponent{
			&costComponent,
		}
	}
}

func GetLogDNACostComponents(r *ResourceInstance) []*schema.CostComponent {
	var q *decimal.Decimal
	if r.LogDNA_GigabyteMonths != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.LogDNA_GigabyteMonths))
	}
	if r.Plan == "lite" {
		costComponent := schema.CostComponent{
			Name:            "Lite plan",
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
			ProductFilter: &schema.ProductFilter{
				VendorName: strPtr("ibm"),
				Region:     strPtr(r.Location),
				Service:    &r.Service,
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "planName", Value: &r.Plan},
				},
			},
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
		return []*schema.CostComponent{
			&costComponent,
		}
	} else {
		return []*schema.CostComponent{{
			Name:            fmt.Sprintf("Gigabyte Months (%s)", r.Plan),
			Unit:            "Gigabyte Months",
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
				Unit: strPtr("GIGABYTE_MONTHS"),
			},
		}}
	}
}

func GetActivityTrackerCostComponents(r *ResourceInstance) []*schema.CostComponent {
	var q *decimal.Decimal
	if r.ActivityTracker_GigabyteMonths != nil {
		q = decimalPtr(decimal.NewFromFloat(*r.ActivityTracker_GigabyteMonths))
	}
	if r.Plan == "lite" {
		costComponent := &schema.CostComponent{
			Name:            "Lite plan",
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
			ProductFilter: &schema.ProductFilter{
				VendorName: strPtr("ibm"),
				Region:     strPtr(r.Location),
				Service:    &r.Service,
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "planName", Value: &r.Plan},
				},
			},
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
		return []*schema.CostComponent{costComponent}
	} else {
		return []*schema.CostComponent{{
			Name:            "Gigabyte Months",
			Unit:            "Gigabyte Months",
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
				Unit: strPtr("GIGABYTE_MONTHS"),
			},
		}}
	}
}

func GetContinuousDeliveryCostComponenets(r *ResourceInstance) []*schema.CostComponent {
	var q *decimal.Decimal
	if r.ContinuousDelivery_AuthorizedUsers != nil {
		q = decimalPtr(decimal.NewFromInt(*r.ContinuousDelivery_AuthorizedUsers))
	}
	if r.Plan == "lite" {

		costComponent := &schema.CostComponent{
			Name:            "Lite plan",
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
			ProductFilter: &schema.ProductFilter{
				VendorName: strPtr("ibm"),
				Region:     strPtr(r.Location),
				Service:    &r.Service,
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "planName", Value: &r.Plan},
				},
			},
		}
		costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
		return []*schema.CostComponent{costComponent}
	} else {
		return []*schema.CostComponent{{
			Name:            "Authorized Users",
			Unit:            "Authorized Users",
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
				Unit: strPtr("AUTHORIZED_USERS_PER_MONTH"),
			},
		}}
	}
}

// BuildResource builds a schema.Resource from a valid ResourceInstance struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *ResourceInstance) BuildResource() *schema.Resource {
	costComponentsFunc, ok := ResourceInstanceCostMap[r.Service]

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
