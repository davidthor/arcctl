variable "name" {
  description = "CronJob name"
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
  default     = "t3.micro"
}

variable "image" {
  description = "Container image URI"
  type        = string
}

variable "command" {
  description = "Command to execute"
  type        = list(string)
  default     = []
}

variable "environment" {
  description = "Environment variables"
  type        = map(string)
  default     = {}
}

variable "schedule" {
  description = "Cron schedule expression"
  type        = string
}

variable "cpu" {
  description = "CPU (unused)"
  type        = any
  default     = null
}

variable "memory" {
  description = "Memory (unused)"
  type        = any
  default     = null
}

variable "port" {
  description = "Port (unused)"
  type        = any
  default     = null
}

variable "replicas" {
  description = "Replicas (unused)"
  type        = any
  default     = null
}

variable "runtime" {
  description = "Runtime (unused)"
  type        = any
  default     = null
}

variable "timeout" {
  description = "Timeout (unused)"
  type        = any
  default     = null
}
