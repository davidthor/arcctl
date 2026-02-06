variable "name" {
  description = "Key name for identification"
  type        = string
}

variable "key_type" {
  description = "Key type: rsa, ecdsa, or symmetric"
  type        = string
  validation {
    condition     = contains(["rsa", "ecdsa", "symmetric"], var.key_type)
    error_message = "key_type must be one of: rsa, ecdsa, symmetric"
  }
}

variable "key_size" {
  description = "Key size (RSA: 2048/4096 bits, ECDSA: 256/384/521, Symmetric: 128/256 bits)"
  type        = number
  default     = null
}
