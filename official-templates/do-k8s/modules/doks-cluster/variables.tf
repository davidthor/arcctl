variable "name" {
  description = "Name of the Kubernetes cluster"
  type        = string
}

variable "region" {
  description = "DigitalOcean region"
  type        = string
}

variable "token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
}

variable "node_pool" {
  description = "Node pool configuration"
  type = object({
    size       = string
    min_nodes  = number
    max_nodes  = number
    auto_scale = bool
  })
  default = {
    size       = "s-2vcpu-4gb"
    min_nodes  = 2
    max_nodes  = 10
    auto_scale = true
  }
}
