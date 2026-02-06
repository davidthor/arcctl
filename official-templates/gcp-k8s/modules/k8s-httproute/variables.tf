variable "name" {
  description = "Name for the HTTPRoute"
  type        = string
}

variable "namespace" {
  description = "Kubernetes namespace"
  type        = string
}

variable "gateway_name" {
  description = "Name of the parent Gateway"
  type        = string
}

variable "kubeconfig" {
  description = "Kubeconfig for the Kubernetes cluster"
  type = object({
    host                   = string
    cluster_ca_certificate = string
    token                  = string
  })
  sensitive = true
}

variable "target" {
  description = "Backend service name to route to"
  type        = string
}

variable "target_type" {
  description = "Type of target resource"
  type        = string
  default     = "service"
}

variable "port" {
  description = "Backend service port"
  type        = number
}

variable "hostnames" {
  description = "List of hostnames for the route"
  type        = list(string)
  default     = []
}

variable "path" {
  description = "URL path prefix to match"
  type        = string
  default     = null
}
