variable "name" {
  description = "Name for the Cloud Scheduler job"
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

variable "schedule" {
  description = "Cron schedule expression"
  type        = string
}

variable "image" {
  description = "Container image for the scheduled job"
  type        = string
  default     = null
}

variable "command" {
  description = "Command to execute"
  type        = list(string)
  default     = null
}

variable "environment" {
  description = "Environment variables"
  type        = map(string)
  default     = {}
}
