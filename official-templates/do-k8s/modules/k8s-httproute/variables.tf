variable "name" {
  description = "HTTPRoute name"
  type        = string
  default     = null
}

variable "namespace" {
  description = "Kubernetes namespace"
  type        = string
}

variable "kubeconfig" {
  description = "Kubernetes cluster kubeconfig content"
  type        = string
  sensitive   = true
}

variable "gateway_name" {
  description = "Gateway resource name"
  type        = string
}

variable "hostnames" {
  description = "Hostnames for the route"
  type        = list(string)
  default     = []
}

variable "path" {
  description = "URL path prefix"
  type        = string
  default     = null
}

variable "service" {
  description = "Target service name"
  type        = string
  default     = null
}

variable "deployment" {
  description = "Target deployment name"
  type        = string
  default     = null
}

variable "function" {
  description = "Target function name"
  type        = string
  default     = null
}

variable "port" {
  description = "Target service port"
  type        = number
  default     = null
}
