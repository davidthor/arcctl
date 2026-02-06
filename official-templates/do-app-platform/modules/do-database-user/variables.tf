variable "name" {
  description = "Database user name"
  type        = string
}

variable "database" {
  description = "Parent database cluster name"
  type        = string
}

variable "token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
}
