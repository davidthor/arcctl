variable "name" {
  description = "Scheduled task name"
  type        = string
}

variable "cluster" {
  description = "ECS cluster name"
  type        = string
}

variable "image" {
  description = "Container image URI"
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

variable "schedule" {
  description = "EventBridge schedule expression (rate or cron)"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
  default     = null
}

variable "security_group_id" {
  description = "Security group ID"
  type        = string
  default     = null
}

variable "log_group" {
  description = "CloudWatch log group name"
  type        = string
  default     = null
}

variable "cpu" {
  description = "CPU units"
  type        = string
  default     = "256"
}

variable "memory" {
  description = "Memory in MiB"
  type        = string
  default     = "512"
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
  description = "Runtime (unused for cron)"
  type        = any
  default     = null
}

variable "timeout" {
  description = "Task timeout"
  type        = any
  default     = null
}
