variable "name" {
  description = "Service name"
  type        = string
}

variable "target" {
  description = "Target deployment or function ID"
  type        = string
}

variable "target_type" {
  description = "Target type (deployment or function)"
  type        = string
  default     = "deployment"
}

variable "port" {
  description = "Service port"
  type        = number
  default     = null
}

variable "protocol" {
  description = "Service protocol (http or https)"
  type        = string
  default     = "https"
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
