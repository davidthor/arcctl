variable "name" {
  description = "Task name"
  type        = string
}

variable "region" {
  description = "AWS region"
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

variable "log_group" {
  description = "CloudWatch log group name"
  type        = string
}

variable "image" {
  description = "Container image URI"
  type        = string
}

variable "command" {
  description = "Command override"
  type        = list(string)
  default     = null
}

variable "environment" {
  description = "Environment variables"
  type        = map(string)
  default     = {}
}

variable "timeout" {
  description = "Lambda timeout in seconds (max 900)"
  type        = number
  default     = 900
}

variable "memory" {
  description = "Lambda memory in MB"
  type        = number
  default     = 512
}
