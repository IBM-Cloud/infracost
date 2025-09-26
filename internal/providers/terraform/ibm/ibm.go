package ibm

import (
	"encoding/json"
	"fmt"

	"github.com/infracost/infracost/internal/schema"
	"github.com/tidwall/gjson"
)

var DefaultProviderRegion = "us-south"

func GetDefaultRefIDFunc(d *schema.ResourceData) []string {

	defaultRefs := []string{d.Get("id").String()}

	if d.Get("self_link").Exists() {
		defaultRefs = append(defaultRefs, d.Get("self_link").String())
	}

	return defaultRefs
}

func DefaultCloudResourceIDFunc(d *schema.ResourceData) []string {
	return []string{}
}

func GetSpecialContext(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{}
}

func GetResourceRegion(resourceType string, v gjson.Result) string {
	return ""
}

func ParseTags(resourceType string, v gjson.Result) map[string]string {
	tags := make(map[string]string)
	for k, v := range v.Get("labels").Map() {
		tags[k] = v.String()
	}
	return tags
}

type catalogMetadata struct {
	serviceId      string
	childResources []string
	configuration  map[string]any
	pricingUrl     string
}

// Map between terraform type and global catalog id. For ibm_resource_instance, the service
// field already matches the global catalog id, so they do not need to be mapped. eg: "kms"
var globalCatalogServiceId = map[string]catalogMetadata{
	"aiopenscale":                   {"2ad019f3-0fd6-4c25-966d-f3952481a870", []string{}, nil, "https://cloud.ibm.com/catalog/services/watsonxgovernance"},
	"appconnect":                    {"96a0ebf2-2a02-4e32-815f-7c09a1268c78", []string{}, nil, "https://www.ibm.com/products/app-connect/pricing"},
	"appid":                         {"AdvancedMobileAccess-d6aece47-d840-45b0-8ab9-ad15354deeea", []string{}, nil, "https://cloud.ibm.com/catalog/services/appid"},
	"apprapp":                       {"apprapp-d6aece47-d840-45b0-8ab9-ad15354deeea", []string{}, nil, "https://cloud.ibm.com/catalog/services/app-configuration"},
	"cloud-object-storage":          {"dff97f5c-bc5e-4455-b470-411c3edbe49c", []string{}, nil, "https://cloud.ibm.com/objectstorage/create#pricing"},
	"compliance":                    {"compliance", []string{}, nil, "https://cloud.ibm.com/catalog/services/security-and-compliance-center"},
	"continuous-delivery":           {"59b735ee-5938-4ebd-a6b2-541aef2d1f68", []string{}, nil, "https://cloud.ibm.com/catalog/services/continuous-delivery"},
	"conversation":                  {"7045626d-55e3-4418-be11-683a26dbc1e5", []string{}, nil, "https://cloud.ibm.com/catalog/services/watsonx-assistant"},
	"data-science-experience":       {"39ba9d4c-b1c5-4cc3-a163-38b580121e01", []string{}, nil, "https://cloud.ibm.com/catalog/services/watson-studio"},
	"databases-for-elasticsearch":   {"databases-for-elasticsearch", []string{}, nil, "https://cloud.ibm.com/databases/databases-for-elasticsearch/create"},
	"databases-for-postgresql":      {"databases-for-postgresql", []string{}, nil, "https://cloud.ibm.com/databases/databases-for-postgresql/create"},
	"discovery":                     {"76b7bf22-b443-47db-b3db-066ba2988f47", []string{}, nil, "https://cloud.ibm.com/catalog/services/watson-discovery"},
	"dns-svcs":                      {"b4ed8a30-936f-11e9-b289-1d079699cbe5", []string{}, nil, "https://cloud.ibm.com/catalog/services/dns-services"},
	"event-notifications":           {"ecdb4690-c2d8-11eb-bff1-4f7b9d2dfe41", []string{}, nil, "https://cloud.ibm.com/catalog/services/event-notifications"},
	"ibm_cloudant":                  {"Cloudant", []string{}, nil, "https://cloud.ibm.com/catalog/services/cloudant"},
	"ibm_code_engine_app":           {"2ad2fdd0-bba5-11ea-8966-5d6402fed1c7", []string{}, nil, "https://cloud.ibm.com/docs/codeengine?topic=codeengine-pricing"},
	"ibm_code_engine_build":         {"2ad2fdd0-bba5-11ea-8966-5d6402fed1c7", []string{}, nil, "https://cloud.ibm.com/docs/codeengine?topic=codeengine-pricing"},
	"ibm_code_engine_function":      {"2ad2fdd0-bba5-11ea-8966-5d6402fed1c7", []string{}, nil, "https://cloud.ibm.com/docs/codeengine?topic=codeengine-pricing"},
	"ibm_code_engine_job":           {"2ad2fdd0-bba5-11ea-8966-5d6402fed1c7", []string{}, nil, "https://cloud.ibm.com/docs/codeengine?topic=codeengine-pricing"},
	"ibm_container_vpc_cluster":     {"containers-kubernetes", []string{"ibm_container_vpc_worker_pool"}, nil, "https://cloud.ibm.com/kubernetes/catalog/about#pricing"},
	"ibm_container_vpc_worker_pool": {"Worker Pool", []string{}, nil, "https://cloud.ibm.com/kubernetes/catalog/about#pricing"},
	"ibm_cos_bucket":                {"Object Storage Bucket", []string{}, nil, "https://cloud.ibm.com/objectstorage/create#pricing"},
	"ibm_is_floating_ip":            {"is.floating-ip", []string{}, nil, "https://cloud.ibm.com/vpc-ext/provision/vs"},
	"ibm_is_flow_log":               {"is.flow-log-collector", []string{}, nil, "https://cloud.ibm.com/vpc-ext/provision/flowLog"},
	"ibm_is_instance":               {"is.instance", []string{"ibm_is_ssh_key", "ibm_is_floating_ip"}, nil, "https://cloud.ibm.com/vpc-ext/provision/vs"},
	"ibm_is_lb":                     {"is.load-balancer", []string{}, nil, "https://cloud.ibm.com/vpc-ext/provision/loadBalancer"},
	"ibm_is_share":                  {"is.share", []string{}, nil, "https://cloud.ibm.com/docs/vpc?topic=vpc-file-storage-vpc-faqs&interface=ui#faq-fs-billing"},
	"ibm_is_volume":                 {"is.volume", []string{}, nil, "https://cloud.ibm.com/vpc-ext/provision/storage"},
	"ibm_is_vpc":                    {"is.vpc", []string{"ibm_is_flow_log", "ibm_is_share"}, nil, "https://cloud.ibm.com/vpc-ext/provision/vpc"},
	"ibm_is_vpn_gateway":            {"is.vpn", []string{}, nil, "https://cloud.ibm.com/vpc-ext/provision/vpngateway"},
	"ibm_is_vpn_server":             {"is.vpn-server", []string{}, nil, "https://cloud.ibm.com/vpc-ext/provision/vpnserver"},
	"ibm_pi_instance":               {"Power Systems Virtual Server", []string{}, nil, "https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-pricing-virtual-server"},
	"ibm_pi_volume":                 {"Power Systems Storage Volume", []string{}, nil, "https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-pricing-virtual-server#storage-type"},
	"ibm_tg_gateway":                {"f38a4da0-c353-11e9-83b6-a36a57a97a06", []string{}, nil, "https://cloud.ibm.com/interconnectivity/transit/provision"},
	"kms":                           {"ee41347f-b18e-4ca6-bf80-b5467c63f9a6", []string{}, nil, "https://cloud.ibm.com/catalog/services/key-protect"},
	"logdna":                        {"e13e1860-959c-11e8-871e-ad157af61ad7", []string{}, nil, "https://cloud.ibm.com/catalog/services/logdna"},
	"logdnaat":                      {"dcc46a60-e13b-11e8-a015-757410dab16b", []string{}, nil, "https://cloud.ibm.com/catalog/services/logdnaat"},
	"logs":                          {"cd515180-d78a-11ec-b396-db7d306c4f73", []string{}, nil, "https://cloud.ibm.com/catalog/services/cloud-logs"},
	"messagehub":                    {"6a7f4e38-f218-48ef-9dd2-df408747568e", []string{}, nil, "https://cloud.ibm.com/eventstreams-provisioning/6a7f4e38-f218-48ef-9dd2-df408747568e/create"},
	"pm-20":                         {"51c53b72-918f-4869-b834-2d99eb28422a", []string{}, nil, "https://cloud.ibm.com/catalog/services/watson-machine-learning"},
	"power-iaas":                    {"abd259f0-9990-11e8-acc8-b9f54a8f1661", []string{}, nil, "https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-pricing-virtual-server"},
	"roks":                          {"containers.kubernetes.cluster.roks", []string{}, nil, "https://cloud.ibm.com/kubernetes/catalog/about?platformType=openshift"},
	"secrets-manager":               {"ebc0cdb0-af2a-11ea-98c7-29e5db822649", []string{}, nil, "https://cloud.ibm.com/catalog/services/secrets-manager"},
	"sysdig-monitor":                {"090c2c10-8c38-11e8-bec2-493df9c49eb8", []string{}, nil, "https://cloud.ibm.com/observe/catalog/ibm-cloud-monitoring"},
	"sysdig-secure":                 {"e831e900-82d6-11ec-95c5-c12c5a5d9687", []string{}, nil, "https://cloud.ibm.com/workload-protection/catalog/security-and-compliance-center-workload-protection"},
	"watsonx-orchestrate":           {"b69f78c0-11d7-11ef-9bdf-c92eb40d1838", []string{}, nil, "https://cloud.ibm.com/catalog/services/watsonx-orchestrate"},
	"wx":                            {"51c53b72-918f-4869-b834-2d99eb28422a", []string{}, nil, "https://cloud.ibm.com/watsonx/overview"},
}

