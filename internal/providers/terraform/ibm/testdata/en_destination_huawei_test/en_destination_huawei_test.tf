
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

resource "ibm_en_destination_huawei" "destination_huawei_lite" {
  for_each              = toset(local.event_notifications.destination_pre_prod)
  instance_guid         = ibm_resource_instance.event_notifications["lite"].guid
  name                  = "Huawei Destination"
  type                  = "push_huawei"
  collect_failed_events = false
  description           = "Huawei push destination"
  config {
    params {
      client_id     = "clientid"
      client_secret = "clientsecret"
      pre_prod      = tobool(each.value)
    }
  }
}

resource "ibm_en_destination_huawei" "destination_huawei_standard" {
  for_each              = toset(local.event_notifications.destination_pre_prod)
  instance_guid         = ibm_resource_instance.event_notifications["standard"].guid
  name                  = "Huawei Destination"
  type                  = "push_huawei"
  collect_failed_events = false
  description           = "Huawei push destination"
  config {
    params {
      client_id     = "clientid"
      client_secret = "clientsecret" // pragma: allowlist secret
      pre_prod      = tobool(each.value)
    }
  }
}

