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

resource "ibm_en_destination_android" "destination_android_lite" {
  for_each              = toset(local.event_notifications.destination_pre_prod)
  instance_guid         = ibm_resource_instance.event_notifications["lite"].guid
  name                  = "Android Destination"
  type                  = "push_android"
  collect_failed_events = false
  description           = "Android push destination"
  config {
    params {
      project_id   = "projectid"
      private_key  = "privatekey"
      client_email = "user@mail.com"
      pre_prod     = each.value
    }
  }
}

resource "ibm_en_destination_android" "destination_android_standard" {
  for_each              = toset(local.event_notifications.destination_pre_prod)
  instance_guid         = ibm_resource_instance.event_notifications["standard"].guid
  name                  = "Android Destination"
  type                  = "push_android"
  collect_failed_events = false
  description           = "Android push destination"
  config {
    params {
      project_id   = "projectid"
      private_key  = "privatekey"
      client_email = "user@mail.com"
      pre_prod     = tobool(each.value)
    }
  }
}
