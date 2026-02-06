variable "name" {
  description = "Name for the Compute Engine instance"
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

variable "machine_type" {
  description = "Compute Engine machine type"
  type        = string
  default     = null
}

variable "network" {
  description = "VPC network ID"
  type        = string
}

variable "subnet" {
  description = "Subnet ID"
  type        = string
}

variable "runtime" {
  description = "Runtime specification (string like 'node:20' or object)"
  type        = any
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

variable "cpu" {
  description = "CPU specification (informational, machine_type determines actual CPU)"
  type        = string
  default     = null
}

variable "memory" {
  description = "Memory specification (informational, machine_type determines actual memory)"
  type        = string
  default     = null
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

variable "replicas" {
  description = "Number of replicas (for compatibility, VM creates a single instance)"
  type        = number
  default     = null
}
