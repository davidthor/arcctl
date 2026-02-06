variable "name" {
  description = "Name for the VPC connector"
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
