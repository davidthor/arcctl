variable "name" {
  description = "Rule name"
  type        = string
}

variable "load_balancer_id" {
  description = "Load balancer ID"
  type        = string
}

variable "domain" {
  description = "Base domain name"
  type        = string
}

variable "subdomain" {
  description = "Subdomain for the route"
  type        = string
}

variable "do_token" {
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
