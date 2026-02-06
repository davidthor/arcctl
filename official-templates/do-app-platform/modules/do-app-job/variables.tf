variable "name" {
  description = "Job name"
  type        = string
}

variable "image" {
  description = "Container image"
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
  default     = null
}

variable "region" {
  description = "DigitalOcean App Platform region"
  type        = string
}

variable "token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
}

variable "schedule" {
  description = "Cron schedule for recurring jobs"
  type        = string
  default     = null
}

variable "cpu" {
  description = "CPU allocation (unused, controlled by instance_size)"
  type        = string
  default     = null
}

variable "memory" {
  description = "Memory allocation (unused, controlled by instance_size)"
  type        = string
  default     = null
}

variable "runtime" {
  description = "Runtime configuration (unused)"
  type        = any
  default     = null
}

variable "replicas" {
  description = "Replicas (unused)"
  type        = number
  default     = null
}

variable "port" {
  description = "Port (unused)"
  type        = number
  default     = null
}

variable "volumes" {
  description = "Volume configuration (unused)"
  type        = any
  default     = null
}
