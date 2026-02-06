variable "name" {
  description = "Key name (used for identification)"
  type        = string
}

variable "key_type" {
  description = "Key type: rsa, ecdsa, or symmetric"
  type        = string

  validation {
    condition     = contains(["rsa", "ecdsa", "symmetric"], var.key_type)
    error_message = "key_type must be rsa, ecdsa, or symmetric"
  }
}

variable "key_size" {
  description = "Key size in bits (RSA: 2048/4096, ECDSA: 256/384, Symmetric: 128/256)"
  type        = number
  default     = 256
}
