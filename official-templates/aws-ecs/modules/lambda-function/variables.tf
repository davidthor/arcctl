variable "name" {
  description = "Lambda function name"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "security_group_id" {
  description = "Security group ID"
  type        = string
}

variable "image" {
  description = "Container image URI"
  type        = string
}

variable "environment" {
  description = "Environment variables"
  type        = map(string)
  default     = {}
}

variable "timeout" {
  description = "Function timeout in seconds"
  type        = number
  default     = 30
}

variable "memory" {
  description = "Function memory in MB"
  type        = any
  default     = 128
}

variable "command" {
  description = "Command override"
  type        = list(string)
  default     = null
}

variable "cpu" {
  description = "CPU (unused for Lambda)"
  type        = any
  default     = null
}

variable "port" {
  description = "Port (unused for Lambda)"
  type        = any
  default     = null
}

variable "replicas" {
  description = "Replicas (unused for Lambda)"
  type        = any
  default     = null
}

variable "runtime" {
  description = "Runtime (unused for Lambda container images)"
  type        = any
  default     = null
}
