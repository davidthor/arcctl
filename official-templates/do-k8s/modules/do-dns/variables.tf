variable "domain" {
  description = "Base domain name"
  type        = string
}

variable "subdomain" {
  description = "Subdomain for the DNS record"
  type        = string
}

variable "target" {
  description = "Target IP address for the A record"
  type        = string
}

variable "token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
}
