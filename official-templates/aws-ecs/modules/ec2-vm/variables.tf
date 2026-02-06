variable "name" {
  description = "Instance name"
  type        = string
}

variable "subnet_id" {
  description = "Subnet ID for the instance"
  type        = string
}

variable "security_group" {
  description = "Security group ID"
  type        = string
}

variable "runtime" {
  description = "Runtime specification (string shorthand or object)"
  type        = any
  default     = "node:20"
}

variable "command" {
  description = "Command to run"
  type        = list(string)
  default     = []
}

variable "environment" {
  description = "Environment variables"
  type        = map(string)
  default     = {}
}

variable "cpu" {
  description = "CPU specification"
  type        = string
  default     = null
}

variable "memory" {
  description = "Memory specification"
  type        = string
  default     = null
}

variable "replicas" {
  description = "Number of replicas (only 1 supported for VM)"
  type        = number
  default     = 1
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.small"
}

variable "image" {
  description = "Container image (unused in VM mode)"
  type        = string
  default     = null
}

variable "port" {
  description = "Application port"
  type        = number
  default     = 8080
}
