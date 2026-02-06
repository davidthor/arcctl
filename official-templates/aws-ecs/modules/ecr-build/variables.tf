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

variable "region" {
  description = "AWS region"
  type        = string
}
