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

variable "image" {
  description = "Docker image to run"
  type        = string
}

variable "command" {
  description = "Command to execute in the container"
  type        = list(string)
  default     = null
}

variable "environment" {
  description = "Environment variables"
  type        = map(string)
  default     = null
}
