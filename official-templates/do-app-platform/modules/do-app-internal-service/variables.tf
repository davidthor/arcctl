variable "name" {
  description = "Internal service name"
  type        = string
}

variable "token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
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

variable "port" {
  description = "Service port"
  type        = number
  default     = null
}
