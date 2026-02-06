variable "context" {
  description = "Docker build context path"
  type        = string
}

variable "dockerfile" {
  description = "Dockerfile path relative to context"
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

variable "registry" {
  description = "Container registry URL"
  type        = string
}

variable "token" {
  description = "Registry authentication token"
  type        = string
  sensitive   = true
}