func SetCatalogMetadata(d *schema.ResourceData, resourceType string, config map[string]any) {
	metadata := make(map[string]gjson.Result)
	var properties gjson.Result
	var serviceId string = resourceType
	var childResources []string
	var pricingUrl string

	catalogEntry, isPresent := globalCatalogServiceId[resourceType]
	if isPresent {
		serviceId = catalogEntry.serviceId
		pricingUrl = catalogEntry.pricingUrl
		childResources = catalogEntry.childResources
	}

	configString, err := json.Marshal(config)
	if err != nil {
		configString = []byte("{}")
	}

	if len(childResources) > 0 {
		childResourcesString, err := json.Marshal(childResources)
		if err != nil {
			childResourcesString = []byte("[]")
		}

		properties = gjson.Result{
			Type: gjson.JSON,
			Raw:  fmt.Sprintf(`{"serviceId": "%s" , "pricingUrl": "%s", "childResources": %s, "configuration": %s}`, serviceId, pricingUrl, childResourcesString, configString),
		}
	} else {
		properties = gjson.Result{
			Type: gjson.JSON,
			Raw:  fmt.Sprintf(`{"serviceId": "%s", "pricingUrl": "%s", "configuration": %s}`, serviceId, pricingUrl, configString),
		}
	}

	metadata["catalog"] = properties
	d.Metadata = metadata
}
