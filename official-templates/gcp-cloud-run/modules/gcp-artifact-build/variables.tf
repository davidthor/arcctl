variable "context" {
  description = "Docker build context path"
  type        = string
}

variable "dockerfile" {
  description = "Path to Dockerfile relative to context"
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
  default     = {}
}

variable "registry" {
  description = "Artifact Registry repository URL"
  type        = string
}

variable "project" {
  description = "GCP project ID"
  type        = string
}
