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

resource "ibm_cloudant" "cloudant" {
  name     = "test-cloudant"
  location = "us-south"
  plan     = "standard"

  legacy_credentials  = true
  include_data_events = false
  capacity            = 1
  enable_cors         = true

  cors_config {
    allow_credentials = false
    origins           = ["https://example.com"]
  }

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}