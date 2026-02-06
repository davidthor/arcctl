variable "name" {
  description = "Database user/role name"
  type        = string
}

variable "api_key" {
  description = "Neon API key"
  type        = string
  sensitive   = true
}

variable "project_id" {
  description = "Neon project ID"
  type        = string
}

variable "branch" {
  description = "Neon branch ID"
  type        = string
}

variable "database" {
  description = "Database name for connection URL"
  type        = string
}
