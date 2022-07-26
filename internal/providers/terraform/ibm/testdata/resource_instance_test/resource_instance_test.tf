terraform {
  required_providers {
    ibm = {
      source = "IBM-Cloud/ibm"
      version = "~> 1.40.0"
    }
  }
}

provider "ibm" {
    region = "us-south"
}

resource "ibm_resource_instance" "resource_instance_kms" {
  name              = "test"
  service           = "kms"
  plan              = "tiered-pricing"
  location          = "us-south"
  resource_group_id = "default"
}

resource "ibm_resource_instance" "resource_instance_secrets_manager" {
  name              = "test"
  service           = "secrets-manager"
  plan              = "standard"
  location          = "us-south"
  resource_group_id = "default"
}