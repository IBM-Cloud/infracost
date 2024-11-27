terraform {
  required_providers {
    ibm = {
      source = "IBM-Cloud/ibm"
    }
  }
}

locals {
  event_notifications = {
    plans = ["lite", "standard"],
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

resource "ibm_en_subscription_cos" "subscription_cos_lite" {
  instance_guid  = ibm_resource_instance.event_notifications["lite"].guid
  name           = "Cloud Object Storage Subscription"
  description    = "Cloud Object Storage Subscription for Event Notifications"
  destination_id = "somedestinationid"
  topic_id       = "sometopicid"
}

resource "ibm_en_subscription_cos" "subscription_cos_standard" {
  instance_guid  = ibm_resource_instance.event_notifications["standard"].guid
  name           = "Cloud Object Storage Subscription"
  description    = "Cloud Object Storage Subscription for Event Notifications"
  destination_id = "somedestinationid"
  topic_id       = "sometopicid"
}

resource "ibm_en_subscription_cos" "subscription_cos_standard_no_usage" {
  instance_guid  = ibm_resource_instance.event_notifications["standard"].guid
  name           = "Cloud Object Storage Subscription"
  description    = "Cloud Object Storage Subscription for Event Notifications"
  destination_id = "somedestinationid"
  topic_id       = "sometopicid"
}
