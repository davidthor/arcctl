variable "name" {
  description = "Function name"
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

variable "context" {
  description = "Build context path"
  type        = string
  default     = "."
}

variable "command" {
  description = "Custom build command"
  type        = string
  default     = null
}

variable "environment" {
  description = "Environment variables"
  type        = map(string)
  default     = {}
}

variable "framework" {
  description = "Framework hint (e.g., nextjs, remix, astro)"
  type        = string
  default     = null
}

variable "memory" {
  description = "Function memory limit in MB"
  type        = string
  default     = null
}

variable "timeout" {
  description = "Function timeout in seconds"
  type        = number
  default     = null
}

variable "region" {
  description = "Deployment region"
  type        = string
  default     = "iad1"
}

variable "vercel_env" {
  description = "Vercel environment target (production or preview)"
  type        = string
  default     = "preview"
}

variable "alias" {
  description = "Custom domain alias for this function"
  type        = string
  default     = ""
}
