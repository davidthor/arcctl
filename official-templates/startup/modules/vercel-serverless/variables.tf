variable "name" {
  description = "Deployment name"
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

variable "image" {
  description = "Container image to deploy"
  type        = string
  default     = null
}

variable "command" {
  description = "Override command"
  type        = list(string)
  default     = null
}

variable "entrypoint" {
  description = "Override entrypoint"
  type        = list(string)
  default     = null
}

variable "environment" {
  description = "Environment variables"
  type        = map(string)
  default     = {}
}

variable "cpu" {
  description = "CPU allocation"
  type        = string
  default     = null
}

variable "memory" {
  description = "Memory allocation"
  type        = string
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
  description = "Custom domain alias for this deployment"
  type        = string
  default     = ""
}

variable "replicas" {
  description = "Number of replicas"
  type        = number
  default     = null
}
