variable "name" {
  description = "Instance name"
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

variable "key_pair" {
  description = "SSH key pair name"
  type        = string
  default     = ""
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.small"
}

variable "target_group_arn" {
  description = "ALB target group ARN"
  type        = string
  default     = ""
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
  description = "Container command"
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
  type        = any
  default     = null
}

variable "memory" {
  description = "Memory specification"
  type        = any
  default     = null
}

variable "replicas" {
  description = "Replicas"
  type        = any
  default     = null
}

variable "port" {
  description = "Container port"
  type        = number
  default     = 8080
}

variable "runtime" {
  description = "Runtime (unused in Docker mode)"
  type        = any
  default     = null
}
