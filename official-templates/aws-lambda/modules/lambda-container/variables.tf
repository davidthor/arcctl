variable "name" {
  description = "Lambda function name"
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

variable "api_id" {
  description = "API Gateway ID"
  type        = string
  default     = null
}

variable "stage" {
  description = "API Gateway stage name"
  type        = string
  default     = null
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
