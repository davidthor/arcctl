variable "name" {
  description = "Service/route name"
  type        = string
}

variable "domain" {
  description = "Custom domain name"
  type        = string
  default     = ""
}

variable "token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
}

variable "service" {
  description = "Target service name"
  type        = string
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

variable "port" {
  description = "Target port"
  type        = number
  default     = null
}

variable "path" {
  description = "URL path prefix"
  type        = string
  default     = null
}

variable "hostnames" {
  description = "Hostnames for the route"
  type        = list(string)
  default     = []
}
