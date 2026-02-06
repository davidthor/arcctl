variable "name" {
  description = "Name for the key"
  type        = string
}

variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
}

variable "key_type" {
  description = "Key type: rsa, ecdsa, or symmetric"
  type        = string
  validation {
    condition     = contains(["rsa", "ecdsa", "symmetric"], var.key_type)
    error_message = "Key type must be rsa, ecdsa, or symmetric."
  }
}

variable "key_size" {
  description = "Key size in bits (e.g., 2048, 4096 for RSA; 256 for symmetric)"
  type        = number
  default     = 2048
}
