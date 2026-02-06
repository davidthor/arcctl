variable "name" {
  description = "Name for the VPC network"
  type        = string
}

variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region for the subnet"
  type        = string
}
