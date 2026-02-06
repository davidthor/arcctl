variable "name" {
  description = "Service name for DNS"
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

variable "target" {
  description = "IP address of the target VM"
  type        = string
}

variable "target_type" {
  description = "Type of target (deployment or function)"
  type        = string
  default     = "deployment"
}

variable "port" {
  description = "Service port"
  type        = number
  default     = null
}
