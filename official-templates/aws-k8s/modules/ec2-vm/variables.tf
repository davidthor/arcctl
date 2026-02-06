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

variable "runtime" {
  description = "Runtime specification"
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

variable "image" {
  description = "Container image (unused)"
  type        = any
  default     = null
}

variable "port" {
  description = "Application port"
  type        = any
  default     = null
}
