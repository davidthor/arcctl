variable "name" {
  description = "Name for the GCS bucket (must be globally unique)"
  type        = string
}

variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region for the bucket"
  type        = string
}

variable "versioning" {
  description = "Enable object versioning"
  type        = bool
  default     = false
}

variable "public" {
  description = "Make the bucket publicly readable"
  type        = bool
  default     = false
}
