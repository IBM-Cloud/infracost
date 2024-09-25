
terraform {
  required_providers {
    ibm = {
      source = "IBM-Cloud/ibm"
    }
  }
}

provider "ibm" {
  region = "us-south"
}

resource "ibm_is_vpc" "vpc" {
  name = "test-vpc"
}

resource "ibm_is_subnet" "subnet" {
  name            = "test-subnet"
  vpc             = ibm_is_vpc.vpc.id
  zone            = "us-south-1"
  ipv4_cidr_block = "10.240.0.0/24"
}

resource "ibm_is_ssh_key" "ssh_key" {
  name       = "test-ssh"
  public_key = file("~/.ssh/id_rsa.pub")
}

resource "ibm_is_instance" "vsi" {
  name    = "vsi-instance"
  image   = "r006-f137ea64-0d27-4d81-afe0-353fd0557e81"
  profile = "cx3d-2x5"
  keys    = [ibm_is_ssh_key.ssh_key.id]
  vpc     = ibm_is_vpc.vpc.id
  zone    = "us-south-1"
  primary_network_interface {
    subnet = ibm_is_subnet.subnet.id
  }
  network_interfaces {
    name   = "eth1"
    subnet = ibm_is_subnet.subnet.id
  }
}
