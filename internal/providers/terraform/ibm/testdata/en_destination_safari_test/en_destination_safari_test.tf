
terraform {
  required_providers {
    ibm = {
      source = "IBM-Cloud/ibm"
    }
  }
}

locals {
  event_notifications = {
    plans                = ["lite", "standard"],
    destination_pre_prod = ["true", "false"]
  }
}

resource "ibm_resource_instance" "event_notifications" {
  for_each          = toset(local.event_notifications.plans)
  name              = "event-notifications-${each.value}"
  location          = "us-south"
  plan              = each.value
  resource_group_id = "default"
  service           = "event-notifications"
}

resource "ibm_en_destination_safari" "destination_safari_lite" {
  for_each              = toset(local.event_notifications.destination_pre_prod)
  instance_guid         = ibm_resource_instance.event_notifications["lite"].guid
  name                  = "Safari Destination"
  certificate           = "${path.module}/Certificates/safaricert.p12"
  collect_failed_events = false
  description           = "Safari push destination"
  type                  = "push_safari"
  config {
    params {
      cert_type         = "p12"
      password          = "apnscertpassword"
      url_format_string = "https://test.petstorez.com"
      website_name      = "petstore"
      website_push_id   = "petzz"
      website_url       = "https://test.petstorez.com"
      pre_prod          = tobool(each.value)
    }
  }
}

resource "ibm_en_destination_safari" "destination_safari_standard" {
  for_each              = toset(local.event_notifications.destination_pre_prod)
  instance_guid         = ibm_resource_instance.event_notifications["standard"].guid
  name                  = "Safari Destination"
  certificate           = "${path.module}/Certificates/safaricert.p12"
  collect_failed_events = false
  description           = "Safari push destination"
  type                  = "push_safari"
  config {
    params {
      cert_type         = "p12"
      password          = "apnscertpassword"
      url_format_string = "https://test.petstorez.com"
      website_name      = "petstore"
      website_push_id   = "petzz"
      website_url       = "https://test.petstorez.com"
      pre_prod          = tobool(each.value)
    }
  }
}
