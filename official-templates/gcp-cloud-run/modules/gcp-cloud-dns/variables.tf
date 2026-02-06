variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "zone_name" {
  description = "Cloud DNS managed zone name"
  type        = string
}

variable "subdomain" {
  description = "Subdomain to create the record for"
  type        = string
}

variable "target" {
  description = "IP address target for the DNS record"
  type        = string
}
