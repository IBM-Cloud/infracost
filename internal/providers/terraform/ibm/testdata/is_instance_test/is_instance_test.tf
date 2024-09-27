
terraform {
  required_providers {
    ibm = {
      source = "IBM-Cloud/ibm"
    }
    random = {
      source = "hashicorp/random"
    }
  }
}

provider "ibm" {
  region = "us-south"
}

# Access random string generated with random_string.unique_identifier.result
resource "random_string" "unique_identifier" {
  length  = 6
  special = false
  upper   = false
}

resource "ibm_resource_group" "resource_group" {
  name = "${random_string.unique_identifier.result}-rg"
}

resource "ibm_is_vpc" "vpc" {
  name           = "${random_string.unique_identifier.result}-vpc"
  resource_group = ibm_resource_group.resource_group.id
}

resource "ibm_is_subnet" "subnet" {
  name            = "${random_string.unique_identifier.result}-subnet"
  ipv4_cidr_block = "10.240.0.0/24"
  resource_group  = ibm_resource_group.resource_group.id
  vpc             = ibm_is_vpc.vpc.id
  zone            = "us-south-1"
}

resource "ibm_is_ssh_key" "ssh_key" {
  name           = "${random_string.unique_identifier.result}-ssh"
  public_key     = file("~/.ssh/id_ed25519.pub")
  resource_group = ibm_resource_group.resource_group.id
  type           = "ed25519"
}

resource "ibm_is_instance" "vsi" {
  name           = "${random_string.unique_identifier.result}-vsi-instance"
  image          = "r006-f137ea64-0d27-4d81-afe0-353fd0557e81"
  keys           = [ibm_is_ssh_key.ssh_key.id]
  profile        = "cx3d-2x5"
  resource_group = ibm_resource_group.resource_group.id
  vpc            = ibm_is_vpc.vpc.id
  zone           = "us-south-1"
  primary_network_interface {
    subnet = ibm_is_subnet.subnet.id
  }
  network_interfaces {
    name   = "eth1"
    subnet = ibm_is_subnet.subnet.id
  }
}

locals {
  profiles = [

  ]
}
