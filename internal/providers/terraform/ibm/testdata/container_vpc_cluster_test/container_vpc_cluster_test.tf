terraform {
  required_providers {
    ibm = {
      source = "IBM-Cloud/ibm"
    }
  }
}
provider "ibm" {
  region           = "us-south"
  ibmcloud_timeout = "1"
  max_retries      = "1"
}

resource "ibm_is_vpc" "vpc1" {
  name = "myvpc"
}

resource "ibm_is_subnet" "subnet1" {
  name                     = "mysubnet1"
  vpc                      = ibm_is_vpc.vpc1.id
  zone                     = "us-south-1"
  total_ipv4_address_count = 256
}

resource "ibm_is_subnet" "subnet2" {
  name                     = "mysubnet2"
  vpc                      = ibm_is_vpc.vpc1.id
  zone                     = "us-south-2"
  total_ipv4_address_count = 256
}

resource "ibm_container_vpc_cluster" "cluster" {
  name         = "mycluster"
  vpc_id       = ibm_is_vpc.vpc1.id
  flavor       = "bx2.4x16"
  worker_count = 3
  kube_version = "1.17.5"
  zones {
    subnet_id = ibm_is_subnet.subnet1.id
    name      = "us-south-1"
  }
  zones {
    subnet_id = ibm_is_subnet.subnet1.id
    name      = "us-south-2"
  }
}

resource "ibm_container_vpc_cluster" "cluster_without_usage" {
  name         = "mycluster-without-usage"
  vpc_id       = ibm_is_vpc.vpc1.id
  flavor       = "bx2.4x16"
  worker_count = 3
  kube_version = "1.17.5"
  zones {
    subnet_id = ibm_is_subnet.subnet1.id
    name      = "us-south-1"
  }
  zones {
    subnet_id = ibm_is_subnet.subnet1.id
    name      = "us-south-2"
  }
}

resource "ibm_container_vpc_cluster" "roks_cluster_with_usage" {
  name         = "myrokscluster-with-usage"
  vpc_id       = ibm_is_vpc.vpc1.id
  flavor       = "bx2.4x16"
  worker_count = 3
  kube_version = "4.13_openshift"
  zones {
    subnet_id = ibm_is_subnet.subnet1.id
    name      = "us-south-1"
  }
  zones {
    subnet_id = ibm_is_subnet.subnet1.id
    name      = "us-south-2"
  }
}

resource "ibm_container_vpc_cluster" "roks_with_entitlement" {
  name         = "roks-with-entitlement"
  vpc_id       = ibm_is_vpc.vpc1.id
  flavor       = "bx2.4x16"
  worker_count = 3
  kube_version = "4.13_openshift"
  entitlement  = "cloud_pak"
  zones {
    subnet_id = ibm_is_subnet.subnet1.id
    name      = "us-south-1"
  }
  zones {
    subnet_id = ibm_is_subnet.subnet1.id
    name      = "us-south-2"
  }
}

resource "ibm_container_vpc_cluster" "roks_no_entitlement" {
  name         = "roks-no-entitlement"
  vpc_id       = ibm_is_vpc.vpc1.id
  flavor       = "bx2.4x16"
  worker_count = 3
  kube_version = "4.13_openshift"
  zones {
    subnet_id = ibm_is_subnet.subnet1.id
    name      = "us-south-1"
  }
  zones {
    subnet_id = ibm_is_subnet.subnet1.id
    name      = "us-south-2"
  }
}

/*
  Copies the configuration from this DA: 
  https://cloud.ibm.com/catalog/7a4d68b4-cf8b-40cd-a3d1-f49aff526eb3/architecture/deploy-arch-ibm-ocp-vpc-1728a4fd-f561-4cf9-82ef-2b1eeb5da1a8-global

  Uses 2 clusters and 2 worker pools each consisting of 730 workers. So 1460 workers
  used for each resource to represent 2 instances of a cluster/pool.
*/ 
resource "ibm_container_vpc_cluster" "cluster_with_pool" {
  name         = "cluster_with_pool"
  vpc_id       = ibm_is_vpc.vpc1.id
  flavor       = "bx2.8x32"
  worker_count = 1460
  kube_version = "4.17_openshift" // Version as shown in console
  zones {
    subnet_id = ibm_is_subnet.subnet1.id
    name      = "us-south-1"
  }
}

resource "ibm_container_vpc_worker_pool" "cluster_pool" {
  cluster = ibm_container_vpc_cluster.cluster.id
  worker_pool_name = "mywp"
  flavor = "bx2.8x32"
  vpc_id = ibm_is_vpc.vpc1.id
  worker_count = 1460
  zones {
    name = "us-south-2"
    subnet_id = ibm_is_subnet.subnet2.id
  }
}