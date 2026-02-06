variable "name" {
  description = "Droplet name"
  type        = string
}

variable "region" {
  description = "DigitalOcean region"
  type        = string
}

variable "size" {
  description = "Droplet size slug"
  type        = string
  default     = "s-1vcpu-1gb"
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

variable "vpc_uuid" {
  description = "VPC UUID for private networking"
  type        = string
  default     = null
}

variable "schedule" {
  description = "Cron schedule expression"
  type        = string
}

variable "image" {
  description = "Docker container image"
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
