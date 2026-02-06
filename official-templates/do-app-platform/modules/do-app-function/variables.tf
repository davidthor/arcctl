variable "name" {
  description = "Function name"
  type        = string
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

variable "image" {
  description = "Container image"
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
  default     = null
}

variable "cpu" {
  description = "CPU allocation"
  type        = string
  default     = null
}

variable "memory" {
  description = "Memory allocation"
  type        = string
  default     = null
}

variable "port" {
  description = "HTTP port"
  type        = number
  default     = null
}

variable "timeout" {
  description = "Function timeout in seconds"
  type        = number
  default     = 300
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

variable "volumes" {
  description = "Volume configuration (unused)"
  type        = any
  default     = null
}
