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

resource "ibm_database" "rabbitmq_b3c" {
  name              = "rabbitmq_b3c"
  plan              = "standard"
  location          = "us-south"
  service           = "messages-for-rabbitmq"
  service_endpoints = "private"
  group {
    group_id = "member"
    host_flavor {
      id = "b3c.4x16.encrypted"
    }
    disk {
      allocation_mb = 4194304
    }
  } 
}

# -------------------------------------------
# POSTGRES
# -------------------------------------------

resource "ibm_database" "postgresql_standard_flavor" {
  name     = "postgres-standard-flavour"
  service  = "databases-for-postgresql"
  plan     = "standard"
  location = "us-south"
  service_endpoints = "private"
  group { # Note: "memory" not allowed when host_flavor is set
    group_id = "member"
    host_flavor {
      id = "m3c.30x240.encrypted"
    }
    disk { # >= 5120 and <= 4194304 in increments of 1024
      allocation_mb = 4194304
    }
  }
  configuration = <<CONFIGURATION
  {
    "max_connections": 400
  }
  CONFIGURATION
}

resource "ibm_database" "postgresql_standard_multitenant_flavor" {
  name     = "postgres-standard-multitenant-flavour"
  service  = "databases-for-postgresql"
  plan     = "standard"
  location = "us-south"
  service_endpoints = "private"
  group { # Note: "memory" not allowed when host_flavor is set
    group_id = "member"
    host_flavor {
      id = "multitenant"
    }
    disk { # >= 5120 and <= 4194304 in increments of 1024
      allocation_mb = 4194304
    }
    memory { # >= 4096 and <= 114688 in increments of 128
      allocation_mb = 114688
    }
    cpu { # >= 0 and <= 28 in increments of 1
      allocation_count = 28
    }
  }
  configuration = <<CONFIGURATION
  {
    "max_connections": 400
  }
  CONFIGURATION
}

resource "ibm_database" "postgresql_standard" {
  name     = "postgres-standard"
  service  = "databases-for-postgresql"
  plan     = "standard"
  location = "us-south"
  service_endpoints = "private"
  group {
    group_id = "member"
    memory { # >= 1024 and <= 114688 in increments of 128
      allocation_mb = 114688
    }
    disk { # >= 5120 and <= 4194304 in increments of 1024
      allocation_mb = 4194304
    }
    cpu { # >= 0 and <= 28 in increments of 1
      allocation_count = 28
    }
  }
  configuration = <<CONFIGURATION
  {
    "max_connections": 400
  }
  CONFIGURATION
}

# -------------------------------------------
# ELASTICSEARCH
# -------------------------------------------

resource "ibm_database" "elasticsearch_platinum" {
  name     = "elasticsearch-platinum"
  service  = "databases-for-elasticsearch"
  plan     = "platinum"
  location = "us-south"
  service_endpoints = "private"
  group {
    group_id = "member"
    memory { # >= 1024 and <= 114688 in increments of 128
      allocation_mb = 114688
    }
    disk { # >= 5120 and <= 4194304 in increments of 1024
      allocation_mb = 4194304
    }
    cpu { # >= 0 and <= 28 in increments of 1
      allocation_count = 28
    }
  }
}

resource "ibm_database" "elasticsearch_platinum_flavor" {
  name     = "elasticsearch-platinum-flavor"
  service  = "databases-for-elasticsearch"
  plan     = "platinum"
  location = "us-south"
  service_endpoints = "private"
  group { # Note: "memory" not allowed when host_flavor is set
    group_id = "member"
    host_flavor {
      id = "m3c.30x240.encrypted"
    }
    disk { # >= 5120 and <= 4194304 in increments of 1024
      allocation_mb = 4194304
    }
  }
}

resource "ibm_database" "elasticsearch_enterprise" {
  name     = "elasticsearch-enterprise"
  service  = "databases-for-elasticsearch"
  plan     = "enterprise"
  location = "us-south"
  service_endpoints = "private"
  group {
    group_id = "member"

    memory { # >= 1024 and <= 114688 in increments of 128
      allocation_mb = 114688
    }
    disk { # >= 5120 and <= 4194304 in increments of 1024
      allocation_mb = 4194304
    }
    cpu { # >= 0 and <= 28 in increments of 1
      allocation_count = 28
    }
  }
}

resource "ibm_database" "elasticsearch_enterprise_flavor" {
  name     = "elasticsearch-enterprise-flavor"
  service  = "databases-for-elasticsearch"
  plan     = "enterprise"
  location = "us-south"
  service_endpoints = "private"
  group { # Note: "memory" not allowed when host_flavor is set
    group_id = "member"
    host_flavor {
      id = "m3c.30x240.encrypted"
    }
    disk { # >= 5120 and <= 4194304 in increments of 1024
      allocation_mb = 4194304
    }
  }
}

resource "ibm_database" "elasticsearch_enterprise_multitenant_flavor" {
  name     = "elasticsearch-enterprise-multitenant-flavor"
  service  = "databases-for-elasticsearch"
  plan     = "enterprise"
  location = "us-south"
  service_endpoints = "private"
  group {
    group_id = "member"
    host_flavor {
      id = "multitenant"
    }
    disk { # >= 5120 and <= 4194304 in increments of 1024
      allocation_mb = 4194304
    }
    memory { # >= 4096 and <= 114688 in increments of 128
      allocation_mb = 114688
    }
    cpu { # >= 0 and <= 28 in increments of 1
      # allocation_count = 0 # Automatically allocate based on a 1:8 ration with RAM
      allocation_count = 28
    }
  }
}

# Specifications used by Dev RAG stack
resource "ibm_database" "elasticsearch_enterprise_multitenant_flavor_auto_cpu_scale" {
  name     = "elasticsearch-enterprise-multitenant-flavor-auto-cpu-scale"
  service  = "databases-for-elasticsearch"
  plan     = "enterprise"
  location = "us-south"
  service_endpoints = "private"
  group {
    group_id = "member"
    host_flavor {
      id = "multitenant"
    }
    disk { # >= 5120 and <= 4194304 in increments of 1024
      allocation_mb = 5120
    }
    memory { # >= 4096 and <= 114688 in increments of 128
      allocation_mb = 4096
    }
    cpu {                  # >= 0 and <= 28 in increments of 1
      allocation_count = 0 # Automatically allocate based on a 1:8 ration with RAM
    }
  }
}

resource "ibm_database" "rabbitmq_multitenant" {
  name              = "rabbitmq_multitenant"
  plan              = "standard"
  location          = "us-south"
  service           = "messages-for-rabbitmq"
  service_endpoints = "private"
  tags              = ["tag1", "tag2"]
  group {
    group_id = "member"
    host_flavor {
      id = "multitenant"
    }
    cpu {
      allocation_count = 3
    }
    memory {
      allocation_mb = 12288
    }
    disk {
      allocation_mb = 256000
    }
  }
}

