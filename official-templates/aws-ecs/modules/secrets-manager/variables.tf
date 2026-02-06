variable "name" {
  description = "Secret name/path in Secrets Manager"
  type        = string
}

variable "data" {
  description = "Secret data (string or JSON)"
  type        = string
  sensitive   = true
}
