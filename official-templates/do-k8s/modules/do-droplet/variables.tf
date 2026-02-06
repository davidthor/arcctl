variable "name" {
  description = "Droplet name"
  type        = string
}

variable "region" {
  description = "DigitalOcean region"
  type        = string
}

variable "do_token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
}

variable "ssh_key_fingerprint" {
  description = "SSH key fingerprint for Droplet access"
  type        = string
  default     = ""
}

variable "runtime" {
  description = "Runtime configuration (language:version or object)"
  type        = any
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

variable "replicas" {
  description = "Number of replicas"
  type        = number
  default     = 1
}

variable "size" {
  description = "Droplet size override"
  type        = string
  default     = null
}

variable "image" {
  description = "Container image (unused for runtime-based deployments)"
  type        = string
  default     = null
}

variable "port" {
  description = "Application port"
  type        = number
  default     = null
}

variable "volumes" {
  description = "Volume configuration"
  type        = any
  default     = null
}
