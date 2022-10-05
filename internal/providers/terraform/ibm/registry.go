package ibm

import "github.com/infracost/infracost/internal/schema"

var ResourceRegistry []*schema.RegistryItem = []*schema.RegistryItem{
	getIsInstanceRegistryItem(),
	getIbmIsVpcRegistryItem(),
	getIbmCosBucketRegistryItem(),
	getIsFloatingIpRegistryItem(),
	getIsFlowLogRegistryItem(),
	getContainerVpcWorkerPoolRegistryItem(),
	getContainerVpcClusterRegistryItem(),
	getResourceInstanceRegistryItem(),
	getIsVolumeRegistryItem(),
	getIsVpnGatewayRegistryItem(),
	getTgGatewayRegistryItem(),
	getCloudantRegistryItem(),
	getPiInstanceRegistryItem(),
	getIsLbRegistryItem(),
	getIsPublicGatewayRegistryItem(),
}

// FreeResources grouped alphabetically
var FreeResources = []string{
	"ibm_atracker_route",
	"ibm_atracker_target",
	"ibm_iam_access_group",
	"ibm_iam_access_group_dynamic_rule",
	"ibm_iam_access_group_members",
	"ibm_iam_access_group_policy",
	"ibm_iam_account_settings",
	"ibm_iam_authorization_policy",
	"ibm_is_lb_listener",
	"ibm_is_lb_pool",
	"ibm_is_lb_pool_member",
	"ibm_is_network_acl",
	"ibm_is_security_group",
	"ibm_is_security_group_rule",
	"ibm_is_ssh_key",
	"ibm_is_subnet",
	"ibm_is_subnet_reserved_ip",
	"ibm_is_virtual_endpoint_gateway",
	"ibm_is_virtual_endpoint_gateway_ip",
	"ibm_is_vpc_address_prefix",
	"ibm_is_vpn_gateway_connection",
	"ibm_kms_key",
	"ibm_kms_key_rings",
	"ibm_pi_cloud_connection",
	"ibm_pi_cloud_connection_network_attach",
	"ibm_pi_console_language",
	"ibm_pi_dhcp",
	"ibm_pi_ike_policy",
	"ibm_pi_instance_action",
	"ibm_pi_ipsec_policy",
	"ibm_pi_key",
	"ibm_pi_network",
	"ibm_pi_network_port",
	"ibm_pi_network_port_attach",
	"ibm_pi_placement_group",
	"ibm_pi_shared_processor_pool",
	"ibm_pi_vpn_connection",
	"ibm_resource_group",
	"ibm_resource_key",
	"ibm_tg_connection",
}

var UsageOnlyResources = []string{
	"",
}
