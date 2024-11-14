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

resource "ibm_resource_group" "test_group" {
  name = "test-resource-group"
}

resource "ibm_code_engine_project" "ce_project" {
  name              = "ce_project"
  resource_group_id = ibm_resource_group.test_group.id
}

resource "ibm_code_engine_build" "ce_build" {
  project_id    = ibm_code_engine_project.ce_project.id
  name          = "ce-build"
  output_image  = "private.de.icr.io/icr_namespace/image-name"
  output_secret = "ce-auto-icr-private-eu-de"
  source_url    = "https://github.com/IBM/CodeEngine"
  strategy_type = "dockerfile"
  strategy_size = "small"
}

resource "ibm_code_engine_build" "ce_build2" {
  project_id    = ibm_code_engine_project.ce_project.id
  name          = "ce-build2"
  output_image  = "private.de.icr.io/icr_namespace/image-name"
  output_secret = "ce-auto-icr-private-eu-de"
  source_url    = "https://github.com/IBM/CodeEngine"
  strategy_type = "dockerfile"
}

resource "ibm_code_engine_build" "ce_build3" {
  project_id    = ibm_code_engine_project.ce_project.id
  name          = "ce-build3"
  output_image  = "private.de.icr.io/icr_namespace/image-name"
  output_secret = "ce-auto-icr-private-eu-de"
  source_url    = "https://github.com/IBM/CodeEngine"
  strategy_type = "dockerfile"
  strategy_size = "large"
}

resource "ibm_code_engine_build" "ce_build4" {
  project_id    = ibm_code_engine_project.ce_project.id
  name          = "ce-build4"
  output_image  = "private.de.icr.io/icr_namespace/image-name"
  output_secret = "ce-auto-icr-private-eu-de"
  source_url    = "https://github.com/IBM/CodeEngine"
  strategy_type = "dockerfile"
  strategy_size = "xlarge"
}

resource "ibm_code_engine_build" "ce_build5" {
  project_id    = ibm_code_engine_project.ce_project.id
  name          = "ce-build5"
  output_image  = "private.de.icr.io/icr_namespace/image-name"
  output_secret = "ce-auto-icr-private-eu-de"
  source_url    = "https://github.com/IBM/CodeEngine"
  strategy_type = "dockerfile"
  strategy_size = "xxlarge"
}