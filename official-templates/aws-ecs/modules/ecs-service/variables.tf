variable "cluster" {
  description = "ECS cluster name"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "security_group_id" {
  description = "Security group ID for the service"
  type        = string
}

variable "log_group" {
  description = "CloudWatch log group name"
  type        = string
}

variable "image" {
  description = "Container image URI"
  type        = string
}

variable "name" {
  description = "Service name"
  type        = string
  default     = "service"
}

variable "command" {
  description = "Container command override"
  type        = list(string)
  default     = null
}

variable "environment" {
  description = "Environment variables map"
  type        = map(string)
  default     = {}
}

variable "cpu" {
  description = "CPU units (256, 512, 1024, 2048, 4096)"
  type        = string
  default     = "256"
}

variable "memory" {
  description = "Memory in MiB (512, 1024, 2048, ...)"
  type        = string
  default     = "512"
}

variable "replicas" {
  description = "Number of desired tasks"
  type        = number
  default     = 1
}

variable "port" {
  description = "Container port"
  type        = number
  default     = 8080
}

# Accept and ignore additional inputs from merge(node.inputs, ...)
variable "runtime" {
  description = "Runtime specification (unused in ECS container mode)"
  type        = any
  default     = null
}
