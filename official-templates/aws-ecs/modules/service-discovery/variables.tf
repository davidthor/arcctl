variable "name" {
  description = "Service name"
  type        = string
}

variable "namespace" {
  description = "Service discovery namespace"
  type        = string
}

variable "target" {
  description = "Target resource identifier"
  type        = string
}

variable "target_type" {
  description = "Target type (e.g., deployment, function)"
  type        = string
}

variable "port" {
  description = "Service port"
  type        = number
}
