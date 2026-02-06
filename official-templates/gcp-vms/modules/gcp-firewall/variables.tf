variable "name" {
  description = "Name prefix for the firewall rules"
  type        = string
}

variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "network" {
  description = "VPC network ID"
  type        = string
}

variable "tags" {
  description = "Network tags to apply the firewall rules to"
  type        = list(string)
}
