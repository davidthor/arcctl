variable "name" {
  description = "App Platform service name"
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
  description = "Container command override"
  type        = list(string)
  default     = null
}

variable "environment" {
  description = "Environment variables"
  type        = map(string)
  default     = null
}

variable "replicas" {
  description = "Number of instances"
  type        = number
  default     = null
}

variable "port" {
  description = "HTTP port"
  type        = number
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

variable "health_check_path" {
  description = "Health check HTTP path"
  type        = string
  default     = "/healthz"
}

variable "runtime" {
  description = "Runtime configuration (unused in App Platform)"
  type        = any
  default     = null
}

variable "volumes" {
  description = "Volume configuration (unused)"
  type        = any
  default     = null
}
