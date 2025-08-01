terraform {
  required_providers {
    ibm = {
      source = "IBM-Cloud/ibm"
    }
  }
}

provider "ibm" {
  generation = 2
  region     = "eu-de"
}

resource "ibm_is_vpc" "vpc1" {
  name = "myvpc"
}

resource "ibm_is_subnet" "subnet1" {
  name                     = "mysubnet1"
  vpc                      = ibm_is_vpc.vpc1.id
  zone                     = "eu-de-1"
  total_ipv4_address_count = 256
}

resource "ibm_is_subnet" "subnet2" {
  name                     = "mysubnet2"
  vpc                      = ibm_is_vpc.vpc1.id
  zone                     = "eu-de-2"
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
    name      = "eu-de-1"
  }
}

resource "ibm_container_vpc_worker_pool" "cluster_pool" {
  cluster          = ibm_container_vpc_cluster.cluster.id
  worker_pool_name = "mywp"
  flavor           = "bx2.2x8"
  vpc_id           = ibm_is_vpc.vpc1.id
  worker_count     = 3
  zones {
    name      = "eu-de-2"
    subnet_id = ibm_is_subnet.subnet2.id
  }
}

resource "ibm_container_vpc_worker_pool" "cluster_pool_without_usage" {
  cluster          = ibm_container_vpc_cluster.cluster.id
  worker_pool_name = "mywp"
  flavor           = "bx2.2x8"
  vpc_id           = ibm_is_vpc.vpc1.id
  worker_count     = 3
  zones {
    name      = "eu-de-2"
    subnet_id = ibm_is_subnet.subnet2.id
  }
}

resource "ibm_container_vpc_worker_pool" "default_cluster_pool" {
  cluster          = ibm_container_vpc_cluster.cluster.id
  worker_pool_name = "default"
  flavor           = "bx2.2x8"
  vpc_id           = ibm_is_vpc.vpc1.id
  worker_count     = 3
  zones {
    name      = "eu-de-2"
    subnet_id = ibm_is_subnet.subnet2.id
  }
}

resource "ibm_container_vpc_worker_pool" "roks_worker_pool" {
  cluster          = ibm_container_vpc_cluster.cluster.id
  worker_pool_name = "workerpool-roks"
  flavor           = "bx2.4x16"
  operating_system = "RHEL_9_64"
  vpc_id           = ibm_is_vpc.vpc1.id
  worker_count     = 3
  zones {
    name      = "eu-de-2"
    subnet_id = ibm_is_subnet.subnet2.id
  }
}

resource "ibm_container_vpc_worker_pool" "ubuntu_worker_pool" {
  cluster          = ibm_container_vpc_cluster.cluster.id
  worker_pool_name = "workerpool-ubuntu"
  flavor           = "bx2.4x16"
  operating_system = "UBUNTU_24_64"
  vpc_id           = ibm_is_vpc.vpc1.id
  worker_count     = 3
  zones {
    name      = "eu-de-2"
    subnet_id = ibm_is_subnet.subnet2.id
  }
}
