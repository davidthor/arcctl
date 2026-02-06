variable "name" {
  description = "Task name"
  type        = string
}

variable "cluster" {
  description = "ECS cluster name"
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
  description = "Command to execute"
  type        = list(string)
  default     = null
}

variable "environment" {
  description = "Environment variables map"
  type        = map(string)
  default     = {}
}

variable "capacity_provider" {
  description = "Capacity provider (FARGATE or FARGATE_SPOT)"
  type        = string
  default     = "FARGATE_SPOT"
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
