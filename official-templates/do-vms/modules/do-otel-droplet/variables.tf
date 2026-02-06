variable "name" {
  description = "OTel collector Droplet name"
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
