variable "name" {
  description = "Secret name (used as environment variable key)"
  type        = string
}

variable "value" {
  description = "Secret value"
  type        = string
  sensitive   = true
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

variable "vercel_env" {
  description = "Vercel environment target (production or preview)"
  type        = string
  default     = "preview"
}
