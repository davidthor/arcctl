variable "name" {
  description = "Database cluster name"
  type        = string
}

variable "type" {
  description = "Database type (postgres, mysql, redis, mongodb)"
  type        = string
}

variable "engine_version" {
  description = "Database version"
  type        = string
  default     = null
}

variable "region" {
  description = "DigitalOcean region"
  type        = string
}

variable "size" {
  description = "Database cluster size slug"
  type        = string
  default     = "db-s-1vcpu-1gb"
}

variable "token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
}
