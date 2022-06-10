terraform {
  required_providers {
    ibm = {
      source  = "IBM-Cloud/ibm"
      version = "1.42.0-beta0"
    }
  }
}

provider "ibm" {
  # Configuration options
  region = "us-south"
}

resource "ibm_cloudant" "cloudant_lite" {
  name     = "cloudant_prt"
  location = "us-south"
  plan     = "lite"
  capacity = 1

  legacy_credentials  = true
  include_data_events = false
  enable_cors         = true

  cors_config {
    allow_credentials = false
    origins           = ["https://example.com"]
  }
}

resource "ibm_cloudant" "cloudant_std" {
  name     = "cloudant_prt"
  location = "us-south"
  plan     = "standard"
  capacity = 1

  legacy_credentials  = true
  include_data_events = false
  enable_cors         = true

  cors_config {
    allow_credentials = false
    origins           = ["https://example.com"]
  }
}

resource "ibm_cloudant" "with_usage" {
  name     = "cloudant_prt"
  location = "us-south"
  plan     = "standard"
  capacity = 2

  legacy_credentials  = true
  include_data_events = false
  enable_cors         = true

  cors_config {
    allow_credentials = false
    origins           = ["https://example.com"]
  }
}



