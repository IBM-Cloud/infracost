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

resource "ibm_en_destination_firefox" "destination_firefox_lite" {
  for_each              = toset(local.event_notifications.destination_pre_prod)
  instance_guid         = ibm_resource_instance.event_notifications["lite"].guid
  name                  = "Firefox Destination"
  type                  = "push_firefox"
  collect_failed_events = false
  description           = "Firefox push destination"
  config {
    params {
      website_url = "https://testevents.com"
      pre_prod    = tobool(each.value)
    }
  }
}

resource "ibm_en_destination_firefox" "destination_firefox_standard" {
  for_each              = toset(local.event_notifications.destination_pre_prod)
  instance_guid         = ibm_resource_instance.event_notifications["standard"].guid
  name                  = "Firefox Destination"
  type                  = "push_firefox"
  collect_failed_events = false
  description           = "Firefox push destination"
  config {
    params {
      website_url = "https://testevents.com"
      pre_prod    = tobool(each.value)
    }
  }
}

