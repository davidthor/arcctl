variable "name" {
  description = "Secret name"
  type        = string
}

variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "data" {
  description = "Secret data to store"
  type        = string
  sensitive   = true
}
