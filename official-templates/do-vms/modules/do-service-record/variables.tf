variable "name" {
  description = "Service record name"
  type        = string
}

variable "domain" {
  description = "Base domain name"
  type        = string
}

variable "do_token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
}

variable "target" {
  description = "Target IP address (Droplet private IP)"
  type        = string
  default     = "127.0.0.1"
}

variable "port" {
  description = "Service port"
  type        = number
  default     = null
}

variable "deployment" {
  description = "Target deployment name"
  type        = string
  default     = null
}

variable "function" {
  description = "Target function name"
  type        = string
  default     = null
}
