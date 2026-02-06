variable "name" {
  description = "Spaces bucket name"
  type        = string
}

variable "region" {
  description = "DigitalOcean region"
  type        = string
}

variable "versioning" {
  description = "Enable versioning"
  type        = bool
  default     = false
}

variable "public" {
  description = "Enable public access"
  type        = bool
  default     = false
}

variable "token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
}
