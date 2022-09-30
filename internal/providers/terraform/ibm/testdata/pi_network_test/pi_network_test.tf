
terraform {
  required_providers {
    ibm = {
      source  = "IBM-Cloud/ibm"
      version = "~> 1.40.0"
    }
  }
}

provider "ibm" {
  region = "us-south"
  zone   = "dal12"
}

locals {
  service_type = "power-iaas"
}

resource "ibm_resource_group" "resource_group" {
  name = "default"
}

resource "ibm_resource_instance" "powervs_service" {
  name              = "Power instance"
  service           = local.service_type
  plan              = "power-virtual-server-group"
  location          = "us-south"
  resource_group_id = ibm_resource_group.resource_group.id
}

resource "ibm_pi_network" "public_vlan_power_network" {
  count                = 1
  pi_network_name      = "power-network"
  pi_cloud_instance_id = ibm_resource_instance.powervs_service.guid
  pi_network_type      = "pub-vlan"
}
