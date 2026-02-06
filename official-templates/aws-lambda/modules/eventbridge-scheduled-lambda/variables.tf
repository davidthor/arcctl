variable "name" {
  description = "Scheduled Lambda name"
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

variable "schedule" {
  description = "EventBridge schedule expression"
  type        = string
}

variable "timeout" {
  description = "Lambda timeout in seconds"
  type        = number
  default     = 900
}

variable "memory" {
  description = "Lambda memory in MB"
  type        = any
  default     = 512
}

variable "cpu" {
  description = "CPU (unused for Lambda)"
  type        = any
  default     = null
}

variable "port" {
  description = "Port (unused for cron)"
  type        = any
  default     = null
}

variable "replicas" {
  description = "Replicas (unused for cron)"
  type        = any
  default     = null
}

variable "runtime" {
  description = "Runtime (unused)"
  type        = any
  default     = null
}
