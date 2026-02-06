variable "name" {
  description = "Name for the cron VM"
  type        = string
}

variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "zone" {
  description = "GCP zone"
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

variable "schedule" {
  description = "Cron schedule expression"
  type        = string
}

variable "image" {
  description = "Docker image to run on schedule"
  type        = string
}

variable "command" {
  description = "Command to execute"
  type        = list(string)
  default     = null
}

variable "environment" {
  description = "Environment variables"
  type        = map(string)
  default     = {}
}

variable "ssh_key" {
  description = "SSH public key for access"
  type        = string
  default     = ""
}

variable "tags" {
  description = "Network tags"
  type        = list(string)
  default     = null
}
