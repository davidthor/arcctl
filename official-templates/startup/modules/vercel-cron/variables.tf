variable "name" {
  description = "Cron job name"
  type        = string
}

variable "token" {
  description = "Vercel API token"
  type        = string
  sensitive   = true
}

variable "team_id" {
  description = "Vercel team ID (optional for personal accounts)"
  type        = string
  default     = ""
}

variable "project_id" {
  description = "Vercel project ID"
  type        = string
}

variable "schedule" {
  description = "Cron schedule expression (e.g., '0 */6 * * *')"
  type        = string
}

variable "image" {
  description = "Container image for the cron job"
  type        = string
  default     = null
}

variable "command" {
  description = "Command to run"
  type        = list(string)
  default     = null
}

variable "environment" {
  description = "Environment variables for the cron job"
  type        = map(string)
  default     = {}
}

variable "region" {
  description = "Deployment region"
  type        = string
  default     = "iad1"
}
