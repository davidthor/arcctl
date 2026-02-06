variable "name" {
  description = "Name for the Cloud SQL instance"
  type        = string
}

variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
}

variable "engine" {
  description = "Database engine: POSTGRES or MYSQL"
  type        = string
  validation {
    condition     = contains(["POSTGRES", "MYSQL"], var.engine)
    error_message = "Engine must be POSTGRES or MYSQL."
  }
}

variable "engine_version" {
  description = "Database engine version (e.g., 16 for Postgres, 8.0 for MySQL)"
  type        = string
}

variable "tier" {
  description = "Cloud SQL machine tier"
  type        = string
  default     = "db-f1-micro"
}

variable "network" {
  description = "VPC network ID for private IP connectivity"
  type        = string
}

variable "database_name" {
  description = "Name of the database to create"
  type        = string
}
