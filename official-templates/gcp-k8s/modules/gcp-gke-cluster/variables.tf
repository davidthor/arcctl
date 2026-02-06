variable "name" {
  description = "Name of the GKE cluster"
  type        = string
}

variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
}

variable "network" {
  description = "VPC network ID"
  type        = string
}

variable "subnet" {
  description = "Subnet ID"
  type        = string
}

variable "node_pool" {
  description = "Node pool configuration"
  type = object({
    machine_type = string
    min_nodes    = number
    max_nodes    = number
    auto_scale   = bool
  })
  default = {
    machine_type = "e2-standard-4"
    min_nodes    = 1
    max_nodes    = 10
    auto_scale   = true
  }
}
