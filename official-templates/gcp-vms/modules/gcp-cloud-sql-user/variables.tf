variable "database" {
  description = "Database connection information from the parent database resource"
  type = object({
    instance_name = string
    database      = string
    host          = string
    port          = number
    scheme        = string
  })
}

variable "username" {
  description = "Username for the new database user"
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
