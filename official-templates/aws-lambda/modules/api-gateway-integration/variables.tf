variable "api_id" {
  description = "API Gateway ID"
  type        = string
}

variable "stage" {
  description = "API Gateway stage name"
  type        = string
}

variable "name" {
  description = "Service name"
  type        = string
}

variable "target" {
  description = "Target endpoint"
  type        = string
}

variable "target_type" {
  description = "Target type"
  type        = string
}

variable "port" {
  description = "Target port"
  type        = number
}
