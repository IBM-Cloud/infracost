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

resource "ibm_en_destination_ios" "destination_ios_lite" {
  for_each                 = toset(local.event_notifications.destination_pre_prod)
  instance_guid            = ibm_resource_instance.event_notifications["lite"].guid
  name                     = "iOS Destination Auth"
  type                     = "push_ios"
  collect_failed_events    = false
  certificate_content_type = "p8"
  certificate              = "${path.module}/Certificates/Auth.p8"
  description              = "iOS destination with P8"
  config {
    params {
      cert_type  = "p8"
      is_sandbox = true
      pre_prod    = tobool(each.value)
    }
  }
}

resource "ibm_en_destination_ios" "destination_ios_standard" {
  for_each                 = toset(local.event_notifications.destination_pre_prod)
  instance_guid            = ibm_resource_instance.event_notifications["standard"].guid
  name                     = "iOS Destination Auth"
  type                     = "push_ios"
  collect_failed_events    = false
  certificate_content_type = "p8"
  certificate              = "${path.module}/Certificates/Auth.p8"
  description              = "iOS destination with P8"
  config {
    params {
      cert_type  = "p8"
      is_sandbox = true
      pre_prod    = tobool(each.value)
    }
  }
}
