variable "context" {
  description = "Docker build context path"
  type        = string
}

variable "dockerfile" {
  description = "Path to Dockerfile (relative to context)"
  type        = string
  default     = null
}

variable "target" {
  description = "Docker build target stage"
  type        = string
  default     = null
}

variable "args" {
  description = "Docker build arguments"
  type        = map(string)
  default     = null
}

variable "token" {
  description = "Vercel API token (used for registry auth)"
  type        = string
  sensitive   = true
}

variable "team_id" {
  description = "Vercel team ID"
  type        = string
  default     = ""
}
