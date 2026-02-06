variable "name" {
  description = "Vercel project name"
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
